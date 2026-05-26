package CommonWidgets

import (
	"Zbolt/basic"
	draw_color "Zbolt/internal/draw"
	"image/color"

	opener "codeberg.org/udan-jayanith/Opener"
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Hyperlink struct {
	widget.Text
	clicked bool
	url     string
}

func (h *Hyperlink) color(ctx *gui.Context) color.Color {
	return draw_color.Color2(ctx.ColorMode(), draw_color.ColorTypeAccent, 0.95, 0.4)
}

func (h *Hyperlink) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	h.Text.SetColor(h.color(ctx))
	h.Text.SetValue(h.url)
	h.Text.Build(ctx, adder)
	return nil
}

func (h *Hyperlink) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	change := widget.UnitSize(ctx) / 7
	vector.StrokeLine(dst, float32(b.Min.X), float32(b.Max.Y-change), float32(b.Max.X), float32(b.Max.Y-change), basic.LineWidth(ctx), h.color(ctx), false)
	h.Text.Draw(ctx, widgetBounds, dst)
}

func (h *Hyperlink) SetRefrence(url string) {
	h.url = url
}

func (h *Hyperlink) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	handle := h.Text.HandlePointingInput(ctx, widgetBounds)
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		opener.Open(h.url)
	}
	return handle
}

func (h *Hyperlink) CursorShape(ctx *gui.Context, widgetBounds *gui.WidgetBounds) (ebiten.CursorShapeType, bool) {
	if widgetBounds.IsHitAtCursor() {
		return ebiten.CursorShapePointer, true
	}
	return 0, true
}