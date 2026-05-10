package requests_handler

import (
	attr "API-Client/pages/request/requests-handler/attributes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func (data *HTTP_Data) open_request() {
	data.request.cancel = make(chan struct{}, 1)
	data.request.on_complete = make(chan error)
	data.request.is_fetching.Store(true)
	data.request.canceled.Store(false)
	data.set_response_data(HTTP_Response_Data{})
}

func (data *HTTP_Data) close_request(err error) {
	close(data.request.cancel)
	data.request.is_fetching.Store(false)
	go func() {
		data.request.on_complete <- err
		close(data.request.on_complete)
	}()
}

func (data *HTTP_Data) IsFetching() bool {
	return data.request.is_fetching.Load()
}

func (data *HTTP_Data) HeadersChanged() bool {
	return data.request.headers_changed.CompareAndSwap(true, false)
}

func (data *HTTP_Data) OnComplete() chan error {
	return data.request.on_complete
}

func (data *HTTP_Data) CancelRequest() error {
	if !data.IsFetching() {
		return errors.New("HTTP request is no fetching")
	} else if data.request.canceled.Load() {
		return errors.New("Request is already being canceled")
	}
	data.request.canceled.Store(true)
	data.request.cancel <- struct{}{}
	return nil
}

// Do performs the http request
// Response data can be revised through ResponseData method
// Calling Do updates Headers so headers must be update in the HTTP_widget
func (data *HTTP_Data) Do() bool {
	if data.request.is_fetching.Load() {
		panic("Request is already being requested")
	}
	data.open_request()

	method := strings.ToUpper(data.Method)
	var body io.Reader
	if method == "POST" || method == "PUT" || method == "PATCH" {
		body = strings.NewReader(data.Body)
	}

	req, err := http.NewRequest(method, data.FullURL().String(), body)
	if err != nil {
		data.close_request(err)
		return false
	}
	for _, header := range data.Headers {
		if header.Checked {
			req.Header.Set(header.Key, header.Value)
		}
	}

	go data.do(req)

	return true
}

func (data *HTTP_Data) set_response_data(res_data HTTP_Response_Data) {
	body_content_copied := make([]byte, len(res_data.Body.content))
	copy(body_content_copied, res_data.Body.content)

	headers_copied := make([]attr.AttrCheck, len(res_data.Headers))
	copy(headers_copied, res_data.Headers)

	data.ResponseData(func(value *HTTP_Response_Data) {
		*value = res_data
		res_data.SelectedResponseTab = value.SelectedResponseTab
		value.Headers = headers_copied
		value.Body.content = body_content_copied
	})

	data.request.headers_changed.Store(true)
}

func (data *HTTP_Data) do(req *http.Request) {
	res_data := HTTP_Response_Data{}
	response_time := time.Now()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		data.close_request(err)
		return
	}
	defer res.Body.Close()

	res_data.Status_code = res.StatusCode
	res_data.Version = Version{
		Major: res.ProtoMajor,
		Minor: res.ProtoMinor,
	}
	res_data.Body.ContentType = ContentType(res.Header.Get("Content-Type"))
	res_data.Headers = http_headers_to_attr_check(res.Header)
	res_data.ResponseTime = time.Since(response_time)
	data.set_response_data(res_data)

	body_content := make([]byte, 0, 2048)
	buffer := make([]byte, 1024)
	update_time := time.Now()

loop:
	for {
		n, e := res.Body.Read(buffer)
		if e != nil && e != io.EOF {
			err = e
			break
		} else if e == io.EOF && n == 0 {
			break
		}

		select {
		case <-data.request.cancel:
			break loop
		default:
		}

		body_content = append(body_content, buffer[:n]...)
		res_data.ContentLenght = len(body_content)
		res_data.ResponseTime = time.Since(response_time)
		if time.Since(update_time).Milliseconds() >= 500 {
			data.set_response_data(res_data)
		}
	}
	res_data.ContentLenght = len(body_content)
	res_data.ResponseTime = time.Since(response_time)

	// TODO: decode the body content if it uses some kind of encription
	res_data.Body.content = body_content
	data.set_response_data(res_data)
	data.close_request(err)
}
