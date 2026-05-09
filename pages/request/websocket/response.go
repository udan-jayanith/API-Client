package websocket_widget

import (
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_widget struct {
	gui.DefaultWidget

	header response_header_widget
	tab    CommonWidgets.Tab

	response_header widget.Table[struct{}]
}

func (rw *response_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.header)

	rw.tab.SetTabItems([]CommonWidgets.TabItem{
		{
			Text: "Frames",
		},
		{
			Text: "Response header",
		},
	})
	adder.AddWidget(&rw.tab)
	return nil
}

func (rw *response_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.header,
			},
			{
				Widget: &rw.tab,
			},
			{
				Size: gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
