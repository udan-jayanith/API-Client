package CommonWidgets

import (
	"API-Client/basic"
	draw_color "API-Client/internal/draw"
	"API-Client/icons"
	attr "API-Client/pages/request/requests-handler/attributes"
	"image"
	"image/color"
	"slices"
	"strings"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type table_row_widget struct {
	gui.DefaultWidget

	table *attribute_table

	index                int
	checkbox             widget.Checkbox
	key_cell, value_cell EditableText
	row_delete_btn       icons.Icon
}

func (w *table_row_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 8)
}

func (w *table_row_widget) gap(ctx *gui.Context) int {
	return basic.Gap(ctx)
}

func (w *table_row_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !w.table.checkbox_disabled {
		adder.AddWidget(&w.checkbox)
	}

	w.key_cell.SetEditable(!w.table.key_not_editable)
	w.key_cell.SetAutoWrap(true)
	w.key_cell.SetEllipsisString("...")
	adder.AddWidget(&w.key_cell)

	w.value_cell.SetEllipsisString("...")
	w.value_cell.SetAutoWrap(true)
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
	gap, padding := w.gap(ctx), w.padding(ctx)

	b1 := widgetBounds.Bounds()
	b1.Min.X += padding.Start
	b1.Max.X -= padding.End
	b1.Min.Y += padding.Top
	b1.Max.Y -= padding.Bottom

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

	padding := row_widget.padding(ctx)
	constraints = gui.FixedWidthConstraints(point.X - (padding.Start + padding.End))
	y := max(row_widget.key_cell.Measure(ctx, constraints).Y, row_widget.value_cell.Measure(ctx, constraints).Y)
	point.Y = padding.Top + padding.Bottom
	point.Y += y
	return point
}

func (row_widget *table_row_widget) on_delete(fn func(index int)) {
	row_widget.row_delete_btn.OnClick(func() {
		fn(row_widget.index)
	})
}

type attribute_table struct {
	gui.DefaultWidget
	key_header, value_header widget.Text
	rows                     []*table_row_widget

	disable_auto_add                   bool
	checkbox_disabled, delete_disabled bool
	key_not_editable                   bool
	rwo_delete_fn                      func(index int)
}

func (at *attribute_table) push_row(row attr.AttrCheck) {
	row_widget := table_row_widget{}
	row_widget.index = len(at.rows)
	row_widget.table = at
	row_widget.checkbox.SetValue(row.Checked)
	row_widget.key_cell.SetValue(row.Key)
	row_widget.value_cell.SetValue(row.Value)
	at.rows = append(at.rows, &row_widget)
}

func (at *attribute_table) delete_row(index int) {
	at.rows = slices.Delete(at.rows, index, index+1)
	gui.RequestRebuild(at)
}

func (at *attribute_table) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	at.key_header.SetValue("Key")
	at.key_header.SetBold(true)
	at.key_header.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&at.key_header)

	at.value_header.SetValue("Value")
	at.value_header.SetBold(true)
	at.value_header.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&at.value_header)

	l := len(at.rows)
	if !at.disable_auto_add && (l == 0 || strings.TrimSpace(at.rows[l-1].key_cell.Value()) != "") {
		at.push_row(attr.AttrCheck{
			Checked: true,
		})
	}

	for i, _ := range at.rows {
		row_widget := at.rows[i]
		if !at.delete_disabled {
			row_widget.on_delete(at.delete_row)
		}
		row_widget.index = i
		adder.AddWidget(at.rows[i])
	}

	return nil
}

func (at *attribute_table) header_height(ctx *gui.Context) int {
	return widget.LineHeight(ctx) + (widget.UnitSize(ctx)/8)*2
}

func (at *attribute_table) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items:     make([]gui.LinearLayoutItem, 0, len(at.rows)),
	}

	u := widget.UnitSize(ctx)
	header_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 2,
		Padding:   basic.NewPadding(u / 8),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &at.key_header,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &at.value_header,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Layout: header_layout,
		Size:   gui.FixedSize(at.header_height(ctx)),
	})

	for i, _ := range at.rows {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: at.rows[i],
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (at *attribute_table) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx) * 6
	}

	point.Y = at.header_height(ctx)
	for i, _ := range at.rows {
		point.Y += at.rows[i].Measure(ctx, gui.FixedWidthConstraints(point.X)).Y
	}

	return point
}

