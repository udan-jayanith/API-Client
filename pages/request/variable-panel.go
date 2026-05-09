package request_page

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	attr "API-Client/pages/request/requests-handler/attributes"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type variable_panel_widget struct {
	gui.DefaultWidget

	public_table_header, private_table_header       widget.Text
	public_description, private_description         CommonWidgets.Description
	public_variables_table, private_variables_table CommonWidgets.AttributeTable
	line                                            CommonWidgets.HorizontalLine
}

func (w *variable_panel_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.public_table_header.SetValue("Public variables")
	w.public_table_header.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.public_table_header.SetVerticalAlign(widget.VerticalAlignBottom)
	adder.AddWidget(&w.public_table_header)

	w.private_table_header.SetValue("Private variables")
	w.private_table_header.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.private_table_header.SetVerticalAlign(widget.VerticalAlignBottom)
	adder.AddWidget(&w.private_table_header)

	w.public_description.SetDescription(`Public variables are visible to everyone.`)
	adder.AddWidget(&w.public_description)

	w.private_description.SetDescription(`Only variable names are visible to others.`)
	adder.AddWidget(&w.private_description)

	w.public_variables_table.DisableCheckbox(true)
	if w.public_variables_table.Count() == 0 {
		w.public_variables_table.SetRows([]attr.Attribute{{}, {Key: "api-key", Value: "gagj9a8gu2an9gih"}})
	}
	adder.AddWidget(&w.public_variables_table)

	w.private_variables_table.DisableCheckbox(true)
	adder.AddWidget(&w.private_variables_table)

	adder.AddWidget(&w.line)
	return nil
}

func (w *variable_panel_widget) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (w *variable_panel_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 2)
}

func (w *variable_panel_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := gui.LinearLayoutItem{
		Size: gui.FixedSize(w.gap(ctx)),
	}
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   w.padding(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.public_table_header,
				Size:   gui.FixedSize(widget.LineHeight(ctx)),
			},
			{
				Widget: &w.public_description,
			},
			gap,
			{
				Widget: &w.public_variables_table,
				Size:   gui.FlexibleSize(1),
			},
			{Size: gui.FixedSize(w.gap(ctx) * 2)},
			{
				Widget: &w.line,
			},
			gap,
			{
				Widget: &w.private_table_header,
			},
			{
				Widget: &w.private_description,
			},
			gap,
			{
				Widget: &w.private_variables_table,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *variable_panel_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = u * 26
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u * 18
	}
	return point
}
