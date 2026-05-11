package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	attr "Zbolt/pages/request/requests-handler/attributes"
	url_utils "Zbolt/pages/request/requests-handler/url-utils"
	"image"
	"net/url"
	"strings"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_widget struct {
	gui.DefaultWidget
	input_bar_widget          request_input_bar_widget
	on_input_bar_value_change func(ctx *gui.Context, text string, committed bool, by_user bool)
	on_method_change          func(method string)

	t           time.Time
	url_preview CommonWidgets.URLPreview

	tab         CommonWidgets.Tab
	tab_content struct {
		params, header  []attr.AttrCheck
		table           CommonWidgets.AttributeTable
		body            request_body_widget
		selected_widget gui.Widget
	}
}

// sets the http method
func (rw *request_widget) SetMethod(method string) {
	rw.input_bar_widget.select_method(method)
}

func (rw *request_widget) OnMethodChanged(fn func(method string)) {
	rw.on_method_change = fn
}

func (rw *request_widget) Method() string {
	return rw.input_bar_widget.method()
}

func (rw *request_widget) OnOpenIn(fn func(ctx *gui.Context)) {
	rw.input_bar_widget.on_open_in_clicked(fn)
}

func (rw *request_widget) OnAutowrap(fn func(ctx *gui.Context, value bool)) {
	rw.tab_content.body.OnAutowrapToggle(fn)
}

func (rw *request_widget) SetAutowrap(value bool) {
	rw.tab_content.body.SetAutowrap(value)
}

func (rw *request_widget) Body() string {
	return rw.tab_content.body.Body()
}

func (rw *request_widget) OnContentTypeChanged(fn func(context *gui.Context, value string, committed bool)) {
	rw.tab_content.body.OnContentTypeChanged(fn)
}

func (rw *request_widget) SetContentType(content_type requests_handler.ContentType) {
	rw.tab_content.body.SetContentType(content_type)
}

func (rw *request_widget) SetURL(u *url.URL) {
	raw_query := u.RawQuery
	url_utils.CleanURL(u)

	parameters, _ := url_utils.ParseParametersAsCheck(raw_query)
	merged_parameters := attr.MergeAttrCheckList(rw.Parameters(), parameters, true)
	rw.SetParameters(merged_parameters)
	rw.input_bar_widget.set_url_input_value(u.String())
	rw.update_url_preview()
}

func (rw *request_widget) FullURL() string {
	return rw.update_url_preview()
}

func (rw *request_widget) update_url_preview() string {
	var parameters []attr.AttrCheck
	_, tab := rw.tab.SelectedTab()
	if tab.Value == "parameters" {
		parameters = rw.tab_content.table.RowsCheck()
	} else {
		parameters = rw.tab_content.params
	}

	u, _ := url.Parse(rw.input_bar_widget.url_input_value())
	if u == nil {
		return rw.url_preview.URL()
	}
	url_utils.CleanURL(u)
	u.RawQuery = url_utils.EncodeParameters(parameters)

	u_str := u.String()
	rw.url_preview.SetURL(u_str)
	return u_str
}

// Value is 'Request' or 'Cancel'
func (rw *request_widget) OnRequestButtonClicked(fn func(ctx *gui.Context, value string)) {
	rw.input_bar_widget.on_request_button_clicked(fn)
}

// Value is 'Request' or 'Cancel'
func (rw *request_widget) SetRequestButtonText(value string) {
	rw.input_bar_widget.set_request_button_value(value)
}

func (rw *request_widget) DisableURLInput(disabled bool) {
	rw.input_bar_widget.disable_url_input(disabled)
}

func (rw *request_widget) OnURLInputChanged(fn func(context *gui.Context, text string, committed bool)) {
	rw.input_bar_widget.on_url_input_value_changed(fn)
}

func (rw *request_widget) SetParameters(parameters []attr.AttrCheck) {
	_, item := rw.tab.SelectedTab()
	if item.Value == "parameters" {
		rw.tab_content.table.SetRowsCheck(parameters)
	}

	rw.tab_content.params = parameters
}

func (rw *request_widget) Parameters() []attr.AttrCheck {
	_, selected_tab := rw.tab.SelectedTab()
	if selected_tab.Value == "parameters" {
		rw.tab_content.params = rw.tab_content.table.RowsCheck()
	}
	return rw.tab_content.params
}

