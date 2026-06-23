package http_widget

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	attr "Zbolt/pages/request/requests-handler/attributes"
	"image"
	"sort"
	"strings"

	gui "github.com/guigui-gui/guigui"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type response_header_table struct {
	gui.DefaultWidget

	search_bar CommonWidgets.SearchBar
	table      CommonWidgets.WidgetWithLazyLoading[*HttpHeaderTable]

	search_query         string
	rows, search_results []attr.Attribute
}

func (table *response_header_table) search(query string) {
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)
	if query == "" {
		table.search_query = ""
		table.search_results = nil
		table.table.Widget().SetRows(table.rows)
		return
	}

	var rows []attr.Attribute
	if table.search_query != "" && strings.Contains(query, table.search_query) {
		rows = table.search_results
	} else {
		rows = table.rows
	}

	results := make([]attr.Attribute, 0, len(rows))
	for _, row := range rows {
		if fuzzy.Match(query, strings.ToLower(row.Key)) {
			results = append(results, row)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return fuzzy.RankMatch(query, strings.ToLower(results[i].Key)) < fuzzy.RankMatch(query, strings.ToLower(results[j].Key))
	})

	table.search_query = query
	table.search_results = results
	table.table.Widget().SetRows(results)
}

func (table *response_header_table) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	headers_table := table.table.Widget()
	headers_table.DisableCheckbox(true)
	headers_table.DisableDelete(true)
	headers_table.KeyEditable(false)
	headers_table.ValueEditable(false)
	headers_table.AutoAddRow(false)
	adder.AddWidget(&table.table)

	table.search_bar.OnSearch(func(context *gui.Context, query string) {
		table.search(query)
	})
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
	table.search(query)
}

func (table *response_header_table) SearchQuery() string {
	return table.search_bar.Query()
}

func (table *response_header_table) SetLazyLoad(lazy_load bool) {
	table.table.SetLazyLoad(lazy_load)
}

func (table *response_header_table) SetRows(rows []attr.Attribute) {
	table.search_results = nil
	table.search_query = ""
	table.rows = rows
	table.table.Widget().SetRows(rows)
	table.SetSearchQuery("")
}
