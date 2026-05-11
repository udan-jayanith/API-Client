package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_body_header struct {
	gui.DefaultWidget

	auto_wrap struct {
		text   widget.Text
		toggle widget.Toggle
	}
	format       CommonWidgets.ButtonWithTooltip
	content_type widget.Combobox
}

func (w *request_body_header) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	input_widget := &w.content_type
	input_widget.SetAllowFreeInput(true)
	input_widget.SetItems([]string{"application/json", "application/octet-stream", "text/html", "text/plain", "image/png", "image/jpeg"})
	adder.AddWidget(input_widget)

	w.auto_wrap.text.SetValue("Auto wrap")
	w.auto_wrap.text.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&w.auto_wrap.text)
	adder.AddWidget(&w.auto_wrap.toggle)

	w.format.SetText("Format")
	w.format.SetTooltip("Ctrl+S")
	adder.AddWidget(&w.format)
	return nil
}

func (rbh *request_body_header) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	toggle_size := gui.FixedSize(u*2 - u/3)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rbh.auto_wrap.text,
			},
			{
				Widget: &rbh.auto_wrap.toggle,
				Size:   toggle_size,
			},
			{
				Widget: &rbh.format,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rbh.content_type,
				Size:   gui.FlexibleSize(1),
				// TODO: use Min and Max size
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rbh *request_body_header) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += rbh.auto_wrap.text.Measure(ctx, constraints).X + gap
		point.X += rbh.auto_wrap.toggle.Measure(ctx, constraints).X + gap

		point.X += rbh.content_type.Measure(ctx, constraints).X + gap
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u
	}
	return point
}

type request_body_widget struct {
	gui.DefaultWidget
	header request_body_header

	body CommonWidgets.WidgetWithLazyLoading[*CommonWidgets.TextInputWithContextMenu]
}

func (w *request_body_widget) SetLazyLoad(lazy_load bool) {
	w.body.SetLazyLoad(lazy_load)
}

func (w *request_body_widget) LazyLoad() bool {
	return w.body.LazyLoad()
}

func (w *request_body_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.header)

	body := w.body.Widget()
	if w.header.auto_wrap.toggle.Value() {
		body.SetWrapMode(widget.WrapModeAnywhere)
	}
	body.SetMultiline(true)
	body.SetEditable(true)
	adder.AddWidget(&w.body)
	return nil
}

func (w *request_body_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.header,
			},
			{
				Widget: &w.body,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (body *request_body_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	var point image.Point
	point.Y = gap

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = body.header.Measure(ctx, constraints).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += body.header.Measure(ctx, constraints).Y
		point.Y += body.body.Measure(ctx, constraints).Y
	}

	return point
}

func (body *request_body_widget) SetBody(content string) {
	body.body.Widget().ForceSetValue(content)
}

func (body *request_body_widget) Body() string {
	return body.body.Widget().Value()
}

func (body *request_body_widget) OnAutowrapToggle(fn func(ctx *gui.Context, value bool)) {
	body.header.auto_wrap.toggle.OnValueChanged(fn)
}

func (body *request_body_widget) SetAutowrap(autowrap bool) {
	body.header.auto_wrap.toggle.SetValue(autowrap)
}

func (body *request_body_widget) OnContentTypeChanged(fn func(context *gui.Context, value string, committed bool)) {
	body.header.content_type.OnValueChanged(fn)
}

func (body *request_body_widget) SetContentType(content_type requests_handler.ContentType) {
	body.header.content_type.SetValue(string(content_type))
}