func (rw *request_widget) SetHeaders(headers []attr.AttrCheck) {
	_, selected_tab := rw.tab.SelectedTab()
	if selected_tab.Value == "headers" {
		rw.tab_content.table.SetRowsCheck(headers)
	}
	rw.tab_content.header = headers
}

func (rw *request_widget) Headers() []attr.AttrCheck {
	_, selected_tab := rw.tab.SelectedTab()
	if selected_tab.Value == "headers" {
		rw.tab_content.header = rw.tab_content.table.RowsCheck()
	}
	return rw.tab_content.header
}

// SetBody set the http request body
func (rw *request_widget) SetBody(body string) {
	rw.tab_content.body.SetBody(body)
}

func (rw *request_widget) SelectedTab() int {
	rw.set_tab_items()
	i, _ := rw.tab.SelectedTab()
	return i
}

func (rw *request_widget) SelectTab(index int) {
	rw.set_tab_items()
	rw.tab.SelectTab(index)
}

func (rw *request_widget) set_tab_items() {
	method := strings.ToUpper(rw.input_bar_widget.method())
	selected_tab, _ := rw.tab.SelectedTab()

	if method == "POST" || method == "PUT" || method == "PATCH" {
		rw.tab.SetTabItems([]CommonWidgets.TabItem{
			{
				Text:  "Parameters",
				Value: "parameters",
			},
			{
				Text:  "Headers",
				Value: "headers",
			},
			{
				Text:  "Body",
				Value: "body",
			},
		})
		rw.tab.SelectTab(selected_tab)
	} else {
		rw.tab.SetTabItems([]CommonWidgets.TabItem{
			{
				Text:  "Parameters",
				Value: "parameters",
			},
			{
				Text:  "Headers",
				Value: "headers",
			},
		})
		rw.tab.SelectTab(min(selected_tab, 1))
	}
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.input_bar_widget)
	adder.AddWidget(&rw.url_preview)

	rw.input_bar_widget.on_method_changed(func(method string) {
		rw.set_tab_items()
		if rw.on_method_change != nil {
			rw.on_method_change(method)
		}
	})

	if rw.tab.ItemCount() == 0 {
		rw.set_tab_items()
	}

	rw.tab.OnSelect(func(from, to CommonWidgets.TabItemContainer, _ bool) {
		if from.Item.Value == "parameters" && to.Item.Value == "headers" {
			rw.tab_content.params = rw.tab_content.table.RowsCheck()
		} else if from.Item.Value == "headers" && to.Item.Value == "parameters" {
			rw.tab_content.header = rw.tab_content.table.RowsCheck()
		}

		if to.Item.Value == "parameters" {
			rw.tab_content.table.SetRowsCheck(rw.tab_content.params)
		} else if to.Item.Value == "headers" {
			rw.tab_content.table.SetRowsCheck(rw.tab_content.header)
		}
		// TODO: This runs every time when the build runs
	})

	_, selected_tab := rw.tab.SelectedTab()
	switch selected_tab.Value {
	case "parameters", "headers":
		rw.tab_content.selected_widget = &rw.tab_content.table
	case "body":
		rw.tab_content.selected_widget = &rw.tab_content.body
	default:
		panic("Unknown tab was selected")
	}

	adder.AddWidget(rw.tab_content.selected_widget)
	adder.AddWidget(&rw.tab)

	// TODO: implement a method on the input bar widget which returns whether the input bar widget is focused or not.
	if time.Since(rw.t).Seconds() >= 1 && !ctx.IsFocused(&rw.input_bar_widget) {
		rw.update_url_preview()
		rw.t = time.Now()
	}
	return nil
}

func (rw *request_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.input_bar_widget,
				Size:   gui.FixedSize(u),
			},
			{
				Widget: &rw.url_preview,
			},
			{
				Widget: &rw.tab,
			},
			{
				Widget: rw.tab_content.selected_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rw *request_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := rw.input_bar_widget.Measure(ctx, constraints)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += rw.url_preview.Measure(ctx, constraints).Y
		point.Y += rw.tab.Measure(ctx, constraints).Y
		point.Y += rw.tab_content.selected_widget.Measure(ctx, constraints).Y
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}

	return point
}
