package http_widget

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	message_model "Zbolt/message-model"
	requests_handler "Zbolt/pages/request/requests-handler"
	attr "Zbolt/pages/request/requests-handler/attributes"
	url_utils "Zbolt/pages/request/requests-handler/url-utils"
	"image"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	opener "codeberg.org/udan-jayanith/Opener"
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/dialog"
)

type HTTP_Widget struct {
	gui.DefaultWidget

	loading_bar     CommonWidgets.InfiniteLoadingBar
	request_widget  request_widget
	response_widget response_widget

	url_panel_widget url_panel_widget
	popup_widget     widget.Popup

	req         *requests_handler.Request
	data        *requests_handler.HTTP_Data
	t           time.Time
	is_fetching bool
}

func (brp *HTTP_Widget) RequestType() requests_handler.RequestType {
	return requests_handler.HTTP
}

// SetReq runs when switching tabs and tab data are passed to this.
func (brp *HTTP_Widget) SetReq(req *requests_handler.Request) {
	if req.Type != requests_handler.HTTP {
		panic("Invalid request type")
	}

	if brp.req == req {
		return
	}
	brp.req = req

	data, ok := req.Data().(*requests_handler.HTTP_Data)
	if !ok {
		panic("Invalid data type")
	}

	// Setup request widget
	brp.data = data
	brp.is_fetching = data.IsFetching()
	brp.setup_request_widget()
	brp.setup_response_widget()
	gui.RequestRebuild(brp)
}

func (brp *HTTP_Widget) setup_request_widget() {
	data := brp.data
	brp.request_widget.SetHeaders(data.Headers)
	brp.request_widget.SetParameters(data.Parameters)
	brp.request_widget.SetAutowrap(data.RequestConfig.AutoWrap)
	brp.request_widget.SetBody(data.Body)
	if data.Method == "" {
		data.Method = "Get"
	}
	brp.request_widget.SetMethod(data.Method)
	brp.request_widget.SetContentType("")
	for _, h := range data.Headers {
		if h.Key == "Content-Type" && h.Checked {
			brp.request_widget.SetContentType(requests_handler.ContentType(h.Value))
			break
		}
	}
	brp.request_widget.SelectTab(data.SelectedRequestTab())

	u, err := url.Parse(data.URL.BaseURL)
	if err != nil {
		message_model.Show(err.Error(), message_model.Alert, nil)
		return
	}
	u.Path = data.URL.URL_Path()
	brp.request_widget.SetURL(u)
	brp.request_widget.DisableURLInput(data.URL.IsPattern())
}

func (brp *HTTP_Widget) setup_response_widget() {
	// Setup response widget
	data := brp.data
	data.ResponseData(func(res_data *requests_handler.HTTP_Response_Data) {
		if brp.is_fetching {
			brp.response_widget.SetLazyLoading(true, len(res_data.Headers) == 0)
		} else {
			brp.response_widget.SetLazyLoading(false, false)
		}

		brp.response_widget.SetHeaders(res_data.Headers)

		brp.response_widget.OnAutowrapToggle(func(ctx *gui.Context, value bool) {})
		brp.response_widget.OnFormatToggle(func(ctx *gui.Context, value bool) {})
		brp.response_widget.SetAutowrap(data.ResponseConfig.AutoWrap)
		brp.response_widget.SetFormat(data.ResponseConfig.Formate)
		brp.response_widget.SetResponseBody(&res_data.Body)
		brp.response_widget.SetSelectedTab(res_data.SelectedResponseTab)

		brp.response_widget.SetHTTPVersion(res_data.Version)
		brp.response_widget.SetResponseTime(res_data.ResponseTime)
		brp.response_widget.SetStatus(res_data.Status_code)
		brp.response_widget.SetContentLength(res_data.ContentLenght)
	})
	brp.response_widget.SearchHeaders(data.ResponseHeaderSearchQuery())
}

// TODO: SyncData should be run to save data before switching tabs, closing tabs or closing the app.
func (brp *HTTP_Widget) SyncData() {
	brp.data.Parameters = brp.request_widget.Parameters()
	brp.data.Headers = brp.request_widget.Headers()
	brp.data.Body = brp.request_widget.Body()

	// TODO: sync response data
	brp.data.SetSelectedRequestTab(brp.request_widget.SelectedTab())
	brp.data.ResponseData(func(value *requests_handler.HTTP_Response_Data) {
		value.SelectedResponseTab = brp.response_widget.SelectedTab()
	})
	brp.data.SetResponseHeaderSearchQuery(brp.response_widget.HeaderSearchQuery())
}

