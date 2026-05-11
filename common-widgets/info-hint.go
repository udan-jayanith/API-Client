package CommonWidgets

import (
	"Zbolt/basic"
	"Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type InfoHint struct {
	gui.DefaultWidget
	icon    icons.Icon
	tooltip widget.TooltipArea
}

func (w *InfoHint) SetHint(text string) {
	w.tooltip.SetText(text)
}

func (w *InfoHint) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.icon.SetIcon("i")
	w.icon.SetSize(widget.LineHeight(ctx) - 8)
	adder.AddWidget(&w.icon)

	adder.AddWidget(&w.tooltip)
	return nil
}

func (w *InfoHint) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(&w.tooltip, b)

	icon_bounds := image.Rectangle{
		Min: image.Point{
			X: b.Min.X + 4,
			Y: b.Min.Y + 4,
		},
		Max: image.Point{
			X: b.Max.X - 4,
			Y: b.Max.Y - 4,
		},
	}
	layouter.LayoutWidget(&w.icon, icon_bounds)
}

func (w *InfoHint) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	radius := w.Measure(ctx, gui.Constraints{}).X

	clr1, clr2 := basicwidgetdraw.BorderAccentColors(ctx.ColorMode(), basicwidgetdraw.RoundedRectBorderTypeRegular)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, radius, 2, basicwidgetdraw.RoundedRectBorderTypeRegular)
}

func (w *InfoHint) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	size := w.icon.Measure(ctx, constraints)
	size.X += 8
	size.Y += 8
	return size
}

type TextWithInfoHint struct {
	widget.Text
	info_hint InfoHint
}

func (w *TextWithInfoHint) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	text_size := w.Text.Measure(ctx, constraints)
	info_hint_size := w.info_hint.Measure(ctx, constraints)
	if text_size.Y < info_hint_size.Y {
		text_size.Y = info_hint_size.Y
	}
	text_size.X += basic.Gap(ctx) + info_hint_size.X
	return text_size
}

func (w *TextWithInfoHint) SetHint(text string) {
	w.info_hint.SetHint(text)
}

func (w *TextWithInfoHint) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.info_hint)
	return w.Text.Build(ctx, adder)
}

func (w *TextWithInfoHint) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.Text,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionVertical,
					Items: []gui.LinearLayoutItem{
						{Size: gui.FlexibleSize(1)},
						{
							Widget: &w.info_hint,
						},
						{Size: gui.FlexibleSize(1)},
					},
				},
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
