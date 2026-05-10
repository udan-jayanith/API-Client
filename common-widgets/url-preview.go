package CommonWidgets

import (
	"API-Client/basic"
	"API-Client/icons"
	message_model "API-Client/message-model"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type URLPreview struct {
	gui.DefaultWidget

	url_preview TextInputWithContextMenu
	copy_button widget.Button
	tooltip     basic.TooltipHelper
}

func (up *URLPreview) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	up.url_preview.SetWrapMode(widget.WrapModeAnywhere)
	up.url_preview.SetEditable(false)
	up.url_preview.SetHorizontalAlign(widget.HorizontalAlignLeft)
	up.url_preview.SetVerticalAlign(widget.VerticalAlignTop)
	adder.AddWidget(&up.url_preview)

	up.copy_button.SetIcon(icons.Store.Open("copy-all"))
	up.copy_button.OnDown(func(context *gui.Context) {
		message_model.Show("Copied", message_model.Notify, nil)
	})
	adder.AddWidget(&up.copy_button)

	if up.tooltip.IsOpen {
		adder.AddWidget(&up.tooltip.Widget)
	}
	return nil
}

func (up *URLPreview) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (up *URLPreview) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if up.tooltip.IsOpen {
		layouter.LayoutWidget(&up.tooltip.Widget, up.tooltip.Bounds)
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       up.gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &up.url_preview,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &up.copy_button,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (up *URLPreview) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !widgetBounds.IsHitAtCursor() {
		return gui.HandleInputResult{}
	}

	gap := up.gap(ctx)
	b := widgetBounds.Bounds()
	tooltip_bounds := b
	tooltip_bounds.Min.X = 0
	tooltip_bounds.Max.X = 0
	cursor_x, _ := ebiten.CursorPosition()

	w := up.copy_button.Measure(ctx, gui.Constraints{}).X
	if cursor_x <= b.Max.X && cursor_x >= b.Max.X-w {
		tooltip_bounds.Max.X = b.Max.X
		tooltip_bounds.Min.X = b.Max.X - w
		up.tooltip.Open(true, "Copy URL", tooltip_bounds)
		return gui.HandleInputResult{}
	}

	w += gap
	if cursor_x >= b.Min.X && cursor_x <= b.Max.X-w {
		tooltip_bounds.Max.X = b.Max.X - w
		tooltip_bounds.Min.X = b.Min.X
		up.tooltip.Open(true, "URL preview", tooltip_bounds)
		return gui.HandleInputResult{}
	}

	up.tooltip.Open(false, "", tooltip_bounds)
	return gui.HandleInputResult{}
}

func (up *URLPreview) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += u / 4
		point.X += up.copy_button.Measure(ctx, gui.Constraints{}).X
		point.X += up.url_preview.Measure(ctx, gui.Constraints{}).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u * 2
	}
	return point
}

func (up *URLPreview) SetURL(url string) {
	up.url_preview.SetValue(url)
}

func (up *URLPreview) URL() string {
	return up.url_preview.Value()
}
