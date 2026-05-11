package websocket_widget

import (
	"Zbolt/basic"
	"Zbolt/pages/request/requests-handler"

	gui "github.com/guigui-gui/guigui"
)

type WebsocketWidget struct {
	gui.DefaultWidget

	request_widget  request_widget
	response_widget response_widget
}

func (ww *WebsocketWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&ww.request_widget)
	adder.AddWidget(&ww.response_widget)
	return nil
}

func (ww *WebsocketWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &ww.request_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &ww.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (ww *WebsocketWidget) RequestType() requests_handler.RequestType {
	return requests_handler.Websocket
}

func (ww *WebsocketWidget) SetReq(req *requests_handler.Request) {}

func (ww *WebsocketWidget) SyncData() {}
