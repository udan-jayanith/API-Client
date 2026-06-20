package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	attr "Zbolt/pages/request/requests-handler/attributes"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/go-units"
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_widget struct {
	gui.DefaultWidget
	header_widget response_header_widget
	tab_container CommonWidgets.TabContainer[string]
	tab_content   struct {
		response_headers CommonWidgets.WidgetWithLazyLoading[*HttpHeaderTable]
		response_body    response_body_widget
	}
}

func (rw *response_widget) SetLazyLoading(body, headers bool) {
	rw.tab_content.response_body.SetLazyLoad(body)
	rw.tab_content.response_headers.SetLazyLoad(headers)
}

func (rw *response_widget) Clear() {
	rw.header_widget.clear()
	rw.tab_content.response_headers.Widget().SetRowsCheck(nil)
	rw.tab_content.response_body.SetBody(nil)
	rw.tab_content.response_body.SetContentType("")
}

func (rw *response_widget) OnAutowrapToggle(fn func(ctx *gui.Context, value bool)) {
	rw.tab_content.response_body.OnAutowrapToggle(fn)
}

func (rw *response_widget) OnFormatToggle(fn func(ctx *gui.Context, value bool)) {
	rw.tab_content.response_body.OnFormatToggle(fn)
}

func (rw *response_widget) SetAutowrap(autowrap bool) {
	rw.tab_content.response_body.SetAutowrap(autowrap)
}

func (rw *response_widget) SetFormat(format bool) {
	rw.tab_content.response_body.SetFormat(format)
}

// SetStatus sets the http status code
func (rw *response_widget) SetStatus(status_code int) {
	if status_code < 200 {
		rw.header_widget.status.SetValue("")
		return
	}
	status := fmt.Sprintf("%v %s", status_code, http.StatusText(status_code))
	rw.header_widget.status.SetValue(status)
}

// SetResponseTime sets the http response time
func (rw *response_widget) SetResponseTime(response_time time.Duration) {
	rw.header_widget.response_time.SetValue(response_time.Round(time.Millisecond).String())
}

func (rw *response_widget) SetHTTPVersion(version requests_handler.Version) {
	rw.header_widget.proto.SetValue(fmt.Sprintf("HTTP v%v.%v", version.Major, version.Minor))
}

func (rw *response_widget) SetContentLength(lenght int) {
	rw.header_widget.content_lenght.SetValue(units.HumanSize(float64(lenght)))
}

func (rw *response_widget) SetHeaders(headers []attr.Attribute) {
	rw.tab_content.response_headers.Widget().SetRows(headers)
}

func (rw *response_widget) SetResponseBody(body *requests_handler.HTTP_Response_Body) {
	rw.tab_content.response_body.SetBody(body)
	rw.tab_content.response_body.SetContentType(body.ContentType)
}

func (rw *response_widget) SetSelectedTab(index int) {
	rw.set_tab_items()
	rw.tab_container.SelectTab(index)
}

func (rw *response_widget) SelectedTab() int {
	_, index := rw.tab_container.SelectedTabContainer()
	return index
}

func (rw *response_widget) set_tab_items() {
	if rw.tab_container.Count() != 0 {
		return
	}
	rw.tab_container.SetItems([]CommonWidgets.TabContainerItem[string]{
		{
			TabItem: CommonWidgets.TabItem[string]{
				Text:  "Body",
				Value: "body",
			},
			Widget: &rw.tab_content.response_body,
		},
		{
			TabItem: CommonWidgets.TabItem[string]{
				Text:  "Header",
				Value: "header",
			},
			Widget: &rw.tab_content.response_headers,
		},
	})
}

func (rw *response_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.header_widget)
	rw.set_tab_items()

	headers_table := rw.tab_content.response_headers.Widget()
	headers_table.DisableCheckbox(true)
	headers_table.DisableDelete(true)
	headers_table.KeyEditable(false)
	headers_table.ValueEditable(false)
	headers_table.AutoAddRow(false)

	adder.AddWidget(&rw.tab_container)

	return nil
}

func (rw *response_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	main_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.header_widget,
			},
			{
				Widget: &rw.tab_container,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	main_layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
