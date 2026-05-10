package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	requests_handler "API-Client/pages/request/requests-handler"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_body_header struct {
	gui.DefaultWidget

	content_type widget.Text

	options struct {
		auto_wrap struct {
			text   widget.Text
			toggle widget.Toggle
		}
		format struct {
			text   widget.Text
			toggle widget.Toggle
		}
		open    CommonWidgets.ButtonWithTooltip
		on_open func(context *gui.Context, content_type string)
	}
}

func (w *response_body_header) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.content_type.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&w.content_type)

	w.options.open.SetText("Open with")
	w.options.open.SetTooltip("Open externally")
	if w.options.on_open != nil {
		w.options.open.OnUp(func(ctx *gui.Context) {
			w.options.on_open(ctx, w.content_type.Value())
		})
	}
	adder.AddWidget(&w.options.open)
	// Add open with button to the middle if the content type is unkown.

	{
		w.options.auto_wrap.text.SetValue("Auto wrap")
		w.options.auto_wrap.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&w.options.auto_wrap.text)

		adder.AddWidget(&w.options.auto_wrap.toggle)
	}
	{
		w.options.format.text.SetValue("Format")
		w.options.format.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&w.options.format.text)

		adder.AddWidget(&w.options.format.toggle)
	}
	return nil
}

func (w *response_body_header) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	toggle_size := gui.FixedSize(u*2 - u/3)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.options.auto_wrap.text,
			},
			{
				Widget: &w.options.auto_wrap.toggle,
				Size:   toggle_size,
			},
			{
				Widget: &w.options.format.text,
			},
			{
				Widget: &w.options.format.toggle,
				Size:   toggle_size,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &w.content_type,
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rbh *response_body_header) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += rbh.options.auto_wrap.text.Measure(ctx, constraints).X + gap
		point.X += rbh.options.auto_wrap.toggle.Measure(ctx, constraints).X + gap
		point.X += rbh.options.format.text.Measure(ctx, constraints).X + gap
		point.X += rbh.options.format.toggle.Measure(ctx, constraints).X + gap
		point.X += rbh.content_type.Measure(ctx, constraints).X + gap
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u
	}
	return point
}

type response_body_widget struct {
	gui.DefaultWidget

	header response_body_header
	body   CommonWidgets.WidgetWithLazyLoading[*CommonWidgets.TextInputWithContextMenu]
}

func (w *response_body_widget) SetLazyLoad(lazy_load bool) {
	w.body.SetLazyLoad(lazy_load)
}

func (w *response_body_widget) LazyLoad() bool {
	return w.body.LazyLoad()
}

func (w *response_body_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.header)

	body := w.body.Widget()
	if w.header.options.auto_wrap.toggle.Value() {
		body.SetWrapMode(widget.WrapModeAnywhere)
	}
	body.SetMultiline(w.header.options.auto_wrap.toggle.Value())
	adder.AddWidget(&w.body)
	// make the view handle images and text.
	// Show content type not supported if content type is not jpg, png or text type.
	return nil
}

func (w *response_body_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
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

func (body *response_body_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
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

func (body *response_body_widget) SetBody(content string, content_type requests_handler.ContentType) {
	t, sub_t := content_type.Parse()
	if t == "text" || (t == "application" && sub_t == "json") || content_type == "" {
		body.body.Widget().ForceSetValue(content)
	} else {
		// TODO: handle images
		body.body.Widget().ForceSetValue("")
	}
}

func (body *response_body_widget) Body() string {
	return body.body.Widget().Value()
}

func (body *response_body_widget) ContentType() requests_handler.ContentType {
	return requests_handler.ContentType(body.header.content_type.Value())
}

func (body *response_body_widget) SetContentType(content_type requests_handler.ContentType) {
	body.header.content_type.SetValue(string(content_type))
}

func (body *response_body_widget) OnAutowrapToggle(fn func(ctx *gui.Context, value bool)) {
	body.header.options.auto_wrap.toggle.OnValueChanged(fn)
}

func (body *response_body_widget) OnFormatToggle(fn func(ctx *gui.Context, value bool)) {
	body.header.options.format.toggle.OnValueChanged(fn)
}

func (body *response_body_widget) SetAutowrap(autowrap bool) {
	body.header.options.auto_wrap.toggle.SetValue(autowrap)
}

func (body *response_body_widget) SetFormat(format bool) {
	body.header.options.format.toggle.SetValue(format)
}
