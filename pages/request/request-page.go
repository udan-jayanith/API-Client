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

//TODO: Close the response body when closing

type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       Sidebar[requests_handler.FolderOrFile]
	sidebar_items []SidebarItem[requests_handler.FolderOrFile]

	tab_container       CommonWidgets.TabContainer
	tab_container_items []CommonWidgets.TabContainerItem

	blank_widget blank_widget

	tab_container_widgets struct {
		HTTP      http_widget.HTTP_Widget
		Websocket websocket_widget.WebsocketWidget
	}
}

// TODO: implement Env for path

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&rp.background)

	rp.sidebar.SetSidebarItems(rp.sidebar_items)
	adder.AddWidget(&rp.sidebar)

	switch len(rp.tab_container_items) {
	case 0:
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
	}else if w > 1400 {
		layout.Items[0].Size = gui.FixedSize(widget.UnitSize(ctx) * 12)
	}else if w > 1500 {
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
