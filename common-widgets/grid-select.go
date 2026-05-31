package CommonWidgets

import (
	"Zbolt/basic"
	"Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type GridSelectItem[T any] struct {
	IconName string
	Title    string
	Value    T
}

type grid_select_item_widget[T any] struct {
	gui.DefaultWidget

	title widget.Text
	icon  icons.Icon
	value T
}

func (grid_select *grid_select_item_widget[T]) SetIcon(icon_name string) {
	grid_select.icon.SetIcon(icon_name)
}

func (grid_select *grid_select_item_widget[T]) SetTitle(title string) {
	grid_select.title.SetValue(title)
}

func (grid_select *grid_select_item_widget[T]) SetValue(val T) {
	grid_select.value = val
}

func (grid_select *grid_select_item_widget[T]) Value() T {
	return grid_select.value
}

func (grid_select *grid_select_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	grid_select.icon.SetSize(widget.UnitSize(ctx) * 2)
	adder.AddWidget(&grid_select.icon)

	grid_select.title.SetEllipsisString("...")
	grid_select.title.SetHorizontalAlign(widget.HorizontalAlignCenter)
	adder.AddWidget(&grid_select.title)
	return nil
}

func (grid_select *grid_select_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()

	icon_size := grid_select.icon.Measure(ctx, gui.Constraints{})
	icon_b := b
	icon_b.Min.Y += b.Dy()/2 - icon_size.Y/2
	icon_b.Max.Y = icon_b.Min.Y + icon_size.Y
	icon_b.Min.X += b.Dx()/2 - icon_size.X/2
	icon_b.Max.X = icon_b.Min.X + icon_size.X
	layouter.LayoutWidget(&grid_select.icon, icon_b)

	title_size := grid_select.title.Measure(ctx, gui.FixedWidthConstraints(b.Dx()))
	title_bounds := b
	title_bounds.Min.Y = icon_b.Max.Y + basic.Gap(ctx)
	title_bounds.Max.Y = title_bounds.Min.Y + title_size.Y
	title_bounds.Min.X += b.Dx()/2 - title_size.X/2
	title_bounds.Max.X = title_bounds.Min.X + title_size.X
	layouter.LayoutWidget(&grid_select.title, title_bounds)
}

func (grid_select *grid_select_item_widget[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, w)
	} else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(h, h)
	}

	s := widget.UnitSize(ctx) * 3
	return image.Pt(s, s)
}

type grid_select_content[T any] struct {
	gui.DefaultWidget

	grid_select_items           gui.WidgetSlice[*grid_select_item_widget[T]]
	grid_select_item_containers gui.WidgetSlice[*widget.Button]
}

func (grid_select *grid_select_content[T]) SetItems(items []GridSelectItem[T]) {
	grid_select.grid_select_items.SetLen(len(items))
	for i, item := range items {
		grid_select := grid_select.grid_select_items.At(i)
		grid_select.SetIcon(item.IconName)
		grid_select.SetTitle(item.Title)
		grid_select.SetValue(item.Value)
	}

	grid_select.grid_select_item_containers.SetLen(len(items))
	for i := range len(items) {
		container := grid_select.grid_select_item_containers.At(i)
		container.SetContent(grid_select.grid_select_items.At(i))
	}
}

func (grid_select *grid_select_content[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := range grid_select.grid_select_item_containers.Len() {
		adder.AddWidget(grid_select.grid_select_item_containers.At(i))
	}
	return nil
}

func (grid_select *grid_select_content[T]) item_size(ctx *gui.Context) int {
	return widget.UnitSize(ctx) * 4
}

func (grid_select *grid_select_content[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, grid_select.grid_select_items.Len()),
	}
	for i := range layout.Items {
		layout.Items[i].Widget = grid_select.grid_select_item_containers.At(i)
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
	/*
			b := widgetBounds.Bounds()
			item_size := grid_select.item_size(ctx)
			gap := basic.Gap(ctx)

			item_bounds := b
			item_bounds.Max.X = item_bounds.Min.X + item_size
			item_bounds.Max.Y = item_bounds.Min.Y + item_size

			var index int

		outer_loop:
			for {
				for {
					if grid_select.grid_select_items.Len() == index {
						break outer_loop
					} else if item_bounds.Max.X+item_size > b.Max.X {
						break
					}
					layouter.LayoutWidget(grid_select.grid_select_item_containers.At(index), item_bounds)
					item_bounds.Min.X = item_bounds.Max.X + gap
					item_bounds.Max.X = item_bounds.Min.X + item_size
					index++
				}

				item_bounds.Min.X = b.Min.X
				item_bounds.Max.X = b.Min.X + item_size
				item_bounds.Min.Y += gap
				item_bounds.Max.Y = item_bounds.Min.Y + item_size
			}
	*/
}

func (grid_select *grid_select_content[T]) calculate_items_per_row(ctx *gui.Context, w int) int {
	gap := basic.Gap(ctx)
	count := (w + gap) / (grid_select.item_size(ctx) + gap)
	return min(count, grid_select.grid_select_items.Len())
}

func (grid_select *grid_select_content[T]) calculate_number_of_rows(items_per_row int) int {
	if items_per_row == 0 {
		return 0
	}
	rows := grid_select.grid_select_items.Len() / items_per_row
	for rows*items_per_row < grid_select.grid_select_items.Len() {
		rows++
	}
	return rows
}

func (grid_select *grid_select_content[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if w, ok := constraints.FixedWidth(); ok {
		items_per_row := grid_select.calculate_items_per_row(ctx, w)
		rows := grid_select.calculate_number_of_rows(items_per_row)
		size := grid_select.item_size(ctx)
		return image.Pt(items_per_row*size, rows*size)
	} else if h, ok := constraints.FixedHeight(); ok {
		items_per_col := grid_select.calculate_items_per_row(ctx, h)
		no_of_cols := grid_select.calculate_number_of_rows(items_per_col)
		size := grid_select.item_size(ctx)
		return image.Pt(no_of_cols*size, items_per_col*size)
	} else {
		item_size := grid_select.item_size(ctx)
		items := grid_select.grid_select_items.Len()
		return image.Pt(items*item_size, item_size)
	}
}

type GridSelect[T any] struct {
	gui.DefaultWidget

	panel   widget.Panel
	content grid_select_content[T]
}

func (grid_select *GridSelect[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	grid_select.panel.SetContent(&grid_select.content)
	grid_select.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&grid_select.panel)
	return nil
}

func (grid_select *GridSelect[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&grid_select.panel, widgetBounds.Bounds())
}

func (grid_select *GridSelect[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return grid_select.Measure(ctx, constraints)
}

func (grid_select *GridSelect[T]) SetItems(items []GridSelectItem[T]) {
	grid_select.content.SetItems(items)
}
