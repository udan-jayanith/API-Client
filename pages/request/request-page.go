package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	http_widget "Zbolt/pages/request/http"
	requests_handler "Zbolt/pages/request/requests-handler"
	websocket_widget "Zbolt/pages/request/websocket"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

//TODO: Close the response body when closing.

type sidebar_item struct {
	IsFolder bool
	Data     requests_handler.FolderOrFile
}

func (si *sidebar_item) GetPath() string {
	return si.Data.Path()
}

type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar gui.WidgetWithPadding[*Sidebar[sidebar_item]]

	// made sidebar_items map[string][]SidebarItem[sidebar_item]
	current_directory string
	sidebar_items     []SidebarItem[sidebar_item]

	tabs_handler   TabsHandler
	nothing_widget blank_widget

	request_widget   CommonWidgets.WidgetWithPadding[requests_handler.RequestWidget]
	http_widget      http_widget.HTTP_Widget
	websocket_widget websocket_widget.WebsocketWidget

	request_create_widget request_create_panel
	variable_panel_widget variable_panel_widget

	popup_content gui.Widget
	popup_size    image.Point
	popup_widget  widget.Popup
}

func (rp *RequestPage) create_sidebar_item(request *requests_handler.Request) {
	request_container := sidebar_item{
		Data: request,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		IconName: request.Type.IconName(),
		Text:     request.Name(),
		Value:    request_container,
	})

	//TODO: Open a the sidebar item if there were no items before on creation.
}

func (rp *RequestPage) create_folder(path string, name string) {
	folder := requests_handler.NewFolder(path, name)
	request_container := sidebar_item{
		IsFolder: true,
		Data:     &folder,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		IconName: "folder",
		Text:     folder.Name(),
		Value:    request_container,
	})

	//TODO: Open the the sidebar item if there were no items before on creation.
}

func (rp *RequestPage) open_popup(widget gui.Widget, ctx *gui.Context) {
	rp.popup_content = widget
	rp.popup_size = widget.Measure(ctx, gui.Constraints{})
	rp.popup_widget.SetContent(rp.popup_content)
	rp.popup_widget.SetOpen(true)
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)

	adder.AddWidget(&rp.background)
	padding := basic.NewPadding(basic.Gap(ctx), 0)

	//sidebar.SetItems(rp.sidebar_items)
	rp.sidebar.SetPadding(padding)

	/*
		sidebar.OnFolderCreate(func(ctx *gui.Context, folder_name string) {
			rp.create_folder(rp.current_directory, folder_name)
		})

		sidebar.OnItemClicked(func(item sidebar_item) {
			if item.IsFolder {
				return
			}

			request := item.Data.(*requests_handler.Request)
			rp.tabs_handler.Open(request, ctx)
		})
	*/
	rp.sidebar.Widget().OnItemDelete(func(ctx *gui.Context, path string, item SidebarItem[sidebar_item]) {
		println(path, item.Text)
	})
	adder.AddWidget(&rp.sidebar)

	if rp.tabs_handler.IsEmpty() {
		adder.AddWidget(&rp.nothing_widget)
	} else {
		rp.tabs_handler.Add(adder)

		data := rp.tabs_handler.GetData(rp.tabs_handler.SelectedTab())
		if data == nil {
			panic("Invalid")
		}
		switch data.Type {
		case requests_handler.HTTP:
			rp.http_widget.SetReq(data)
			rp.request_widget.SetWidget(&rp.http_widget)
		case requests_handler.Websocket:
			rp.websocket_widget.SetReq(data)
			rp.request_widget.SetWidget(&rp.websocket_widget)
		default:
			panic("request type not handled")
		}

		rp.tabs_handler.OnSelect(func(from, to CommonWidgets.TabItemContainer, by_user bool) {
			rp.request_widget.Widget().SyncData()
		})

		rp.tabs_handler.OnClose(func(closed CommonWidgets.TabItemContainer) {
			rp.request_widget.Widget().SyncData()
		})

		rp.request_widget.SetPadding(padding)
		adder.AddWidget(&rp.request_widget)
	}

	if rp.popup_widget.IsOpen() {
		rp.popup_widget.SetBackgroundDark(true)
		rp.popup_widget.SetCloseByClickingOutside(true)
		rp.popup_widget.SetBackgroundBlurred(true)
		adder.AddWidget(&rp.popup_widget)
	}

	/*
		sidebar.OnRequestCreate(func(ctx *gui.Context) {
			//rp.request_create_widget.Clear()
			rp.open_popup(&rp.request_create_widget, ctx)
		})
	*/
	/*
			rp.request_create_widget.OnCreateButtonClicked(func(request *requests_handler.Request) {
				rp.create_sidebar_item(request)
				rp.popup_widget.SetOpen(false)
			})


		rp.sidebar.Widget().OnVariableClicked(func(ctx *gui.Context) {
			rp.open_popup(&rp.variable_panel_widget, ctx)
		})
	*/
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())

	b := widgetBounds.Bounds()
	if rp.popup_widget.IsOpen() {
		popup_content_bounds := rp.popup_size

		popup_size := image.Rectangle{
			Min: image.Pt(b.Min.X+b.Max.X/2-popup_content_bounds.X/2, b.Min.Y+b.Max.Y/2-popup_content_bounds.Y/2),
		}
		popup_size.Max = image.Pt(popup_size.Min.X+popup_content_bounds.X, popup_size.Min.Y+popup_content_bounds.Y)

		layouter.LayoutWidget(&rp.popup_widget, popup_size)
	}

	tab_container_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
	}

	if !rp.tabs_handler.IsEmpty() {
		tab_container_layout.Items = []gui.LinearLayoutItem{
			{
				Size: gui.FixedSize(basic.Gap(ctx)),
			},
			{
				Widget: rp.tabs_handler.Widget(),
			},
			{
				Widget: &rp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
		}
	} else {
		tab_container_layout.Items = append(tab_container_layout.Items, gui.LinearLayoutItem{
			Size:   gui.FlexibleSize(1),
			Widget: &rp.nothing_widget,
		})
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items: []gui.LinearLayoutItem{
			{},
			{
				Widget: &rp.sidebar,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: tab_container_layout,
				Size:   gui.FlexibleSize(4),
			},
			{},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
