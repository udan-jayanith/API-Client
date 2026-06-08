package CommonWidgets

import (
	"Zbolt/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
)

type TabContainerItem struct {
	TabItem TabItem
	Widget  gui.Widget
}

type TabContainer struct {
	gui.DefaultWidget

	tab Tab

	selected_widget gui.Widget
	tab_items       []TabItem
	on_select       func(item TabItem, index int)
	widgets         []gui.Widget
}

func (widget *TabContainer) Count() int {
	return len(widget.tab_items)
}

func (widget *TabContainer) SelectedTabContainer() (TabContainerItem, int) {
	container := TabContainerItem{}

	index, item := widget.tab.SelectedTab()
	container.TabItem = item
	container.Widget = widget.widgets[index]
	return container, index
}

func (widget *TabContainer) SelectTab(index int) {
	widget.tab.SelectTab(index)
}

func (widget *TabContainer) OnSelect(fn func(item TabItem, index int)) {
	widget.on_select = fn
}

func (widget *TabContainer) OnClose(fn func(closed TabItemContainer)) {
	widget.tab.OnClose(fn)
}

func (widget *TabContainer) SetClosable(closable bool) {
	widget.tab.SetClosable(closable)
}

func (widget *TabContainer) OnSwap(fn func(from TabItemContainer, to TabItemContainer)) {
	widget.tab.OnSwap(fn)
}

func (widget *TabContainer) SetItems(items []TabContainerItem) {
	widget.tab_items = make([]TabItem, 0, len(items))
	widget.widgets = make([]gui.Widget, 0, len(items))

	for _, item := range items {
		widget.tab_items = append(widget.tab_items, item.TabItem)
		widget.widgets = append(widget.widgets, item.Widget)
	}

	widget.tab.OnSelect(func(from, to TabItemContainer, by_user bool) {
		widget.selected_widget = widget.widgets[to.Index]
		if widget.on_select != nil {
			widget.on_select(to.Item, to.Index)
		}
	})
	widget.tab.SetTabItems(widget.tab_items)
}

func (widget *TabContainer) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&widget.tab)

	if widget.selected_widget != nil {
		adder.AddWidget(widget.selected_widget)
	}
	return nil
}

func (widget *TabContainer) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &widget.tab,
			},
			{
				Widget: widget.selected_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (widget *TabContainer) get_selected_widget() gui.Widget {
	if widget.selected_widget != nil {
		return widget.selected_widget
	} else if len(widget.tab_items) == 0 {
		panic("Tab items must be set")
	}

	return widget.widgets[0]
}

func (widget *TabContainer) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var size image.Point

	tab_h := widget.tab.Measure(ctx, gui.Constraints{}).Y

	selected_widget := widget.get_selected_widget()
	var selected_widget_size image.Point

	if selected_widget != nil {
		if h, ok := constraints.FixedHeight(); ok {
			selected_widget_size = selected_widget.Measure(ctx, gui.FixedHeightConstraints(h-tab_h))
		} else {
			selected_widget_size = selected_widget.Measure(ctx, constraints)
		}
	}

	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = selected_widget_size.X
	}

	if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	} else {
		size.Y = tab_h + selected_widget_size.Y
	}
	return size
}
