package requests_handler

import (
	"Zbolt/internal/readreader"
	attr "Zbolt/pages/request/requests-handler/attributes"
	lazy_atomic "Zbolt/pages/request/requests-handler/internal/lazy-atomic"
	url_utils "Zbolt/pages/request/requests-handler/url-utils"
	"net/url"
	"sync/atomic"
	"time"
)

type URL struct {
	BaseURL string `json:"base-url"` // Everything before the path.

	Path struct {
		RawPath string `json:"raw-path"`
		Pattern struct {
			Pattern    string           `json:"pattern"`
			Attributes []attr.Attribute `json:"attributes"`
		} `json:"pattern"`
	} `json:"path"` // Both path and pattern can't exists at once.
}

func (u *URL) IsPattern() bool {
	return len(u.Path.Pattern.Attributes) > 0
}

// EncodedPath returns the encoded path
func (u *URL) URL_Path() string {
	if !u.IsPattern() {
		return u.Path.RawPath
	}

	pattern, _ := url_utils.ParsePattern(u.Path.Pattern.Pattern)
	for _, attr := range u.Path.Pattern.Attributes {
		pattern.Set(attr.Key, attr.Value)
	}
	return pattern.Path()
}

// RawPath returns the raw-path if exists otherwise returns the raw-path-pattern without being encoded.
func (u *URL) RawPath() string {
	if u.IsPattern() {
		return u.Path.Pattern.Pattern
	}
	return u.Path.RawPath
}

func (u *URL) SetPattern(pattern string, attributes []attr.Attribute) {
	u.Path.RawPath = ""
	u.Path.Pattern.Pattern = pattern
	u.Path.Pattern.Attributes = attributes
}

func (u *URL) SetPath(path string) {
	u.Path.RawPath = path
	u.Path.Pattern.Pattern = ""
	u.Path.Pattern.Attributes = []attr.Attribute{}
}

// TODO: Implement Io.Closer
type HTTP_Data struct {
	Method string `json:"method"` // HTTP method

	URL URL `json:"url"`

	Parameters    []attr.AttrCheck `json:"parameters"`
	Headers       []attr.AttrCheck `json:"headers"`
	Body          string           `json:"body"`
	RequestConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"request-config"`

	ResponseConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`

	selected_request_tab int
	// TODO: Store wether the url panel is opne or not

	request struct {
		is_fetching, canceled, headers_changed atomic.Bool
		on_complete                            chan error
		cancel                                 chan struct{}
	}

	// Do not use this directly use ResponseData function.
	response_data lazy_atomic.Value[HTTP_Response_Data]
}

func (data *HTTP_Data) SetSelectedRequestTab(index int) {
	data.selected_request_tab = index
}

func (data *HTTP_Data) SelectedRequestTab() int {
	return data.selected_request_tab
}

/*
Adapted from Golang net/http package.
example: username=edger&age=20
*/
func (data *HTTP_Data) EncodedParameters() string {
	return url_utils.EncodeParameters(data.Parameters)
}

// GetUrl return the full url.
func (data *HTTP_Data) FullURL() *url.URL {
	u, _ := url.Parse(data.URL.BaseURL)
	u.Path = data.URL.URL_Path()
	u.RawQuery = data.EncodedParameters()
	return u
}

func (data *HTTP_Data) ResponseData(fn func(value *HTTP_Response_Data)) {
	data.response_data.Mutex.Lock()
	defer data.response_data.Mutex.Unlock()
	fn(data.response_data.LoadUnsafe())
}

func (data *HTTP_Data) Close() error {
	data.ResponseData(func(value *HTTP_Response_Data) {
		value.Body.Content().Close()
	})
	return nil
}

type HTTP_Response_Body struct {
	ContentType       ContentType
	is_formated       bool
	formattee_content *readreader.ReadReader
	content           *readreader.ReadReader
}

func (c *HTTP_Response_Body) set_content(rr *readreader.ReadReader) {
	if c.content != nil {
		c.content.Close()
	}
	if c.formattee_content != nil {
		c.formattee_content.Close()
		c.formattee_content = nil
	}

	c.content = rr
	c.is_formated = false
}

func (c *HTTP_Response_Body) Content() *readreader.ReadReader {
	return c.content
}

type Version struct {
	Major, Minor int
}

type HTTP_Response_Data struct {
	Status_code   int
	ResponseTime  time.Duration
	ContentLenght int // In bytes
	Version       Version
	Headers       []attr.Attribute
	Body          HTTP_Response_Body

	SelectedResponseTab int
}
