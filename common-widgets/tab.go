package CommonWidgets

import (
	"image"

	"Zbolt/basic"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type TabItemContainer struct {
	Index int
	Item  TabItem
}

type tab_bar struct {
	gui.DefaultWidget
	tab_items gui.WidgetSlice[*tab_item]
	closable  bool

	holding struct {
		is_holding               bool
		tab_item_index           int
		relative_cursor_position int
		closest_index            int // closest_index is the closest tab item that is closest to the holding tab item
	}
	selected_item_index int

	listeners struct {
		on_select func(from TabItemContainer, to TabItemContainer, by_user bool)
		on_swap   func(from TabItemContainer, to TabItemContainer)
		on_close  func(closed TabItemContainer)
	}
}

func (tab *tab_bar) item_count() int {
	return tab.tab_items.Len()
}

func (tab *tab_bar) on_select(index int, tab_item TabItem, by_user bool) {
	if tab.listeners.on_select != nil {
		from := TabItemContainer{
			Index: tab.selected_item_index,
			Item:  tab.tab_items.At(tab.selected_item_index).tab_item,
		}
		to := TabItemContainer{
			Index: index,
			Item:  tab_item,
		}
		tab.listeners.on_select(from, to, by_user)
	}
	tab.selected_item_index = index
	gui.RequestRedraw(tab)
}

func (tab *tab_bar) on_holding(index int, relative_cursor_x int) {
	tab.holding.is_holding = true
	tab.holding.tab_item_index = index
	tab.holding.relative_cursor_position = relative_cursor_x
	gui.RequestRebuild(tab)
}

func (tab *tab_bar) on_mouse_up(index int) {
	if !tab.holding.is_holding {
		return
	}
	tab.holding.is_holding = false

	if tab.tab_items.Len() > 0 {
		from := TabItemContainer{
			Index: index,
			Item:  tab.tab_items.At(index).tab_item,
		}
		to := TabItemContainer{
			Index: tab.holding.closest_index,
			Item:  tab.tab_items.At(tab.holding.closest_index).tab_item,
		}
		tab.selected_item_index = tab.holding.closest_index
		if tab.listeners.on_swap != nil {
			tab.listeners.on_swap(from, to)
		}
	}
	gui.RequestRedraw(tab)
}

func (tab *tab_bar) on_close(index int, item TabItem) {
	if tab.listeners.on_close == nil {
		return
	}
	tab.listeners.on_close(TabItemContainer{
		Index: index,
		Item:  item,
	})
	gui.RequestRebuild(tab)
}

func (tab *tab_bar) update_tab_items(tab_items []TabItem) {
	tab.tab_items.SetLen(len(tab_items))
	for i := range tab.tab_items.Len() {
		widget := tab.tab_items.At(i)
		widget.tab_item = tab_items[i]
		widget.index = i
		widget.tabs_container = tab
	}
}

func (tab *tab_bar) set_tab_items(tab_items []TabItem) {
	if tab.tab_items.Len() == len(tab_items) {
		tab.update_tab_items(tab_items)
		return
	}

	tab.selected_item_index = 0
	if len(tab_items) <= tab.selected_item_index {
		tab.selected_item_index = 0
	}
	tab.update_tab_items(tab_items)
	if len(tab_items) > 0 {
		tab.on_select(0, tab_items[0], false)
	}
}

func (tab *tab_bar) select_tab(index int, by_user bool) {
	if index >= tab.tab_items.Len() {
		panic("Invalid index")
	}
	tab.on_select(index, tab.tab_items.At(index).tab_item, by_user)
}

func (tab *tab_bar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := range tab.tab_items.Len() {
		widget := tab.tab_items.At(i)
		if tab.holding.is_holding && tab.holding.tab_item_index == widget.index {
			continue
		}
		adder.AddWidget(widget)
	}

	if tab.holding.is_holding {
		adder.AddWidget(tab.tab_items.At(tab.holding.tab_item_index))
	}
	return nil
}

func (tab *tab_bar) holding_item_bounds(ctx *gui.Context, b image.Rectangle) image.Rectangle {
	cursor_axis, _ := ebiten.CursorPosition()
	tab_item_bounds := image.Rectangle{
		Min: image.Point{
			X: cursor_axis - tab.holding.relative_cursor_position,
			Y: b.Min.Y,
		},
	}

	holding_tab_item := tab.tab_items.At(tab.holding.tab_item_index)
	tab_item_bounds.Max.X = tab_item_bounds.Min.X + holding_tab_item.Measure(ctx, gui.Constraints{}).X
	tab_item_bounds.Max.Y = b.Max.Y
	return tab_item_bounds
}

func (tab *tab_bar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, tab.tab_items.Len()),
	}

	if tab.holding.is_holding {
		b := widgetBounds.Bounds()
		layouter.LayoutWidget(tab.tab_items.At(tab.holding.tab_item_index), tab.holding_item_bounds(ctx, b))
	}

	for i := range tab.tab_items.Len() {
		widget := tab.tab_items.At(i)
		if tab.holding.is_holding && tab.holding.tab_item_index == widget.index {
			holding_tab_item := tab.tab_items.At(tab.holding.tab_item_index)
			w := holding_tab_item.Measure(ctx, gui.Constraints{}).X
			layout.Items = append(layout.Items, gui.LinearLayoutItem{
				Size: gui.FixedSize(w),
			})
			continue
		}

		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: widget,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *tab_bar) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !tab.holding.is_holding || tab.tab_items.Len() == 0 {
		return gui.HandleInputResult{}
	}

	b := widgetBounds.Bounds()
	var closest_item_index int
	cursor_x, _ := ebiten.CursorPosition()

	for i := range tab.tab_items.Len() {
		w := tab.tab_items.At(i).Measure(ctx, gui.Constraints{}).X
		if b.Min.X <= cursor_x && cursor_x <= b.Min.X+w {
			closest_item_index = i
			break
		}
		b.Min.X += w
	}

	tab.holding.closest_index = closest_item_index
	return gui.HandleInputResult{}
}

