package request_page

import (
	//CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type nothing_widget struct {
	gui.DefaultWidget

	icon     icons.Icon
	tooltip  basicwidget.TooltipArea
	on_click func()
}

func (nw *nothing_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	//nw.image.SetTooltip("Create a new request")
	nw.tooltip.SetText("Create a new request")
	adder.AddWidget(&nw.tooltip)

	nw.icon.SetIcon("large-icons/add-box")
	u := widget.UnitSize(ctx)
	nw.icon.SetSize(u * 7)
	if nw.on_click != nil {
		nw.icon.OnClick(nw.on_click)
	}

	adder.AddWidget(&nw.icon)
	return nil
}

func (nw *nothing_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	icon_size := nw.icon.Size()
	b := ctx.AppBounds()
	b.Min.Y = (b.Min.Y + b.Dy()/2) - (icon_size.Y / 2)
	b.Max.Y = b.Min.Y + icon_size.Y

	b.Min.X = (b.Min.X + b.Dx()/2) - (icon_size.X / 2)
	b.Max.X = b.Min.X + icon_size.X

	layouter.LayoutWidget(&nw.icon, b)
	layouter.LayoutWidget(&nw.tooltip, b)
}

func (nw *nothing_widget) OnClick(fn func()) {
	nw.on_click = fn
}
