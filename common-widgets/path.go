package CommonWidgets

import (
	"Zbolt/basic"
	"Zbolt/icons"
	"image"
	"image/color"
	"path/filepath"
	"strings"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	draw "github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type path_segment_widget struct {
	gui.DefaultWidget

	path_name   string
	text_widget widget.Text

	path_widget *path_widget
	index       int
}

func (psw *path_segment_widget) Build(context *gui.Context, adder *gui.ChildAdder) error {
	text := psw.path_name
	text_widget := &psw.text_widget
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	text_widget.SetValue(text)
	adder.AddWidget(text_widget)
	return nil
}

func (psw *path_segment_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(0, widget.RoundedCornerRadius(ctx)),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &psw.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (psw *path_segment_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := psw.text_widget.Measure(ctx, gui.Constraints{})
	point.X += widget.UnitSize(ctx) / 2
	return point
}

func (psw *path_segment_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	var background_color color.Color
	cm := ctx.ColorMode()

	if widgetBounds.IsHitAtCursor() {
		background_color = draw.BackgroundColor(cm)
	} else {
		background_color = draw.BackgroundSecondaryColor(cm)
	}

	draw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, widget.RoundedCornerRadius(ctx))
}

func (psw *path_segment_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !(widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)) || psw.path_widget.on_select == nil {
		return gui.HandleInputResult{}
	}

	var path string
	for i := 0; i <= psw.index; i++ {
		path = filepath.Join(path, psw.path_widget.segments[i].text_widget.Value())
	}
	psw.path_widget.on_select(ctx, path)

	return gui.HandleInputResult{}
}

type path_widget struct {
	gui.DefaultWidget

	segments  []path_segment_widget
	separator []*icons.Icon

	on_select func(ctx *gui.Context, path string)
}

func (pw *path_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	l := len(pw.segments)
	for i, _ := range pw.segments {
		segment := &pw.segments[i]
		adder.AddWidget(segment)

		if i != l-1 {
			separator := pw.separator[i]

			line_height := widget.LineHeight(ctx)
			size := line_height - line_height/3
			separator.SetSize(size)

			adder.AddWidget(separator)
		}
	}
	return nil
}

func (pw *path_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(widget.UnitSize(ctx)/3-10, 0),
		Items:     make([]gui.LinearLayoutItem, 0, len(pw.segments)),
	}

	l := len(pw.segments)
	for i, _ := range pw.segments {
		segment := &pw.segments[i]
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: segment,
		})

		if i != l-1 {
			separator := pw.separator[i]
			layout.Items = append(layout.Items, gui.LinearLayoutItem{
				Widget: separator,
			})
		}
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (pw *path_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	point.Y = widget.LineHeight(ctx) + widget.UnitSize(ctx)/3

	for i, _ := range pw.segments {
		measurements := pw.segments[i].Measure(ctx, constraints)
		point.X += measurements.X
	}

	line_height := widget.LineHeight(ctx)
	size := line_height - line_height/3
	point.X += (len(pw.segments) - 1) * size
	return point
}

type Path struct {
	gui.DefaultWidget

	panel       widget.Panel
	path_widget path_widget
}

func (path *Path) Build(context *gui.Context, adder *gui.ChildAdder) error {
	path.panel.SetContent(&path.path_widget)
	path.panel.SetAutoBorder(true)
	path.panel.SetStyle(widget.PanelStyleSide)
	path.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&path.panel)
	return nil
}

func (path *Path) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(3),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &path.panel,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (path *Path) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	measurements := path.path_widget.Measure(ctx, gui.Constraints{})
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = measurements.X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = measurements.Y
	}

	point.X += 3 * 2
	point.Y += 3 * 2
	return point
}

func (path *Path) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	r := basic.BorderRadius(ctx)
	border_type := draw.RoundedRectBorderTypeRegular

	background_color := draw.BackgroundSecondaryColor(cm)
	draw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := draw.BorderColors(cm, border_type)
	draw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 2, border_type)
}

func (path_widget *Path) SetPath(directory_path string) {
	if directory_path == path_widget.Path() {
		return
	}
	directory_path = filepath.Clean(directory_path)
	list := strings.Split(filepath.ToSlash(directory_path), string(filepath.Separator))
	l := len(list)

	path_widget.path_widget.segments = make([]path_segment_widget, 0, l)
	path_widget.path_widget.separator = make([]*icons.Icon, 0, l-1)

	for i, path_name := range list {
		if path_name == "" {
			continue
		}
		path_widget.path_widget.segments = append(path_widget.path_widget.segments, path_segment_widget{
			path_name:   path_name,
			path_widget: &path_widget.path_widget,
			index:       i,
		})

		if i != l-1 {
			path_widget.path_widget.separator = append(path_widget.path_widget.separator, icons.NewIcon("arrow_forward"))
		}
	}
	gui.RequestRebuild(&path_widget.path_widget)
}

func (path_widget *Path) OnSelect(fn func(ctx *gui.Context, path string)) {
	path_widget.path_widget.on_select = fn
}

func (path_widget *Path) Path() string {
	var path string
	for i, _ := range path_widget.path_widget.segments {
		seg := path_widget.path_widget.segments[i].text_widget.Value()
		path = filepath.Join(path, seg)
	}
	return path
}
