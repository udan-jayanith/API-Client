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

func (w *table_row_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !w.table.checkbox_disabled {
		adder.AddWidget(&w.checkbox)
	}

	w.key_cell.SetEditable(!w.table.key_not_editable)
	w.key_cell.SetWrapMode(widget.WrapModeAnywhere)
	w.key_cell.SetEllipsisString("...")
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
	w.value_cell.SetEllipsisString("...")
	w.value_cell.SetWrapMode(widget.WrapModeAnywhere)
	if w.table.on_type != nil {
		w.value_cell.OnType(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_type(ctx, "value", &w.key_cell, widget_bounds)
		})
	}
	if w.table.on_hover != nil {
		w.value_cell.OnHover(func(ctx *gui.Context, widget_bounds *gui.WidgetBounds) {
			w.table.on_hover(ctx, "value", &w.key_cell, widget_bounds)
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

func (w *table_row_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := w.gap(ctx)

	b1 := widgetBounds.Bounds()

	if !w.table.checkbox_disabled {
		size := w.checkbox.Measure(ctx, gui.Constraints{})
		checkbox_bounds := b1
		checkbox_bounds.Max.X = checkbox_bounds.Min.X + size.X
		checkbox_bounds.Max.Y = checkbox_bounds.Min.Y + size.Y
		b1.Min.X += gap + size.X
		layouter.LayoutWidget(&w.checkbox, checkbox_bounds)
	}

	if !w.table.delete_disabled {
		size := w.row_delete_btn.Measure(ctx, gui.Constraints{})
		btn_bounds := b1
		btn_bounds.Min.X = btn_bounds.Max.X - size.X
		btn_bounds.Min.Y += gap / 4
		btn_bounds.Max.Y = btn_bounds.Min.Y + size.Y
		b1.Max.X -= (gap + size.X)
		layouter.LayoutWidget(&w.row_delete_btn, btn_bounds)
	}

	b2 := widgetBounds.Bounds()
	middle := b2.Min.X + b2.Dx()/2

	key_bounds := b1
	key_bounds.Max.X = middle - gap
	layouter.LayoutWidget(&w.key_cell, key_bounds)

	val_bounds := b1
	val_bounds.Min.X = middle + gap
	layouter.LayoutWidget(&w.value_cell, val_bounds)
}

func (row_widget *table_row_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx) * 6
	}

	constraints = gui.FixedWidthConstraints(point.X)
	y := max(row_widget.key_cell.Measure(ctx, constraints).Y, row_widget.value_cell.Measure(ctx, constraints).Y)
	point.Y += y
	return point
}

func (row_widget *table_row_widget) on_delete(fn func(index int)) {
	row_widget.row_delete_btn.OnClick(func() {
		fn(row_widget.index)
	})
}

type AttributeTable struct {
	gui.DefaultWidget
	key_header, value_header widget.Text
	rows                     []*table_row_widget
	list                     widget.List[struct{}]

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
	padding := basic.NewPadding(widget.UnitSize(ctx) / 8)
	for i, _ := range at.rows {
		row_widget := at.rows[i]
		if !at.delete_disabled {
			row_widget.on_delete(at.delete_row)
		}
		row_widget.index = i
		//adder.AddWidget(at.rows[i])
		list_items[i].Content = row_widget
		list_items[i].Padding = padding
	}
	list_items = append([]widget.ListItem[struct{}]{
		{
			Header: true,
			Border: true,
			Text:   "Hello world",
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
	//TODO: BUG: optimizations doesn't work figure out what is happening.

	//if len(table_rows) > len(rows) {
	//table_rows = table_rows[:len(rows)]
	//} else if len(table_rows) != len(rows) {
	table_rows = make([]*table_row_widget, len(rows))
	//}
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
		if strings.TrimSpace(table_row.key_cell.Value()) == "" {
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
		if strings.TrimSpace(table_row.key_cell.Value()) == "" {
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
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) DisableDelete(disable bool) {
	t.delete_disabled = disable
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) KeyEditable(editable bool) {
	t.key_not_editable = !editable
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) ValueEditable(editable bool) {
	t.value_not_editable = !editable
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) AutoAddRow(auto_add bool) {
	t.disable_auto_add = !auto_add
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) Count() int {
	return len(t.rows)
}

func (t *AttributeTable) PushRow(row attr.AttrCheck) {
	t.push_row(row)
	gui.RequestRebuild(&t.list)
}

func (t *AttributeTable) OnHover(fn func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)) {
	t.on_hover = fn
}

func (t *AttributeTable) OnType(fn func(ctx *gui.Context, t string, widget *EditableText, widget_bounds *gui.WidgetBounds)) {
	t.on_type = fn
}
