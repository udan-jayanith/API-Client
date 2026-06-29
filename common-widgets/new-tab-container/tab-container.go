package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
)

type TabItem[T any] struct {
	Text     string
	IconName string
	Closable bool
	Movable  bool
	Value    T
}

type TabContainerItem[T any] struct {
	TabItem    TabItem[T]
	Widget     gui.Widget
	WidgetFunc func(ctx *gui.Context, index int, tab_item TabItem[T]) gui.Widget
}

var (
	// func(ctx *gui.Context, index int, tab_item TabItem[T])
	on_tab_item_select gui.EventKey = gui.GenerateEventKey()
	// func(ctx *gui.Context, index int, tab_item TabItem[T])
	on_tab_item_close = gui.GenerateEnvKey()
	// func(ctx *gui.Context, from, to int, from_item, to_item TabItem[T])
	on_tab_item_move = gui.GenerateEnvKey()
)

type TabContainer struct {
	tab_bar tab_bar_widget
}
