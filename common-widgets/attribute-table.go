package CommonWidgets

import (
	"Zbolt/basic"
	"Zbolt/icons"
	attr "Zbolt/pages/request/requests-handler/attributes"
	"image"
	"slices"
	"strings"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type table_row_widget struct {
	gui.DefaultWidget

	table *AttributeTable

	index                int
	checkbox             widget.Checkbox
	key_cell, value_cell EditableText
	row_delete_btn       icons.Icon
}

func (w *table_row_widget) gap(ctx *gui.Context) int {
	return basic.Gap(ctx)
}

func (w *table_row_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 8)
}

func (w *table_row_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !w.table.checkbox_disabled {
		adder.AddWidget(&w.checkbox)
	}

	w.key_cell.SetEditable(!w.table.key_not_editable)
	key_cell := &w.key_cell
	key_cell.SetWrapMode(widget.WrapModeAnywhere)
	if w.table.on_type != nil {
		w.key_cell.OnType(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_type(ctx, "key", &w.key_cell, widget_bounds)
		})
	}
	if w.table.on_hover != nil {
		w.key_cell.OnHover(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_hover(ctx, "key", &w.key_cell, widget_bounds)
		})
	}
	adder.AddWidget(&w.key_cell)

	w.value_cell.SetEditable(!w.table.value_not_editable)
	value_cell := &w.value_cell
	value_cell.SetWrapMode(widget.WrapModeAnywhere)
	if w.table.on_type != nil {
		w.value_cell.OnType(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_type(ctx, "value", &w.value_cell, widget_bounds)
		})
	}
	if w.table.on_hover != nil {
		w.value_cell.OnHover(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_hover(ctx, "value", &w.value_cell, widget_bounds)
		})
	}
	adder.AddWidget(&w.value_cell)

	if !w.table.delete_disabled {
		l := widget.LineHeight(ctx)
		w.row_delete_btn.SetSize(l - (l / 6))
		w.row_delete_btn.SetIcon("delete")
		adder.AddWidget(&w.row_delete_btn)
	}
	return nil
}

func (w *table_row_widget) layout(ctx *gui.Context) gui.LinearLayout {
	layout := gui.LinearLayout{
		Padding: w.padding(ctx),
		Gap:     w.gap(ctx),
		Items:   make([]gui.LinearLayoutItem, 0, 4),
	}

	gap := w.gap(ctx)/6
	if !w.table.checkbox_disabled {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Layout: gui.LinearLayout{
				Gap: gap,
				Items: []gui.LinearLayoutItem{
					{
						Widget: &w.checkbox,
						Size:   gui.FixedSize(widget.UnitSize(ctx)),
					},
					{
						Widget: &w.key_cell,
						Size:   gui.FlexibleSize(1),
					},
				},
			},
			Size: gui.FlexibleSize(1),
		})
	} else {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &w.key_cell,
			Size:   gui.FlexibleSize(1),
		})
	}

	if !w.table.delete_disabled {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Layout: gui.LinearLayout{
				Gap: gap,
				Items: []gui.LinearLayoutItem{
					{
						Widget: &w.value_cell,
						Size:   gui.FlexibleSize(1),
					},
					{
						Widget: &w.row_delete_btn,
						Size:   gui.FixedSize(w.row_delete_btn.Measure(ctx, gui.Constraints{}).X),
					},
				},
			},
			Size: gui.FlexibleSize(1),
		})
	} else {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &w.value_cell,
			Size:   gui.FlexibleSize(1),
		})
	}
	return layout
}

func (w *table_row_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	w.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (row_widget *table_row_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return row_widget.layout(ctx).Measure(ctx, constraints)
}

func (row_widget *table_row_widget) on_delete(fn func(index int)) {
	row_widget.row_delete_btn.OnClick(func() {
		fn(row_widget.index)
	})
}

type table_header struct {
	gui.DefaultWidget

	key, value widget.Text
	vr         VerticalLine
}

func (t *table_header) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	t.key.SetValue("Key")
	t.value.SetValue("Value")
	adder.AddWidget(&t.key)
	adder.AddWidget(&t.vr)
	adder.AddWidget(&t.value)
	return nil
}

func (t *table_header) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	padding := widget.UnitSize(ctx) / 8
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Padding:   basic.NewPadding(0, padding),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &t.key,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &t.vr,
			},
			{
				Widget: &t.value,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (t *table_header) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var s image.Point
	if w, ok := constraints.FixedWidth(); ok {
		s.X = w
	} else {
		s.X = widget.UnitSize(ctx) * 6
	}

	s.Y = widget.LineHeight(ctx)
	return s
}

type AttributeTable struct {
	gui.DefaultWidget
	header table_header
	hr     HorizontalLine
	rows   []*table_row_widget
	list   widget.List[struct{}]

	disable_auto_add                     bool
	checkbox_disabled, delete_disabled   bool
	key_not_editable, value_not_editable bool
	rwo_delete_fn                        func(index int)
	on_hover                             func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)
	on_type                              func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)
	on_scroll                            func(offset_x, height int)
}

