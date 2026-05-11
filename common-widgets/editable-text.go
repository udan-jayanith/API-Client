package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type EditableText struct {
	gui.DefaultWidget
	not_editable bool
	text         widget.Text
	on_hover     func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)
	on_type      func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)
}

func (et *EditableText) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	et.text.SetEditable(!et.not_editable)
	adder.AddWidget(&et.text)
	return nil
}

func (et *EditableText) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&et.text, widgetBounds.Bounds())
}

func (et *EditableText) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return et.text.Measure(ctx, constraints)
}

func (et *EditableText) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if et.on_hover != nil && widgetBounds.IsHitAtCursor() {
		et.on_hover(ctx, widgetBounds)
	}
	return gui.HandleInputResult{}
}

func (et *EditableText) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if et.on_type != nil {
		et.on_type(ctx, widgetBounds)
	}
	return gui.HandleInputResult{}
}

func (et *EditableText) Widget() *widget.Text {
	return &et.text
}

func (et *EditableText) SetEditable(editable bool) {
	et.not_editable = !editable
}

func (et *EditableText) OnHover(fn func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)) {
	et.on_hover = fn
}

func (et *EditableText) OnType(fn func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)) {
	et.on_type = fn
}
