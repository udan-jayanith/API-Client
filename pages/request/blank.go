package request_page

import (
	"Zbolt/icons"
	requests_handler "Zbolt/pages/request/requests-handler"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type blank_widget struct {
	gui.DefaultWidget

	icon    icons.Icon
	tooltip basicwidget.TooltipArea

	request_create_panel   request_create_panel
	on_request_item_create OnRequestItemCreateFunc
	popup                  widget.Popup
}

func (nw *blank_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	nw.tooltip.SetText("Create a new request")
	adder.AddWidget(&nw.tooltip)

	nw.icon.SetIcon("large-icons/add-box")
	u := widget.UnitSize(ctx)
	nw.icon.SetSize(u * 7)
	nw.icon.OnClick(func() {
		nw.popup.SetModal(true)
		nw.popup.SetCloseByClickingOutside(true)
		nw.popup.SetBackgroundBlurred(true)
		nw.popup.SetBackgroundDark(true)
		nw.popup.SetAnimated(false)
		nw.popup.SetContent(&nw.request_create_panel)
		nw.request_create_panel.OnCreate(func(ctx *gui.Context, path string, request_name string, request_type requests_handler.RequestType) {
			if nw.on_request_item_create != nil {
				nw.on_request_item_create(ctx, path, request_name, request_type)
			}
			nw.popup.SetOpen(false)
		})
		nw.popup.SetOpen(true)
	})
	adder.AddWidget(&nw.icon)

	if nw.popup.IsOpen() {
		nw.popup.OnClose(func(context *gui.Context, reason widget.PopupCloseReason) {
			nw.request_create_panel.Clear()
		})
		adder.AddWidget(&nw.popup)
	}
	return nil
}

func (nw *blank_widget) layout_popup(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !nw.popup.IsOpen() {
		return
	}

	content_size := nw.request_create_panel.Measure(ctx, gui.Constraints{})
	b := ctx.AppBounds().Bounds()
	b.Min.X += b.Dx()/2 - content_size.X/2
	b.Max.X = b.Min.X + content_size.X

	b.Min.Y += b.Dy()/2 - content_size.Y/2
	b.Max.Y = b.Min.Y + content_size.Y
	layouter.LayoutWidget(&nw.popup, b)
}

func (nw *blank_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	icon_size := nw.icon.Size()
	b := widgetBounds.Bounds()
	b.Min.Y = (b.Min.Y + b.Dy()/2) - (icon_size.Y / 2)
	b.Max.Y = b.Min.Y + icon_size.Y

	b.Min.X = (b.Min.X + b.Dx()/2) - (icon_size.X / 2)
	b.Max.X = b.Min.X + icon_size.X

	layouter.LayoutWidget(&nw.icon, b)
	layouter.LayoutWidget(&nw.tooltip, b)
	nw.layout_popup(ctx, widgetBounds, layouter)
}

func (nw *blank_widget) OnRequestItemCreate(fn OnRequestItemCreateFunc) {
	nw.on_request_item_create = fn
}
