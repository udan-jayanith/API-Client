// Adapted from guigui
package tooltip

import (
	"image"
	"time"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type TooltipArea struct {
	guigui.DefaultWidget

	popup   basicwidget.Popup
	content guigui.Widget

	t time.Time
}

func (t *TooltipArea) SetContent(widget guigui.Widget) {
	t.content = widget
}

func (t *TooltipArea) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	if t.popup.IsOpen() {
		t.popup.SetContent(t.content)
		t.popup.SetModal(false)

		adder.AddWidget(&t.popup)
	}
	return nil
}

func (t *TooltipArea) tooltip_bounds(ctx *guigui.Context, widgetBounds *guigui.WidgetBounds) image.Rectangle {
	content_bounds := t.content.Measure(ctx, guigui.Constraints{})
	area_bounds := widgetBounds.Bounds()
	u := basicwidget.UnitSize(ctx)
	gap := u / 8

	tooltip_bounds := image.Rectangle{
		Min: image.Pt(area_bounds.Min.X, area_bounds.Min.Y-content_bounds.Y+gap),
	}
	tooltip_bounds.Max = image.Point{
		X: tooltip_bounds.Min.X + content_bounds.X,
		Y: tooltip_bounds.Min.Y + content_bounds.Y - gap,
	}

	// Clamp to app bounds so it doesn't go off screen.
	app_bounds := ctx.AppBounds()
	if app_bounds.Min.X > tooltip_bounds.Min.X {
		// left
		tooltip_bounds.Min.X = app_bounds.Min.X
	}
	if tooltip_bounds.Max.X > app_bounds.Max.X {
		// right
		tooltip_bounds.Max.X = app_bounds.Max.X
	}

	if app_bounds.Min.Y > tooltip_bounds.Min.Y {
		// top
		tooltip_bounds.Min.Y = app_bounds.Min.Y
	}
	return tooltip_bounds
}

func (t *TooltipArea) Layout(ctx *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	if t.popup.IsOpen() {
		layouter.LayoutWidget(&t.popup, t.tooltip_bounds(ctx, widgetBounds))
	}
}

func (t *TooltipArea) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	// Returning zero keeps a TooltipArea from contributing an unexpected size when used as an item
	// in a layout such as LinearLayout, which would otherwise pick up the inherited DefaultWidget size.
	return image.Point{}
}

func (t *TooltipArea) HandlePointingInput(ctx *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	tooltip_bounds := t.tooltip_bounds(ctx, widgetBounds)
	tooltip_bounds.Max.Y += basicwidget.UnitSize(ctx) / 8 // gap
	cursor_pos := image.Pt(ebiten.CursorPosition())

	if cursor_pos.In(widgetBounds.Bounds()) || t.popup.IsOpen() && cursor_pos.In(tooltip_bounds) {
		if t.t.IsZero() {
			t.t = time.Now()
		} else if time.Since(t.t).Milliseconds() >= 500 && !t.popup.IsOpen() {
			t.popup.SetOpen(true)
		}
	} else {
		t.t = time.Time{}
		t.popup.SetOpen(false)
	}
	return guigui.HandleInputResult{}
}
