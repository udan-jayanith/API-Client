package request_page

import (
	CommonWidgets "API-Client/common-widgets"
	"API-Client/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type nothing_widget struct {
	gui.DefaultWidget

	image    CommonWidgets.WidgetWithTooltip[*icons.Icon]
	on_click func()
}

func (nw *nothing_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	nw.image.SetTooltip("Create a new request")
	image := nw.image.Widget()
	image.SetIcon("large-icons/add-box")
	u := widget.UnitSize(ctx)
	image.SetSize(u * 8)
	if nw.on_click != nil {
		image.OnClick(nw.on_click)
	}

	adder.AddWidget(&nw.image)
	return nil
}

func (nw *nothing_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	flex1_item := gui.LinearLayoutItem{
		Size: gui.FlexibleSize(1),
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			flex1_item,
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionVertical,
					Items: []gui.LinearLayoutItem{
						flex1_item,
						{
							Widget: &nw.image,
						},
						flex1_item,
					},
				},
			},
			flex1_item,
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (nw *nothing_widget) OnClick(fn func()) {
	nw.on_click = fn
}
