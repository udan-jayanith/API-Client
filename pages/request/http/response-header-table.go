package http_widget

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	attr "Zbolt/pages/request/requests-handler/attributes"
	"image"

	gui "github.com/guigui-gui/guigui"
)

type response_header_table struct {
	gui.DefaultWidget

	search_bar CommonWidgets.SearchBar
	table      CommonWidgets.WidgetWithLazyLoading[*HttpHeaderTable]
}

func (table *response_header_table) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	headers_table := table.table.Widget()
	headers_table.DisableCheckbox(true)
	headers_table.DisableDelete(true)
	headers_table.KeyEditable(false)
	headers_table.ValueEditable(false)
	headers_table.AutoAddRow(false)
	adder.AddWidget(&table.table)

	adder.AddWidget(&table.search_bar)
	return nil
}

func (table *response_header_table) layout(ctx *gui.Context) gui.LinearLayout {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &table.search_bar,
			},
			{
				Widget: &table.table,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	return layout
}

func (table *response_header_table) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	table.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (table *response_header_table) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return table.layout(ctx).Measure(ctx, constraints)
}

func (table *response_header_table) SetSearchQuery(query string) {
	table.search_bar.SetQuery(query)
}

func (table *response_header_table) OnSearch(fn func(context *gui.Context, query string)) {
	table.search_bar.OnSearch(fn)
}

func (table *response_header_table) SetLazyLoad(lazy_load bool) {
	table.table.SetLazyLoad(lazy_load)
}

func (table *response_header_table) SetRows(rows []attr.Attribute) {
	table.table.Widget().SetRows(rows)
}