func (at *attribute_table) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	header_heaight := at.header_height(ctx)
	b.Min.Y += header_heaight

	width := 1 * ctx.Scale()
	line_color := draw_color.ScaleAlpha(draw_color.Color(ctx.ColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)

	vector.StrokeLine(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Max.X), float32(b.Min.Y), float32(width), line_color, false)

	padding := widget.UnitSize(ctx) / 4
	middle := float32(b.Min.X + b.Dx()/2)
	vector.StrokeLine(dst, middle, float32(b.Min.Y)-float32(header_heaight)+float32(padding), middle, float32(b.Min.Y)-float32(padding), float32(width), line_color, false)

	for i, _ := range at.rows {
		if i%2 == 1 {
			b := b
			b.Max.Y = b.Min.Y + at.rows[i].Measure(ctx, gui.FixedWidthConstraints(b.Dx())).Y
			basicwidgetdraw.DrawRoundedRect(ctx, dst, b, color.Alpha16{0x1111}, basic.BorderRadius(ctx))
		}
		b.Min.Y += at.rows[i].Measure(ctx, gui.FixedWidthConstraints(b.Dx())).Y
	}
}

type AttributeTable struct {
	gui.DefaultWidget
	table gui.WidgetWithPadding[*attribute_table]

	panel widget.Panel
}

func (table *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	padding := basic.NewPadding(widget.UnitSize(ctx) / 3)
	padding.Top = 0
	table.table.SetPadding(padding)

	table.panel.SetContent(&table.table)
	table.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&table.panel)
	return nil
}

func (table *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&table.panel, widgetBounds.Bounds())
}

func (t *AttributeTable) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	// TODO: fix this
	size := image.Pt(12*widget.UnitSize(ctx), 6*widget.UnitSize(ctx))
	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	}

	return size
}

func (t *AttributeTable) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	r := widget.RoundedCornerRadius(ctx)
	background_clr := basicwidgetdraw.ControlColor(ctx.ColorMode(), ctx.IsEnabled(t))
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_clr, r)

	border_clr1, border_clr2 := basicwidgetdraw.BorderColors(ctx.ColorMode(), basicwidgetdraw.RoundedRectBorderTypeInset)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), border_clr1, border_clr2, r, 1, basicwidgetdraw.RoundedRectBorderTypeInset)
}

func (t *AttributeTable) SetRows(rows []attr.Attribute) {
	table_rows := t.table.Widget().rows
	if len(table_rows) > len(rows) {
		table_rows = table_rows[:len(rows)]
	} else if len(table_rows) != len(rows) {
		table_rows = make([]*table_row_widget, len(rows))
	}
	table := t.table.Widget()

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
	t.table.Widget().rows = table_rows
	gui.RequestRebuild(t)
}

func (t *AttributeTable) SetRowsCheck(rows []attr.AttrCheck) {
	table_rows := t.table.Widget().rows
	//TODO: BUG: optimizations doesn't work figure out what is happening.

	//if len(table_rows) > len(rows) {
	//table_rows = table_rows[:len(rows)]
	//} else if len(table_rows) != len(rows) {
	table_rows = make([]*table_row_widget, len(rows))
	//}
	table := t.table.Widget()

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
	t.table.Widget().rows = table_rows
	gui.RequestRebuild(t)
}

func (t *AttributeTable) RowsCheck() []attr.AttrCheck {
	table_rows := t.table.Widget().rows
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
	table_rows := t.table.Widget().rows
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
	t.table.Widget().checkbox_disabled = disable
	gui.RequestRebuild(t)
}

func (t *AttributeTable) DisableDelete(disable bool) {
	t.table.Widget().delete_disabled = disable
	gui.RequestRebuild(t)
}

func (t *AttributeTable) KeyEditable(editable bool) {
	t.table.Widget().key_not_editable = !editable
	gui.RequestRebuild(t)
}

func (t *AttributeTable) AutoAddRow(auto_add bool) {
	t.table.Widget().disable_auto_add = !auto_add
	gui.RequestRebuild(t)
}

func (t *AttributeTable) Count() int {
	return len(t.table.Widget().rows)
}

func (t *AttributeTable) PushRow(row attr.AttrCheck) {
	t.table.Widget().push_row(row)
	gui.RequestRebuild(t.table.Widget())
}
