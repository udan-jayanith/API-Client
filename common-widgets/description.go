package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Description struct {
	gui.DefaultWidget

	text_widget widget.Text
}

func (w *Description) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.text_widget.SetEllipsisString("...")
	w.text_widget.SetMultiline(true)
	w.text_widget.SetWrapMode(widget.WrapModeAnywhere)
	w.text_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	w.text_widget.SetOpacity(0.84)
	adder.AddWidget(&w.text_widget)
	return nil
}

func (w *Description) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.text_widget, widgetBounds.Bounds())
}

func (w *Description) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return w.text_widget.Measure(ctx, constraints)
}

func (w *Description) SetDescription(description string){
	w.text_widget.SetValue(description)
}

func (w *Description) Description() string{
	return w.text_widget.Value()
}