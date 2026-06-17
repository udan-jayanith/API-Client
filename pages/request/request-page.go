package request_page

import (
	"Zbolt/basic"
	http_widget "Zbolt/pages/request/http"
	requests_handler "Zbolt/pages/request/requests-handler"
	websocket_widget "Zbolt/pages/request/websocket"
	"slices"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

// TODO: Close the HTTP_Data.
type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       Sidebar[requests_handler.Item]
	sidebar_items []SidebarItem[requests_handler.Item]

	blank_widget blank_widget

	tab_container_widgets struct {
		selected  requests_handler.RequestWidget
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

func (rp *RequestPage) widget_for_request_type(t requests_handler.RequestType) requests_handler.RequestWidget {
	switch t {
	case requests_handler.HTTP:
		return &rp.tab_container_widgets.HTTP
	case requests_handler.Websocket:
		return &rp.tab_container_widgets.Websocket
	default:
		panic("Not implemented a widget for RequestType")
	}
}

func (rp *RequestPage) on_item_select(ctx *gui.Context, index int) {
	item := rp.sidebar_items[index]
	req, ok := item.Value.(*requests_handler.Request)
	if !ok {
		return
	}

	if rp.tab_container_widgets.selected != nil {
		rp.tab_container_widgets.selected.SyncData()
	}
	rp.tab_container_widgets.selected = rp.widget_for_request_type(req.Type)
	rp.tab_container_widgets.selected.SetReq(req)
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rp.background)

	rp.sidebar.OnRequestItemCreate(rp.on_item_create)
	rp.sidebar.OnFolderCreate(rp.on_folder_create)
	rp.sidebar.OnItemDelete(rp.on_item_delete)
	rp.sidebar.OnItemRename(rp.on_item_rename)
	rp.sidebar.OnItemSelect(rp.on_item_select)

	rp.sidebar.SetSidebarItems(rp.sidebar_items)
	adder.AddWidget(&rp.sidebar)

	if rp.tab_container_widgets.selected == nil {
		rp.blank_widget.OnRequestItemCreate(rp.on_item_create)
		adder.AddWidget(&rp.blank_widget)
	} else {
		adder.AddWidget(rp.tab_container_widgets.selected)
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

	if rp.tab_container_widgets.selected == nil {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &rp.blank_widget,
			Size:   gui.FlexibleSize(1),
		})
	} else {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: rp.tab_container_widgets.selected,
			Size:   gui.FlexibleSize(1),
		})
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
