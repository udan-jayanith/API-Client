package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	http_widget "Zbolt/pages/request/http"
	requests_handler "Zbolt/pages/request/requests-handler"
	websocket_widget "Zbolt/pages/request/websocket"
	"slices"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Close the HTTP_Data.
type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       Sidebar[requests_handler.Item]
	sidebar_items []SidebarItem[requests_handler.Item]

	tab_container       CommonWidgets.TabContainer[requests_handler.Item]
	tab_container_items []CommonWidgets.TabContainerItem[requests_handler.Item]

	blank_widget blank_widget

	tab_container_widgets struct {
		HTTP      http_widget.HTTP_Widget
		Websocket websocket_widget.WebsocketWidget
	}
}

func (rp *RequestPage) Env(ctx *gui.Context, key gui.EnvKey, source *gui.EnvSource) (any, bool) {
	if key == path_env {
		return rp.sidebar.Path(), true
	}
	return nil, false
}

func (rp *RequestPage) on_item_create(ctx *gui.Context, path string, request_name string, request_type requests_handler.RequestType) {
	item := requests_handler.NewRequest(request_type, path, request_name)
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[requests_handler.Item]{
		Text:     item.Name(),
		IconName: request_type.IconName(),
		Value:    item,
	})
}

func (rp *RequestPage) on_folder_create(ctx *gui.Context, path string, folder_name string) {
	item := requests_handler.NewFolder(path, folder_name)
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[requests_handler.Item]{
		Text:     item.Name(),
		IconName: "folder",
		Value:    item,
	})
}

func (rp *RequestPage) find_item_index(item SidebarItem[requests_handler.Item]) int {
	var index int
	for i, sidebar_item := range rp.sidebar_items {
		if sidebar_item.Value == item.Value {
			index = i
			break
		}
	}
	return index
}

func (rp *RequestPage) on_item_delete(ctx *gui.Context, path string, item SidebarItem[requests_handler.Item]) {
	i := rp.find_item_index(item)
	rp.sidebar_items = slices.Delete(rp.sidebar_items, i, i+1)
	// TODO: delete tab_item too.
}

func (rp *RequestPage) on_item_rename(ctx *gui.Context, path string, item SidebarItem[requests_handler.Item], new_name string) {
	i := rp.find_item_index(item)
	rp.sidebar_items[i].Value.Rename(new_name)
	rp.sidebar_items[i].Text = rp.sidebar_items[i].Value.Name()
	// TODO: rename tab_item too.
}

func (rp *RequestPage) find_tab_item_index_by_value(value requests_handler.Item) int {
	for i, item := range rp.tab_container_items {
		if item.TabItem.Value == value {
			return i
		}
	}
	return -1
}

func (rp *RequestPage) on_item_select(ctx *gui.Context, index int) {
	item := rp.sidebar_items[index]
	_, ok := item.Value.(*requests_handler.Folder)
	if ok {
		return
	}

	tab_item_index := rp.find_tab_item_index_by_value(item.Value)
	if tab_item_index == -1 {
		icon := icons.Icon{}
		icon.SetIcon(item.IconName)
		rp.tab_container_items = append(rp.tab_container_items, CommonWidgets.TabContainerItem[requests_handler.Item]{
			TabItem: CommonWidgets.TabItem[requests_handler.Item]{
				Text:  item.Value.Name(),
				Icon:  &icon,
				Value: item.Value,
			},
		})
	} else {
		rp.tab_container.SelectTab(tab_item_index)
	}
}

func (rp *RequestPage) on_tab_item_close(closed CommonWidgets.TabItemContainer[requests_handler.Item]) {
	rp.tab_container_items = slices.Delete(rp.tab_container_items, closed.Index, closed.Index+1)
}

func (rp *RequestPage) on_tab_item_swap(from CommonWidgets.TabItemContainer[requests_handler.Item], to CommonWidgets.TabItemContainer[requests_handler.Item]) {
	temp := rp.tab_container_items[from.Index]
	rp.tab_container_items[from.Index] = rp.tab_container_items[to.Index]
	rp.tab_container_items[to.Index] = temp
}

func (rp *RequestPage) on_tab_item_select(item CommonWidgets.TabItem[requests_handler.Item], index int) {
	selected_item := rp.sidebar.SelectedItemIndex()
	if selected_item >= 0 && rp.sidebar_items[selected_item].Value != item.Value {
		rp.sidebar.SelectItemByIndex(-1)
	}
	// TODO: select sidebar item.
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&rp.background)

	rp.sidebar.OnRequestItemCreate(rp.on_item_create)
	rp.sidebar.OnFolderCreate(rp.on_folder_create)
	rp.sidebar.OnItemDelete(rp.on_item_delete)
	rp.sidebar.OnItemRename(rp.on_item_rename)
	rp.sidebar.OnItemSelect(rp.on_item_select)

	rp.sidebar.SetSidebarItems(rp.sidebar_items)
	adder.AddWidget(&rp.sidebar)

	switch len(rp.tab_container_items) {
	case 0:
		rp.blank_widget.OnRequestItemCreate(rp.on_item_create)
		adder.AddWidget(&rp.blank_widget)
	default:
		rp.tab_container.SetItems(rp.tab_container_items)
		rp.tab_container.SetClosable(true)
		rp.tab_container.OnClose(rp.on_tab_item_close)
		rp.tab_container.OnSwap(rp.on_tab_item_swap)
		rp.tab_container.OnSelect(rp.on_tab_item_select)
		adder.AddWidget(&rp.tab_container)
	}
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Padding:   basic.NewPadding(basic.BorderRadius(ctx)),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rp.sidebar,
			},
		},
	}

	w := ctx.AppBounds().Dx() / 5
	if w <= 260 {
		w = widget.UnitSize(ctx) * 8
	}
	layout.Items[0].Size = gui.FixedSize(w)

	switch len(rp.tab_container_items) {
	case 0:
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &rp.blank_widget,
			Size:   gui.FlexibleSize(1),
		})
	default:
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &rp.tab_container,
			Size:   gui.FlexibleSize(1),
		})
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
