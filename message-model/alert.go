package message_model

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type alert_widget_content struct {
	gui.DefaultWidget

	text      widget.Text
	panel     widget.Panel
	ok_button widget.Button
	on_result result_fn_type
}

func (w *alert_widget_content) SetMessage(message string) {
	w.text.SetValue(message)
}

func (w *alert_widget_content) OnResult(fn result_fn_type) {
	w.on_result = fn
	gui.RequestRebuild(w)
}

func (w *alert_widget_content) Bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	b := widgetBounds.Bounds()
	alert_box_point := w.Measure(ctx, gui.Constraints{})
	middle := image.Rectangle{
		Min: image.Point{
			X: (b.Min.X + b.Dx()/2) - alert_box_point.X/2,
			Y: (b.Min.Y + b.Dy()/2) - alert_box_point.Y/2,
		},
	}
	middle.Max = middle.Min.Add(alert_box_point)
	return middle
}

func (w *alert_widget_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.text.SetWrapMode(widget.WrapModeAnywhere)
	w.text.SetMultiline(true)
	w.text.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.text.SetVerticalAlign(widget.VerticalAlignTop)
	w.text.SetSelectable(true)

	w.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	w.panel.SetContent(&w.text)
	adder.AddWidget(&w.panel)

	w.ok_button.SetText("Ok")
	w.ok_button.SetType(widget.ButtonTypePrimary)
	if w.on_result != nil {
		w.ok_button.OnDown(func(context *gui.Context) {
			w.on_result(true, context)
		})
	}
	adder.AddWidget(&w.ok_button)
	return nil
}

func (w *alert_widget_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.panel,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &w.ok_button,
						},
					},
				},
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *alert_widget_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	point.X = u * 14
	point.Y = u * 8
	return point
}

type alert_widget struct {
	gui.DefaultWidget
	popup  widget.Popup
	widget alert_widget_content
}

func (w *alert_widget) SetMessage(message string) {
	if !w.popup.IsOpen() {
		w.popup.SetOpen(true)
	}
	w.widget.SetMessage(message)
}

func (w *alert_widget) OnResult(fn result_fn_type) {
	w.widget.OnResult(func(ok bool, ctx *gui.Context) {
		w.popup.SetOpen(false)
		if fn != nil {
			fn(ok, ctx)
		}
	})
}

func (w *alert_widget) Bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	return w.widget.Bounds(ctx, widgetBounds)
}

func (w *alert_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.popup.SetAnimated(true)
	w.popup.SetContent(&w.widget)
	adder.AddWidget(&w.popup)
	return nil
}

func (w *alert_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.popup, widgetBounds.Bounds())
}

func (w *alert_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return w.widget.Measure(ctx, constraints)
}
