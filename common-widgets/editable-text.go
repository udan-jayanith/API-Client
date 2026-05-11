package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type EditableText struct {
	not_editable bool
	widget.Text
	on_hover func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)
	on_type  func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)
}

func (et *EditableText) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	et.Text.SetEditable(!et.not_editable)
	et.Text.Build(ctx, adder)
	return nil
}

func (et *EditableText) SetEditable(editable bool) {
	et.not_editable = !editable
}

/*
func (et *EditableText) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	// TODO: this dosen't work fix this later
	et.Text.HandlePointingInput(ctx, widgetBounds)
	if et.on_hover != nil && widgetBounds.IsHitAtCursor() {
		et.on_hover(ctx, widgetBounds)
	}
	return gui.HandleInputResult{}
}
 */
 
func (et *EditableText) OnHover(fn func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)) {
	et.on_hover = fn
}

func (et *EditableText) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	et.Text.HandleButtonInput(ctx, widgetBounds)
	if et.on_type != nil {
		et.on_type(ctx, widgetBounds)
	}
	return gui.HandleInputResult{}
}

func (et *EditableText) OnType(fn func(ctx *gui.Context, widget_bounds *gui.WidgetBounds)) {
	et.on_type = fn
}
