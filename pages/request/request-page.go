package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	http_widget "Zbolt/pages/request/http"
	requests_handler "Zbolt/pages/request/requests-handler"
	websocket_widget "Zbolt/pages/request/websocket"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type sidebar_item struct {
	path      string
	name      string
	is_folder bool
	item_type requests_handler.RequestType
}

// TODO: Close the HTTP_Data.
type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       Sidebar[sidebar_item]
	sidebar_items []SidebarItem[sidebar_item]

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
	item := sidebar_item{
		path:      path,
		name:      request_name,
		item_type: request_type,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		Text:     item.name,
		IconName: item.item_type.IconName(),
		Value:    item,
	})
}

func (rp *RequestPage) on_folder_create(ctx *gui.Context, path string, folder_name string) {
	item := sidebar_item{
		path:      path,
		name:      folder_name,
		is_folder: true,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		Text:     item.name,
		IconName: "folder",
		Value:    item,
	})
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&rp.background)

	rp.sidebar.OnRequestItemCreate(rp.on_item_create)
	rp.sidebar.OnFolderCreate(rp.on_folder_create)
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

	w := ctx.AppBounds().Dx()
	if w > 1200 {
		layout.Items[0].Size = gui.FixedSize(widget.UnitSize(ctx) * 10)
	} else if w > 1400 {
		layout.Items[0].Size = gui.FixedSize(widget.UnitSize(ctx) * 12)
	} else if w > 1500 {
		layout.Items[0].Size = gui.FixedSize(widget.UnitSize(ctx) * 14)
	}

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
