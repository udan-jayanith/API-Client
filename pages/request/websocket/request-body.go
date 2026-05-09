package websocket_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_body struct {
	gui.DefaultWidget

	send_button widget.Button
	text_widget CommonWidgets.TextInputWithContextMenu
}

func (ww *request_body) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ww.send_button.SetType(widget.ButtonTypePrimary)
	ww.send_button.SetText("Send")
	adder.AddWidget(&ww.send_button)

	ww.text_widget.SetMultiline(true)
	ww.text_widget.SetAutoWrap(true)
	adder.AddWidget(&ww.text_widget)
	return nil
}

func (ww *request_body) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := basic.Gap(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       gap,
		Items: []gui.LinearLayoutItem{
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &ww.send_button,
						},
					},
				},
				Size: gui.FixedSize(ww.send_button.Measure(ctx, gui.Constraints{}).Y),
			},
			{
				Widget: &ww.text_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
