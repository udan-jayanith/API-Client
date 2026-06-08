package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
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

	tab_container       CommonWidgets.TabContainer
	tab_container_items []CommonWidgets.TabContainerItem

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
		if sidebar_item.Value.Name() == item.Value.Name() {
			index = i
			break
		}
	}
	return index
}

func (rp *RequestPage) on_item_delete(ctx *gui.Context, path string, item SidebarItem[requests_handler.Item]) {
	i := rp.find_item_index(item)
	rp.sidebar_items = slices.Delete(rp.sidebar_items, i, i+1)
}

func (rp *RequestPage) on_item_rename(ctx *gui.Context, path string, item SidebarItem[requests_handler.Item], new_name string) {
	i := rp.find_item_index(item)
	rp.sidebar_items[i].Value.Rename(new_name)
	rp.sidebar_items[i].Text = rp.sidebar_items[i].Value.Name()
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&rp.background)

	rp.sidebar.OnRequestItemCreate(rp.on_item_create)
	rp.sidebar.OnFolderCreate(rp.on_folder_create)
	rp.sidebar.OnItemDelete(rp.on_item_delete)
	rp.sidebar.OnItemRename(rp.on_item_rename)

	rp.sidebar.SetSidebarItems(rp.sidebar_items)
	adder.AddWidget(&rp.sidebar)

	switch len(rp.tab_container_items) {
	case 0:
		rp.blank_widget.OnRequestItemCreate(rp.on_item_create)
		adder.AddWidget(&rp.blank_widget)
	default:
		rp.tab_container.SetItems(rp.tab_container_items)
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