func (brp *HTTP_Widget) url_panel_bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	b := widgetBounds.Bounds()
	draw_area := b

	w := brp.url_panel_widget.Measure(ctx, gui.FixedWidthConstraints(b.Dx())).X
	draw_area.Min.X += (b.Dx() / 2) - (w / 2)
	draw_area.Max.X = draw_area.Min.X + w

	h := brp.url_panel_widget.Measure(ctx, gui.FixedHeightConstraints(b.Dy())).Y
	draw_area.Min.Y += (b.Dy() / 2) - (h / 2)
	draw_area.Max.Y = draw_area.Min.Y + h

	return draw_area
}

func (brp *HTTP_Widget) on_url_panel_open(ctx *gui.Context) {
	u, _ := url.Parse(brp.data.URL.BaseURL)
	brp.url_panel_widget.Set(u.Scheme, u.Host, brp.data.URL.RawPath(), brp.data.URL.Path.Pattern.Attributes)
	brp.popup_widget.SetOpen(true)
}

func (brp *HTTP_Widget) on_url_panel_close(_ *gui.Context, _ widget.PopupCloseReason) {
	u, err := url.Parse(brp.url_panel_widget.URL())
	if err != nil {
		message_model.Show(err.Error(), message_model.Alert, nil)
		return
	}
	brp.request_widget.SetURL(u)

	pattern, query_list := brp.url_panel_widget.Pattern()
	if len(query_list) > 0 {
		brp.data.URL.SetPattern(pattern, query_list)
		brp.request_widget.DisableURLInput(true)
	} else {
		brp.data.URL.SetPath(u.Path)
		brp.request_widget.DisableURLInput(false)
	}
	u.Path = ""
	url_utils.CleanURL(u)
	brp.data.URL.BaseURL = u.String()

	brp.url_panel_widget.Clear()
}

func (brp *HTTP_Widget) on_request_button_clicked(ctx *gui.Context, value string) {
	if value == RequestButton {
		brp.SyncData()
		brp.data.Do()
		brp.response_widget.Clear()
		brp.setup_request_widget()
		brp.is_fetching = true
	} else if err := brp.data.CancelRequest(); err != nil {
		message_model.Show(err.Error(), message_model.Alert, nil)
	}
}

func (brp *HTTP_Widget) on_url_input_changed(_ *gui.Context, u_str string, committed bool) {
	if !committed || brp.data.URL.IsPattern() || u_str == "" {
		return
	}

	u_str = strings.TrimSpace(u_str)
	if url_utils.IsJustPortNumber(u_str) {
		u_str = "http://localhost" + u_str
	} else if strings.HasPrefix(u_str, "localhost") {
		u_str = "http://" + u_str
	} else if !strings.HasPrefix(u_str, "http") {
		u_str = "http://" + u_str
	}

	u, err := url.Parse(u_str)
	if err != nil {
		message_model.Show(err.Error(), message_model.Alert, nil)
		return
	}

	brp.request_widget.SetURL(u)

	url_utils.CleanURL(u)
	brp.data.URL.SetPath(u.Path)
	u.Path = ""
	brp.data.URL.BaseURL = u.String()
}

func (brp *HTTP_Widget) on_open_externally(context *gui.Context) {
	brp.data.ResponseData(func(value *requests_handler.HTTP_Response_Data) {
		if value.Body.Content() == nil {
			message_model.Show("No content found to open", message_model.Alert, nil)
			return
		}
		r := value.Body.Content().NewReader()
		defer r.Close()
		_, err := opener.OpenStream(r, value.Body.ContentType.Extension())
		if err != nil {
			message_model.Show(err.Error(), message_model.Alert, nil)
		}
	})
}

func (brp *HTTP_Widget) on_save_as(context *gui.Context) {
	go func() {
		var content_type requests_handler.ContentType
		var r io.ReadCloser
		// TODO: use the content-disposition header to determine the defualt file name.

		brp.data.ResponseData(func(value *requests_handler.HTTP_Response_Data) {
			content_type = value.Body.ContentType
			if value.Body.Content() == nil {
				return
			}
			r = value.Body.Content().NewReader()
		})
		if r == nil {
			dialog.Message("No content found to save").Info()
			return
		}
		defer r.Close()

		extensions := content_type.Extensions()
		path, err := dialog.File().Title("Save As").Filter("", extensions...).Save()
		if err != nil {
			return
		}

		file, err := os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
		file.ReadFrom(r)
	}()
}

