package http_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	attr "API-Client/pages/request/requests-handler/attributes"
	url_utils "API-Client/pages/request/requests-handler/url-utils"
	"image"
	"net/url"
	"strings"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type long_text_input_widget struct {
	CommonWidgets.TextInputWithContextMenu
}

func (w *long_text_input_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := w.TextInput.Measure(ctx, gui.Constraints{})
	point.X *= 3
	return point
}

type url_panel_content struct {
	gui.DefaultWidget

	form                 widget.Form
	scheme_text          widget.Text
	host_text, path_text CommonWidgets.TextWithInfoHint
	scheme               widget.Select[struct{}]
	host, path           long_text_input_widget
	// TODO: make the scheme a select to select between http and https

	query_header      widget.Text
	query_description CommonWidgets.Description
	query             CommonWidgets.AttributeTable

	hr1, hr2   CommonWidgets.HorizontalLine
	pseudo_url CommonWidgets.Description

	url_preview_header widget.Text
	url_preview        CommonWidgets.URLPreview

	table_updates_left                   bool
	url_preview_update_t, table_update_t time.Time
}

// update_query_table updates the query table based on the path input
func (w *url_panel_content) update_query_table() {
	pattern, _ := url_utils.ParsePattern(w.path.Value())
	merged_list := attr.MergeAttrList(w.query.Rows(), pattern.List)
	w.query.SetRows(merged_list)
}

// url returns the url for preview.
// For safety update_query_table must be run before this.
func (w *url_panel_content) url() string {
	pattern, _ := url_utils.ParsePattern(w.path.Value())
	for _, attr := range w.query.Rows() {
		pattern.Set(attr.Key, attr.Value)
	}

	u, _ := url.Parse(w.host.Value())
	w.init_scheme()
	selected_item, _ := w.scheme.SelectedItem()
	u.Scheme = strings.ToLower(selected_item.Text)
	u.Host = w.host.Value() // net/url parses the localhost:8080 by path
	u.Path = pattern.Path()
	return u.String()
}

func (w *url_panel_content) safe_url() string {
	w.update_query_table()
	w.table_update_t = time.Now()
	w.table_updates_left = false
	return w.url()
}

func (w *url_panel_content) pattern() (string, []attr.Attribute) {
	w.update_query_table()
	return w.path.Value(), w.query.Rows()
}

// clear clears the table rows for now
func (w *url_panel_content) clear() {
	w.query.SetRows([]attr.Attribute{})
}

func (w *url_panel_content) init_scheme() {
	if w.scheme.SelectedItemIndex() == -1 {
		w.scheme.SetItemsByStrings([]string{"HTTP", "HTTPS"})
		w.scheme.SelectItemByIndex(0)
	}
}

func (w *url_panel_content) set(shceme, host, path string, pattern []attr.Attribute) {
	w.init_scheme()
	if shceme == "https" || shceme == "HTTPS" {
		w.scheme.SelectItemByIndex(1)
	} else {
		w.scheme.SelectItemByIndex(0)
	}
	w.host.SetValue(host)
	w.path.SetValue(path)
	w.query.SetRows(pattern)
	w.update_query_table()
}

func (w *url_panel_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.init_scheme()

	w.scheme_text.SetValue("Scheme")
	w.host_text.SetValue("Host")

	// Warning: text hints are formatted
	w.host_text.SetHint(`"host" or "host:port"
	ex: www.google.com or localhost:8080`)
	w.path_text.SetHint(`URL path
	ex: "/en-US/docs/Web/API/URL/pathname"`)

	w.path_text.SetValue("Path")

	w.path.OnValueChanged(func(context *gui.Context, text string, committed bool) {
		w.table_updates_left = true
	})

	if time.Since(w.table_update_t).Seconds() >= 1 && w.table_updates_left {
		w.update_query_table()

		w.table_update_t = time.Now()
		w.table_updates_left = false
	}

	w.form.SetItems([]widget.FormItem{
		{
			PrimaryWidget:   &w.scheme_text,
			SecondaryWidget: &w.scheme,
		},
		{
			PrimaryWidget:   &w.host_text,
			SecondaryWidget: &w.host,
		},
		{
			PrimaryWidget:   &w.path_text,
			SecondaryWidget: &w.path,
		},
	})
	adder.AddWidget(&w.form)

	w.query_header.SetValue("Query")
	adder.AddWidget(&w.query_header)

	w.query_description.SetDescription("Attributes enclosed by {} in path.")
	adder.AddWidget(&w.query_description)

	w.query.AutoAddRow(false)
	w.query.DisableCheckbox(true)
	w.query.DisableDelete(true)
	w.query.KeyEditable(false)
	adder.AddWidget(&w.query)

	w.pseudo_url.SetDescription("The general form of the URL is:\n``[scheme:][//[host][/]path``")
	adder.AddWidget(&w.pseudo_url)
	adder.AddWidget(&w.hr1)

	adder.AddWidget(&w.hr2)

	w.url_preview_header.SetValue("URL preview")
	adder.AddWidget(&w.url_preview_header)

	if time.Now().Sub(w.url_preview_update_t).Seconds() >= 1 && !w.table_updates_left {
		w.url_preview.SetURL(w.url())
		w.url_preview_update_t = time.Now()
	}
	adder.AddWidget(&w.url_preview)
	return nil
}

func (w *url_panel_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	gap := gui.LinearLayoutItem{
		Size: gui.FixedSize(u / 3),
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(u / 2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.form,
			},
			gap,
			{
				Widget: &w.query_header,
			},
			{
				Widget: &w.query_description,
			},
			{
				Widget: &w.query,
				Size:   gui.FlexibleSize(1),
			},
			gap,
			{
				Widget: &w.hr1,
			},
			gap,
			{
				Widget: &w.pseudo_url,
			},
			gap,
			{
				Widget: &w.hr2,
			},
			gap,
			{
				Widget: &w.url_preview_header,
			},
			gap,
			{
				Widget: &w.url_preview,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *url_panel_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	point.X = u * 36
	point.Y = u * 20

	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
	}
	return point
}

type url_panel_widget struct {
	gui.DefaultWidget

	content url_panel_content
	panel   widget.Panel
}

// URL returns the url up to the path anything beyond the path of the url get excluded.
func (w *url_panel_widget) URL() string {
	return w.content.safe_url()
}

// Pattern if pattern exists length of query_list is greater then 0
func (w *url_panel_widget) Pattern() (pattern string, query_list []attr.Attribute) {
	return w.content.pattern()
}

func (w *url_panel_widget) Clear() {
	w.content.clear()
}

func (w *url_panel_widget) Set(shceme, host, path string, pattern []attr.Attribute) {
	w.content.set(shceme, host, path, pattern)
}

func (w *url_panel_widget) Build(context *gui.Context, adder *gui.ChildAdder) error {
	w.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	w.panel.SetContent(&w.content)
	adder.AddWidget(&w.panel)
	return nil
}

func (w *url_panel_widget) Layout(context *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.panel, widgetBounds.Bounds())
}

func (w *url_panel_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := w.content.Measure(ctx, gui.Constraints{})
	u := widget.UnitSize(ctx)

	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
		point.Y = u * 18
	}
	return point
}