func (tab *tab_bar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	for i := range tab.tab_items.Len() {
		measurement := tab.tab_items.At(i).Measure(ctx, constraints)
		point.X += measurement.X
	}

	u := widget.UnitSize(ctx)
	point.Y = u

	return point
}

type Tab struct {
	gui.DefaultWidget
	panel         widget.Panel
	tab_container tab_bar
}

func (tab *Tab) SetTabItems(items []TabItem) {
	tab.tab_container.set_tab_items(items)
}

func (tab *Tab) SelectTab(index int) {
	tab.tab_container.select_tab(index, false)
}

func (tab *Tab) SelectedTab() (int, TabItem) {
	if tab.tab_container.tab_items.Len() == 0 {
		return 0, TabItem{}
	}
	return tab.tab_container.selected_item_index, tab.tab_container.tab_items.At(tab.tab_container.selected_item_index).tab_item
}

func (tab *Tab) ItemCount() int {
	return tab.tab_container.item_count()
}

func (tab *Tab) OnSelect(fn func(from TabItemContainer, to TabItemContainer, by_user bool)) {
	tab.tab_container.listeners.on_select = fn
}

func (tab *Tab) OnSwap(fn func(from TabItemContainer, to TabItemContainer)) {
	tab.tab_container.listeners.on_swap = fn
}

func (tab *Tab) SetClosable(closable bool) {
	tab.tab_container.closable = closable
}

func (tab *Tab) OnClose(fn func(closed TabItemContainer)) {
	tab.tab_container.listeners.on_close = fn
}

func (tab *Tab) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab_container)
	tab.panel.SetStyle(widget.PanelStyleSide)
	tab.panel.SetAutoBorder(true)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&tab.panel)
	return nil
}

func (tab *Tab) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &tab.panel,
				Size:   gui.FixedSize(tab.tab_container.Measure(ctx, gui.Constraints{}).Y),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *Tab) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx)*4 + 4
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = tab.tab_container.Measure(ctx, constraints).Y + 4
	}

	return point
}

func (tab *Tab) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	r := basic.BorderRadius(ctx)
	border_type := basicwidgetdraw.RoundedRectBorderTypeRegular
	
	background_color := basicwidgetdraw.BackgroundSecondaryColor(cm)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 2, border_type)
}
