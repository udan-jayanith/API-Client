package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/webp"
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

	text_content  CommonWidgets.WidgetWithLazyLoading[*widget.TextInput]
	image_content CommonWidgets.WidgetWithLazyLoading[*widget.Image]
}

func (w *response_body_widget) SetLazyLoad(lazy_load bool) {
	w.text_content.SetLazyLoad(lazy_load)
	w.image_content.SetLazyLoad(lazy_load)
}

func (w *response_body_widget) LazyLoad() bool {
	return w.text_content.LazyLoad() && w.image_content.LazyLoad()
}

func (w *response_body_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.header)

	if w.image_content.Widget().HasImage() {
		adder.AddWidget(&w.image_content)
	} else {
		body := w.text_content.Widget()
		if w.header.options.auto_wrap.toggle.Value() {
			body.SetWrapMode(widget.WrapModeAnywhere)
		} else {
			body.SetWrapMode(widget.WrapModeNone)
		}
		body.SetEditable(false)
		body.SetMultiline(true)
		adder.AddWidget(&w.text_content)
	}
	return nil
}

func (w *response_body_widget) layout(ctx *gui.Context) gui.LinearLayout {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.header,
			},
		},
	}

	if w.image_content.Widget().HasImage() {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Size:   gui.FlexibleSize(1),
			Widget: &w.image_content,
		})
	} else {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Size:   gui.FlexibleSize(1),
			Widget: &w.text_content,
		})
	}
	return layout
}

func (w *response_body_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	w.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (body *response_body_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return body.layout(ctx).Measure(ctx, constraints)
}

func (body *response_body_widget) set_unknow_data(msg string) {
	body.text_content.Widget().ForceSetValue("")
	body.image_content.Widget().SetImage(nil)
}

// TODO: handle the error
func (body *response_body_widget) SetBody(b *requests_handler.HTTP_Response_Body) {
	if b == nil || b.Content() == nil {
		body.text_content.Widget().ForceSetValue("")
		body.image_content.Widget().SetImage(nil)
		return
	}

	r := b.Content().NewReader()
	defer r.Close()

	t, sub_t := b.ContentType.Parse()
	if b.IsTexttual() {
		body.image_content.Widget().SetImage(nil)
		body.text_content.Widget().ReadValueFrom(r)
		r.Close()
	} else if t == "image" && (sub_t == "jpeg" || sub_t == "png" || sub_t == "webp") {
		body.text_content.Widget().ForceSetValue("")

		var img image.Image
		var err error
		switch sub_t {
		case "jpeg":
			img, err = jpeg.Decode(r)
		case "png":
			img, err = png.Decode(r)
		case "webp":
			img, err = webp.Decode(r)
		default:
			panic("Not handled")
		}

		if err != nil {
			body.set_unknow_data(err.Error())
		} else {
			ebit_img := ebiten.NewImageFromImage(img)
			body.image_content.Widget().SetImage(ebit_img)
		}
	} else {
		body.set_unknow_data("")
	}
	return
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
