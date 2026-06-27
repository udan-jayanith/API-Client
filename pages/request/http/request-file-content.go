package http_widget

import (
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_file_content struct {
	gui.DefaultWidget

	close_icon icons.Icon
	filename   CommonWidgets.TextWithTooltip
	select_btn widget.Button
}


func (*DefaultWidget) Build(context *Context, adder *ChildAdder) error {
	return nil
}

func (*DefaultWidget) Layout(context *Context, widgetBounds *WidgetBounds, layouter *ChildLayouter) {
}


func (d *DefaultWidget) Measure(context *Context, constraints Constraints) image.Point {
	var s image.Point
	if d.widgetState().root {
		s = context.app.bounds().Size()
	} else {
		s = image.Pt(int(144*context.Scale()), int(144*context.Scale()))
	}
	if w, ok := constraints.FixedWidth(); ok {
		s.X = w
	}
	if h, ok := constraints.FixedHeight(); ok {
		s.Y = h
	}
	return s
}
