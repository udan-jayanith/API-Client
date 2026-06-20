package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	"fmt"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_body_header struct {
	gui.DefaultWidget

	content_type        requests_handler.ContentType
	content_type_widget CommonWidgets.TextWithTooltip

	options struct {
		disable   bool
		auto_wrap struct {
			text   widget.Text
			toggle widget.Toggle
		}
		format struct {
			text   widget.Text
			toggle widget.Toggle
		}
	}
}

func (w *response_body_header) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.content_type_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	w.content_type_widget.SetEllipsisString("...")
	adder.AddWidget(&w.content_type_widget)

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

	ctx.SetEnabled(&w.options.auto_wrap.text, !w.options.disable)
	ctx.SetEnabled(&w.options.auto_wrap.toggle, !w.options.disable)
	ctx.SetEnabled(&w.options.format.text, !w.options.disable)
	ctx.SetEnabled(&w.options.format.toggle, !w.options.disable)
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
				Widget: &w.content_type_widget,
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
		point.X += rbh.content_type_widget.Measure(ctx, constraints).X + gap
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u
	}
	return point
}

func (rbh *response_body_header) SetContentType(content_type requests_handler.ContentType) {
	rbh.content_type = content_type
	ex := content_type.Extension()
	rbh.content_type_widget.SetValue(ex)
	rbh.content_type_widget.SetTooltip(fmt.Sprintf("Content Type: %s", content_type))
}

func (rbh *response_body_header) ContentType() requests_handler.ContentType {
	return rbh.content_type
}

func (rbh *response_body_header) EnableToggles(enable bool) {
	rbh.options.disable = !enable
}

type response_body_widget struct {
	gui.DefaultWidget

	header response_body_header

	text_content  widget.TextInput
	image_content widget.Image
	body          CommonWidgets.WidgetWithLazyLoading[*CommonWidgets.TextInputWithContextMenu]
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
	} else {
		body.SetWrapMode(widget.WrapModeNone)
	}
	body.SetEditable(false)
	body.SetMultiline(true)
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

func (body *response_body_widget) SetBody(b *requests_handler.HTTP_Response_Body) {
	if b == nil || b.Content() == nil {
		body.body.Widget().ForceSetValue("")
		return
	}

	t, sub_t := b.ContentType.Parse()
	if t == "text" || (t == "application" && sub_t == "json") || b.ContentType == "" {
		body.image_content.SetImage(nil)
		r := b.Content().NewReader()
		body.body.Widget().ReadValueFrom(r)
		r.Close()
	} else if t == "image" && (sub_t == "jpeg" || sub_t == "png" || sub_t == "webp") {
		body.body.Widget().ForceSetValue("")
	}
}

func (body *response_body_widget) Body() string {
	return body.body.Widget().Value()
}

func (body *response_body_widget) ContentType() requests_handler.ContentType {
	return body.header.ContentType()
}

func (body *response_body_widget) SetContentType(content_type requests_handler.ContentType) {
	t, sub_t := content_type.Parse()
	body.header.EnableToggles(t == "text" || (t == "application" && sub_t == "json") || content_type == "")
	body.header.SetContentType(content_type)
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
