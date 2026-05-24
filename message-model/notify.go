package message_model

import (
	"Zbolt/basic"
	"image"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type notify_widget struct {
	gui.DefaultWidget

	text_widget widget.Text
	open        bool
	t           time.Time
	on_result   result_fn_type
}

func (w *notify_widget) SetMessage(message string) {
	w.text_widget.SetValue(message)
	if !w.open {
		w.t = time.Now()
	}
	w.open = true
}

func (w *notify_widget) OnResult(fn result_fn_type) {
	w.on_result = fn
}

func (w *notify_widget) Tick(ctx *gui.Context, widgetBounds *gui.WidgetBounds) error {
	if time.Since(w.t) >= 2*time.Second {
		if w.on_result != nil {
			w.on_result(true, ctx)
		}
		w.open = false
	}
	return nil
}

func (notify_widget *notify_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	text_widget := &notify_widget.text_widget
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	text_widget.SetEllipsisString("...")
	text_widget.SetOpacity(0.9)

	adder.AddWidget(&notify_widget.text_widget)
	return nil
}

func (notify_widget *notify_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if notify_widget.on_result != nil {
			notify_widget.on_result(true, ctx)
		}
		notify_widget.open = false
	}
	return gui.HandleInputResult{}
}

func (notify_widget *notify_widget) padding(ctx *gui.Context) gui.Padding {
	line_height := widget.LineHeight(ctx)
	return basic.NewPadding(line_height/3, line_height/2+line_height/3)
}

func (notify_widget *notify_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   notify_widget.padding(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &notify_widget.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (notify_widget *notify_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	size := notify_widget.text_widget.Measure(ctx, gui.Constraints{})
	padding := notify_widget.padding(ctx)
	size.X += padding.End + padding.Start
	size.Y += padding.Top + padding.Bottom
	max_w := widget.UnitSize(ctx) * 4
	if size.X > max_w {
		size.X = max_w
	}
	return size
}

func (notify_widget *notify_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	background_color := basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, widget.LineHeight(ctx))
}

func (notify_widget *notify_widget) Bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	measurements := notify_widget.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()

	notify_widget_bounds := image.Rectangle{
		Min: image.Point{
			X: (b.Min.X + b.Dx()/2) - measurements.X/2,
			Y: b.Max.Y - (measurements.Y * 2),
		},
	}

	notify_widget_bounds.Max = image.Point{
		X: notify_widget_bounds.Min.X + measurements.X,
		Y: notify_widget_bounds.Min.Y + measurements.Y,
	}

	return notify_widget_bounds
}
