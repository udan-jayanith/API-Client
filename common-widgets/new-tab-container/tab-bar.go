package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
)

type tab_bar_content_widget struct {
	gui.DefaultWidget

	tab_widget_items gui.WidgetSlice[*tab_item_widget]
}

type tab_bar_widget struct {
	gui.DefaultWidget
}
