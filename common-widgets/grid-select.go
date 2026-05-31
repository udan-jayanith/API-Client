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

	title gui.WidgetWithPadding[*widget.Text]
	icon  icons.Icon
	value T
}

func (grid_select *grid_select_item_widget[T]) SetIcon(icon_name string) {
	grid_select.icon.SetIcon(icon_name)
}

func (grid_select *grid_select_item_widget[T]) SetTitle(title string) {
	grid_select.title.Widget().SetValue(title)
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

	grid_select.title.Widget().SetEllipsisString("...")
	grid_select.title.Widget().SetHorizontalAlign(widget.HorizontalAlignCenter)
	grid_select.title.SetPadding(basic.NewPadding(0, basic.Gap(ctx)))
	adder.AddWidget(&grid_select.title)
	return nil
}

func (grid_select *grid_select_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()

	icon_size := grid_select.icon.Measure(ctx, gui.Constraints{})
	title_size := grid_select.title.Measure(ctx, gui.FixedWidthConstraints(b.Dx()))

	icon_b := b
	icon_b.Min.Y += b.Dy()/2 - (icon_size.Y/2 + title_size.Y/2)
	icon_b.Max.Y = icon_b.Min.Y + icon_size.Y
	icon_b.Min.X += b.Dx()/2 - icon_size.X/2
	icon_b.Max.X = icon_b.Min.X + icon_size.X
	layouter.LayoutWidget(&grid_select.icon, icon_b)

	title_bounds := b
	title_bounds.Min.Y = icon_b.Max.Y
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
	selected_index              int
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
	if len(items) > 0 {
		grid_select.on_select(0)
	}
}

func (grid_select *grid_select_content[T]) SelectedItemIndex() (int, bool) {
	return grid_select.selected_index, grid_select.grid_select_items.Len() > 0
}

func (grid_select *grid_select_content[T]) SelectedItem() (T, bool) {
	if grid_select.grid_select_items.Len() == 0 {
		var t T
		return t, false
	}
	return grid_select.grid_select_items.At(grid_select.selected_index).Value(), true
}

func (grid_select *grid_select_content[T]) SelectItemByIndex(index int) {
	grid_select.on_select(index)
}

func (grid_select *grid_select_content[T]) on_select(index int) {
	grid_select.grid_select_item_containers.At(grid_select.selected_index).SetType(widget.ButtonTypeNormal)
	grid_select.selected_index = index
	grid_select.grid_select_item_containers.At(index).SetType(widget.ButtonTypePrimary)
}

func (grid_select *grid_select_content[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := range grid_select.grid_select_item_containers.Len() {
		grid_select.grid_select_item_containers.At(i).OnDown(func(context *gui.Context) {
			grid_select.on_select(i)
		})
		adder.AddWidget(grid_select.grid_select_item_containers.At(i))
	}
	return nil
}

func (grid_select *grid_select_content[T]) item_size(ctx *gui.Context) int {
	return widget.UnitSize(ctx) * 4
}

func (grid_select *grid_select_content[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	item_size := grid_select.item_size(ctx)
	gap := basic.Gap(ctx)

	item_bounds := b
	item_bounds.Max.X = item_bounds.Min.X + item_size
	item_bounds.Max.Y = item_bounds.Min.Y + item_size

	var index int
	items_per_row := grid_select.calculate_items_per_row(ctx, b.Dx())
	rows := grid_select.calculate_number_of_rows(items_per_row)
	items := grid_select.grid_select_items.Len()
	for range rows {
		for range items_per_row {
			if index == items {
				break
			}
			layouter.LayoutWidget(grid_select.grid_select_item_containers.At(index), item_bounds)
			item_bounds.Min.X = item_bounds.Max.X + gap
			item_bounds.Max.X = item_bounds.Min.X + item_size
			index++
		}
		item_bounds.Min.X = b.Min.X
		item_bounds.Max.X = b.Min.X + item_size
		item_bounds.Min.Y = item_bounds.Max.Y + gap
		item_bounds.Max.Y = item_bounds.Min.Y + item_size
	}
}

func (grid_select *grid_select_content[T]) calculate_items_per_row(ctx *gui.Context, w int) int {
	gap := basic.Gap(ctx)
	item_size := grid_select.item_size(ctx)
	count := (w + gap) / (item_size + gap)
	for (count*(item_size+gap))+(item_size+gap) < w {
		count++
	}
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

func (grid_select *grid_select_content[T]) RecommendedSize(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if w, ok := constraints.FixedWidth(); ok {
		items_per_row := grid_select.calculate_items_per_row(ctx, w)
		rows := grid_select.calculate_number_of_rows(items_per_row)
		size := grid_select.item_size(ctx)
		return image.Pt(items_per_row*size, rows*size+(rows-1)*basic.Gap(ctx))
	} else {
		item_size := grid_select.item_size(ctx)
		items := grid_select.grid_select_items.Len()
		return image.Pt((items*item_size)+((items-1)*basic.Gap(ctx)), item_size)
	}
}

func (grid_select *grid_select_content[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	size := grid_select.RecommendedSize(ctx, constraints)
	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	}
	return size
}

/*
func (grid_select *grid_select_content[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		ctx.SetFocused(grid_select, true)
	} else {
		ctx.SetFocused(grid_select, false)
	}
	return gui.HandleInputResult{}
}

func (grid_select *grid_select_content[T]) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !ctx.IsFocused(grid_select) {
		return gui.HandleInputResult{}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if grid_select.selected_index == 0 {
			index := grid_select.grid_select_items.Len() - 1
			grid_select.on_select(index)
		} else {
			grid_select.on_select(grid_select.selected_index - 1)
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		println("Pressed")
		if grid_select.selected_index == grid_select.grid_select_items.Len()-1 {
			grid_select.on_select(0)
		} else {
			grid_select.on_select(grid_select.selected_index + 1)
		}
	}

	return gui.HandleInputResult{}
}
*/

type GridSelect[T any] struct {
	gui.DefaultWidget

	panel   widget.Panel
	content grid_select_content[T]
}

func (grid_select *GridSelect[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	//grid_select.content.SetPadding(basic.NewPadding(basic.BorderRadius(ctx)))
	grid_select.panel.SetContent(&grid_select.content)
	grid_select.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&grid_select.panel)
	return nil
}

func (grid_select *GridSelect[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&grid_select.panel, widgetBounds.Bounds())
}

func (grid_select *GridSelect[T]) RecommendedSize(ctx *gui.Context) image.Point {
	size := grid_select.content.RecommendedSize(ctx, gui.Constraints{})
	return size
}

func (grid_select *GridSelect[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return grid_select.panel.Measure(ctx, constraints)
}

func (grid_select *GridSelect[T]) SetItems(items []GridSelectItem[T]) {
	grid_select.content.SetItems(items)
}

func (grid_select *GridSelect[T]) SelectedItem() (T, bool) {
	return grid_select.content.SelectedItem()
}

func (grid_select *GridSelect[T]) SelectedItemIndex() (int, bool) {
	return grid_select.content.SelectedItemIndex()
}

func (grid_select *GridSelect[T]) SelectItemByIndex(index int) {
	grid_select.content.SelectItemByIndex(index)
}
