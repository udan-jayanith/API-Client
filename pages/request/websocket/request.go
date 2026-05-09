package websocket_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_widget struct {
	gui.DefaultWidget

	url_input_bar struct {
		url_input      CommonWidgets.TextInputWithContextMenu
		connect_button widget.Button
	}
	
	url_preview CommonWidgets.URLPreview

	tab_widget CommonWidgets.Tab
	content    struct {
		params, headers CommonWidgets.AttributeTable
		body            request_body // TODO: make a editor widget and replace with it.
		selected        gui.Widget
	}
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.url_input_bar.url_input)

	rw.url_input_bar.connect_button.SetText("Connect")
	rw.url_input_bar.connect_button.SetType(widget.ButtonTypePrimary)
	adder.AddWidget(&rw.url_input_bar.connect_button)

	adder.AddWidget(&rw.url_preview)
	
	rw.tab_widget.SetTabItems([]CommonWidgets.TabItem{
		{
			Text: "Parameters",
		},
		{
			Text: "Headers",
		},
		{
			Text: "Body",
		},
	})

	if rw.content.selected == nil {
		rw.content.selected = &rw.content.params
	}



	adder.AddWidget(&rw.tab_widget)
	adder.AddWidget(rw.content.selected)
	return nil
}

func (rw *request_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := basic.Gap(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap: gap,
		Items: []gui.LinearLayoutItem{
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap: gap,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &rw.url_input_bar.url_input,
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &rw.url_input_bar.connect_button,
						},
					},
				},
			},
			{
				Widget: &rw.url_preview,
			},
			{
				Widget: &rw.tab_widget,
			},
			{
				Widget: rw.content.selected,
				Size: gui.FlexibleSize(1),
			},
		},
	}
	
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
