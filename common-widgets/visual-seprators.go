package CommonWidgets

import (
	"API-Client/basic"
	"image"

	draw_color "API-Client/internal/draw"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func LinePadding(ctx *gui.Context) {
	basic.NewPadding(widget.UnitSize(ctx) / 2)
}

type HorizontalLine struct {
	gui.DefaultWidget
	padding gui.Padding
}

//func (line *HorizontalLine) SetPadding(padding gui.Padding) {
//line.padding = padding
//}

func (line *HorizontalLine) width(ctx *gui.Context) float32 {
	return 1 * float32(ctx.Scale())
}

func (line *HorizontalLine) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	b.Min.Add(image.Pt(line.padding.Start, 0))
	b.Max.Sub(image.Pt(line.padding.End, 0))

	line_color := draw_color.ScaleAlpha(draw_color.Color(	ctx.ColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)
	width := line.width(ctx)

	vector.StrokeLine(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Max.X), float32(b.Min.Y), width, line_color, false)
}

func (line *HorizontalLine) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = int(line.width(ctx))
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx) * 20
	}

	return point
}

type VerticalLine struct {
	HorizontalLine
}

func (line *VerticalLine) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = widget.UnitSize(ctx) * 20
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = int(line.width(ctx))
	}

	return point
}

func (line *VerticalLine) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()

	line_color := draw_color.ScaleAlpha(draw_color.Color(	ctx.ColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)
	width := line.width(ctx)

	vector.StrokeLine(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Min.X), float32(b.Max.Y), width, line_color, false)
}
