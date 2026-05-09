package websocket_widget

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_header_widget struct {
	gui.DefaultWidget

	http_status_code, response_status widget.Text
	data_type                         widget.Select[struct{}]
	http_version                      widget.Text

	selected_data_type int
}

func (rw *response_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rw.http_status_code.SetValue("101")
	rw.http_status_code.SetVerticalAlign(widget.VerticalAlignMiddle)
	rw.http_status_code.SetTabular(true)
	adder.AddWidget(&rw.http_status_code)

	rw.response_status.SetValue("Switching Protocols")
	rw.response_status.SetVerticalAlign(widget.VerticalAlignMiddle)
	rw.response_status.SetTabular(true)
	adder.AddWidget(&rw.response_status)

	rw.data_type.SetItemsByStrings([]string{"Json", "Text", "Binary", "Other"})
	rw.data_type.OnItemSelected(func(context *gui.Context, index int) {
		rw.selected_data_type = index
	})
	rw.data_type.SelectItemByIndex(rw.selected_data_type)
	adder.AddWidget(&rw.data_type)

	rw.http_version.SetValue("HTTP v1.1")
	rw.http_version.SetVerticalAlign(widget.VerticalAlignMiddle)
	rw.http_version.SetTabular(true)
	adder.AddWidget(&rw.http_version)
	return nil
}

func (rw *response_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(u / 4),
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.http_status_code,
			},
			{
				Widget: &rw.response_status,
			},
			{
				Widget: &rw.data_type,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rw.http_version,
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rw *response_header_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	padding := u / 2

	point := image.Pt(padding+gap*4, padding)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += rw.http_status_code.Measure(ctx, gui.Constraints{}).X
		point.X += rw.response_status.Measure(ctx, gui.Constraints{}).X
		point.X += rw.data_type.Measure(ctx, gui.Constraints{}).X
		point.X += rw.http_version.Measure(ctx, gui.Constraints{}).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += widget.LineHeight(ctx)
	}

	return point
}