func (at *AttributeTable) push_row(row attr.AttrCheck) {
	row_widget := table_row_widget{}
	row_widget.index = len(at.rows)
	row_widget.table = at
	row_widget.checkbox.SetValue(row.Checked)
	row_widget.key_cell.SetValue(row.Key)
	row_widget.value_cell.SetValue(row.Value)
	at.rows = append(at.rows, &row_widget)
}

func (at *AttributeTable) delete_row(index int) {
	at.rows = slices.Delete(at.rows, index, index+1)
	gui.RequestRebuild(at)
}

func (at *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	l := len(at.rows)
	if !at.disable_auto_add && (l == 0 || strings.TrimSpace(at.rows[l-1].key_cell.Value()) != "") {
		at.push_row(attr.AttrCheck{
			Checked: true,
		})
	}

	list_items := make([]widget.ListItem[struct{}], len(at.rows))
	for i, _ := range at.rows {
		row_widget := at.rows[i]
		if !at.delete_disabled {
			row_widget.on_delete(at.delete_row)
		}
		row_widget.index = i
		list_items[i].Content = row_widget
	}
	list_items = append([]widget.ListItem[struct{}]{
		{
			Content:      &at.header,
			Unselectable: true,
			Header:       true,
		},
		{
			Content:      &at.hr,
			Unselectable: true,
			Padding:      basic.NewPadding(widget.UnitSize(ctx)/8, 0),
		},
	}, list_items...)
	at.list.SetItems(list_items)

	adder.AddWidget(&at.list)

	return nil
}

func (at *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&at.list, widgetBounds.Bounds())
}

func (at *AttributeTable) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return at.list.Measure(ctx, constraints)
}

func (t *AttributeTable) SetRows(rows []attr.Attribute) {
	table_rows := t.rows
	if len(table_rows) > len(rows) {
		table_rows = table_rows[:len(rows)]
	} else if len(table_rows) != len(rows) {
		table_rows = make([]*table_row_widget, len(rows))
	}
	table := t

	for i, row := range rows {
		if table_rows[i] == nil {
			table_rows[i] = &table_row_widget{}
		}
		table_row := table_rows[i]
		table_row.table = table
		table_row.index = i

		table_row.key_cell.SetValue(row.Key)
		table_row.value_cell.SetValue(row.Value)
	}
	t.rows = table_rows
}

func (t *AttributeTable) SetRowsCheck(rows []attr.AttrCheck) {
	table_rows := t.rows
	if len(table_rows) > len(rows) {
		table_rows = table_rows[:len(rows)]
	} else if len(table_rows) != len(rows) {
		table_rows = make([]*table_row_widget, len(rows))
	}
	table := t

	for i, row := range rows {
		if table_rows[i] == nil {
			table_rows[i] = &table_row_widget{}
		}

		table_row := table_rows[i]
		table_row.table = table
		table_row.index = i

		table_row.checkbox.SetValue(row.Checked)
		table_row.key_cell.SetValue(row.Key)
		table_row.value_cell.SetValue(row.Value)
	}
	t.rows = table_rows
}

func (t *AttributeTable) RowsCheck() []attr.AttrCheck {
	table_rows := t.rows
	rows := make([]attr.AttrCheck, 0, len(table_rows))

	for _, table_row := range table_rows {
		if strings.TrimSpace(table_row.key_cell.Value()) == "" && strings.TrimSpace(table_row.value_cell.Value()) == "" {
			continue
		}
		rows = append(rows, attr.AttrCheck{
			Key:     table_row.key_cell.Value(),
			Value:   table_row.value_cell.Value(),
			Checked: table_row.checkbox.Value(),
		})
	}

	return rows
}

func (t *AttributeTable) Rows() []attr.Attribute {
	table_rows := t.rows
	rows := make([]attr.Attribute, 0, len(table_rows))

	for _, table_row := range table_rows {
		if strings.TrimSpace(table_row.key_cell.Value()) == "" && strings.TrimSpace(table_row.value_cell.Value()) == "" {
			continue
		}
		rows = append(rows, attr.Attribute{
			Key:   table_row.key_cell.Value(),
			Value: table_row.value_cell.Value(),
		})
	}

	return rows
}

func (t *AttributeTable) DisableCheckbox(disable bool) {
	t.checkbox_disabled = disable
}

func (t *AttributeTable) DisableDelete(disable bool) {
	t.delete_disabled = disable
}

func (t *AttributeTable) KeyEditable(editable bool) {
	t.key_not_editable = !editable
}

func (t *AttributeTable) ValueEditable(editable bool) {
	t.value_not_editable = !editable
}

func (t *AttributeTable) AutoAddRow(auto_add bool) {
	t.disable_auto_add = !auto_add
}

func (t *AttributeTable) Count() int {
	return len(t.rows)
}

func (t *AttributeTable) PushRow(row attr.AttrCheck) {
	t.push_row(row)
}

func (t *AttributeTable) OnHover(fn func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)) {
	t.on_hover = fn
}

func (t *AttributeTable) OnType(fn func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)) {
	t.on_type = fn
}