func (brp *HTTP_Widget) on_request_content_type_changed(ctx *gui.Context, content_type string, committed bool) {
	if !committed {
		return
	}

	brp.data.Headers = brp.request_widget.Headers()
	for i, header := range brp.data.Headers {
		if header.Key == "Content-Type" {
			brp.data.Headers[i].Checked = content_type != ""
			brp.data.Headers[i].Value = content_type
			brp.request_widget.SetHeaders(brp.data.Headers)
			return
		}
	}

	if content_type != "" {
		brp.data.Headers = append([]attr.AttrCheck{
			{
				Checked: true,
				Key:     "Content-Type",
				Value:   content_type,
			},
		}, brp.data.Headers...)
		brp.request_widget.SetHeaders(brp.data.Headers)
	}
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	brp.request_widget.OnOpenIn(brp.on_url_panel_open)
	adder.AddWidget(&brp.popup_widget)
	if brp.popup_widget.IsOpen() {
		brp.popup_widget.SetCloseByClickingOutside(true)
		brp.popup_widget.SetContent(&brp.url_panel_widget)
		brp.popup_widget.SetModal(true)
		brp.url_panel_widget.OnDoneButtonClicked(func(context *gui.Context) {
			brp.popup_widget.SetOpen(false)
		})
		brp.popup_widget.OnClose(brp.on_url_panel_close)
	}

	brp.request_widget.OnMethodChanged(func(method string) {
		brp.data.Method = method

		brp.data.Headers = brp.request_widget.Headers()
		method = strings.ToLower(method)
		allowed := method == "post" || method == "put" || method == "patch"
		for i, header := range brp.data.Headers {
			if header.Key == "Content-Type" {
				brp.data.Headers[i].Checked = allowed && header.Value != ""
				break
			}
		}
		brp.request_widget.SetHeaders(brp.data.Headers)
	})

	brp.request_widget.OnContentTypeChanged(brp.on_request_content_type_changed)

	brp.request_widget.OnURLInputChanged(brp.on_url_input_changed)

	brp.request_widget.OnRequestButtonClicked(brp.on_request_button_clicked)

	brp.request_widget.OnAutowrap(func(ctx *gui.Context, value bool) {
		brp.data.RequestConfig.AutoWrap = value
	})

	brp.response_widget.OnAutowrapToggle(func(ctx *gui.Context, value bool) {
		brp.data.ResponseConfig.AutoWrap = value
		brp.response_widget.SetAutowrap(value)
		brp.data.ResponseData(func(value *requests_handler.HTTP_Response_Data) {
			brp.response_widget.SetResponseBody(&value.Body)
		})
	})

	brp.response_widget.OnFormatToggle(func(ctx *gui.Context, value bool) {
		brp.data.SetFormatResponseBody(value)
		brp.response_widget.SetFormat(value)
		brp.data.ResponseData(func(value *requests_handler.HTTP_Response_Data) {
			brp.response_widget.SetResponseBody(&value.Body)
		})
	})

	brp.response_widget.OnOpenExternally(brp.on_open_externally)
	brp.response_widget.OnSaveAs(brp.on_save_as)

	if brp.is_fetching {
		adder.AddWidget(&brp.loading_bar)
	} else {
	}

	if brp.data.HeadersChanged() {
		brp.setup_response_widget()
	}

	select {
	case err, ok := <-brp.data.OnComplete():
		if !ok {
			break
		}
		brp.is_fetching = false
		brp.setup_response_widget()
		if err != nil {
			message_model.Show(err.Error(), message_model.Alert, nil)
		}
	default:
	}

	if brp.is_fetching {
		brp.request_widget.SetRequestButtonText(CancelButton)
	} else {
		brp.request_widget.SetRequestButtonText(RequestButton)
	}
	adder.AddWidget(&brp.request_widget)
	adder.AddWidget(&brp.response_widget)
	return nil
}

func (brp *HTTP_Widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if brp.popup_widget.IsOpen() {
		brp.popup_widget.SetBackgroundBounds(ctx.AppBounds())
		layouter.LayoutWidget(&brp.popup_widget, brp.url_panel_bounds(ctx, widgetBounds))
	}

	if brp.is_fetching {
		loading_bar_bounds := widgetBounds.Bounds()
		loading_bar_size := brp.loading_bar.Measure(ctx, gui.Constraints{})
		loading_bar_bounds.Max.Y = loading_bar_bounds.Min.Y
		loading_bar_bounds.Min.Y -= loading_bar_size.Y
		layouter.LayoutWidget(&brp.loading_bar, loading_bar_bounds)
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &brp.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
