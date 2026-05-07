package request_page

import (
	CommonWidgets "API-Client/common-widgets"
	"API-Client/icons"
	"API-Client/pages/request/requests-handler"
	"slices"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TabsHandler struct {
	tab_widget CommonWidgets.Tab
	tab_items  []CommonWidgets.TabItem
	tabs_data  []*requests_handler.Request

	on_select func(from CommonWidgets.TabItemContainer, to CommonWidgets.TabItemContainer, by_user bool)
	on_close  func(closed CommonWidgets.TabItemContainer)
}

func (tabs *TabsHandler) OnSelect(fn func(from CommonWidgets.TabItemContainer, to CommonWidgets.TabItemContainer, by_user bool)) {
	tabs.on_select = fn
}

func (tabs *TabsHandler) OnClose(fn func(closed CommonWidgets.TabItemContainer)) {
	tabs.on_close = fn
}

// Open opens
func (tabs *TabsHandler) Open(request *requests_handler.Request, ctx *gui.Context) {
	for i, _ := range tabs.tab_items {
		if tabs.tabs_data[i].Path() == request.Path() {
			tabs.tab_widget.SelectTab(i)
			return
		}
	}

	line_height := widget.LineHeight(ctx)
	size := line_height - line_height/4
	tabs.tab_items = append(tabs.tab_items, CommonWidgets.TabItem{
		Text: request.Name(),
		Icon: icons.NewIcon(request.Type.IconName(), size),
	})
	tabs.tabs_data = append(tabs.tabs_data, request)

	tabs.tab_widget.SetTabItems(tabs.tab_items)
	tabs.tab_widget.SelectTab(len(tabs.tab_items) - 1)
}

func (tabs *TabsHandler) SelectTab(index int) {
	tabs.tab_widget.SetTabItems(tabs.tab_items)
	tabs.tab_widget.SelectTab(index)
}

func (tabs *TabsHandler) SelectedTab() int {
	tabs.tab_widget.SetTabItems(tabs.tab_items)
	index, _ := tabs.tab_widget.SelectedTab()
	return index
}

func (tabs *TabsHandler) IsEmpty() bool {
	return len(tabs.tab_items) == 0
}

func (tabs *TabsHandler) GetData(index int) *requests_handler.Request {
	if index >= len(tabs.tabs_data) {
		return nil
	}
	return tabs.tabs_data[index]
}

func (tabs *TabsHandler) Add(adder *gui.ChildAdder) {
	if tabs.on_select != nil {
		tabs.tab_widget.OnSelect(tabs.on_select)
	}

	tabs.tab_widget.SetClosable(true)
	tabs.tab_widget.OnClose(func(closed CommonWidgets.TabItemContainer) {
		tabs.tab_items = slices.Delete(tabs.tab_items, closed.Index, closed.Index+1)
		tabs.tabs_data = slices.Delete(tabs.tabs_data, closed.Index, closed.Index+1)
		if tabs.on_close != nil {
			tabs.on_close(closed)
		}
	})

	tabs.tab_widget.OnSwap(func(from, to CommonWidgets.TabItemContainer) {
		from_request := tabs.tabs_data[from.Index]
		from_tab_item := tabs.tab_items[from.Index]

		tabs.tab_items[from.Index] = tabs.tab_items[to.Index]
		tabs.tabs_data[from.Index] = tabs.tabs_data[to.Index]

		tabs.tab_items[to.Index] = from_tab_item
		tabs.tabs_data[to.Index] = from_request
	})

	tabs.tab_widget.SetTabItems(tabs.tab_items)
	adder.AddWidget(&tabs.tab_widget)
}

func (tabs *TabsHandler) Widget() *CommonWidgets.Tab {
	return &tabs.tab_widget
}
