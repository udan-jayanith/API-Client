package CommonWidgets

import (
	"Zbolt/basic"
	draw_color "Zbolt/internal/draw"
	"image"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type InfiniteLoadingBar struct {
	gui.DefaultWidget
	p1, p2 int
	t      time.Time
}

func (wi *InfiniteLoadingBar) Tick(ctx *gui.Context, widgetBounds *gui.WidgetBounds) error {
	if time.Since(wi.t).Milliseconds() < 12 {
		return nil
	}

	if wi.p1 < 100 {
		wi.p1++
	}
	if wi.p1 > 20 {
		wi.p2++
	}

	if wi.p2 == 100 {
		wi.p1 = 0
		wi.p2 = 0
	}

	wi.t = time.Now()
	gui.RequestRedraw(wi)
	return nil
}

func (wi *InfiniteLoadingBar) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	bar_clr := draw_color.Color(ctx.ColorMode(), draw_color.ColorTypeAccent, 0.6)
	background_clr := basicwidgetdraw.BackgroundColor(ctx.ColorMode())

	b := widgetBounds.Bounds()
	vector.FillRect(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Dx()), float32(b.Dy()), background_clr, false)

	x1 := b.Min.X + (b.Dx()*wi.p1)/100
	x2 := b.Min.X + (b.Dx()*wi.p2)/100
	if wi.p2 == 0 {
		x2 = b.Min.X
	}
	vector.FillRect(dst, float32(x2), float32(b.Min.Y), float32(x1-x2), float32(b.Dy()), bar_clr, false)
}

func (wi *InfiniteLoadingBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var size image.Point
	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = widget.UnitSize(ctx) * 4
	}

	size.Y = basic.Gap(ctx)
	return size
}
