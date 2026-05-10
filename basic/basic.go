package basic

import (
	draw_color "API-Client/internal/draw"
	"image/color"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

const (
	HeadingSize float64 = 1.8
)

func NewPadding(padding ...int) gui.Padding {
	l := len(padding)
	if l > 4 {
		panic("Extra arguments to NewPadding function")
	} else if l == 3 || l == 0 {
		panic("Invalid number of arguments to NewPadding function")
	} else if l == 4 {
		return gui.Padding{
			Top:    padding[0],
			Bottom: padding[2],
			Start:  padding[3],
			End:    padding[1],
		}
	} else if l == 2 {
		return gui.Padding{
			Top:    padding[0],
			Bottom: padding[0],
			Start:  padding[1],
			End:    padding[1],
		}
	}
	return gui.Padding{
		Top:    padding[0],
		Bottom: padding[0],
		Start:  padding[0],
		End:    padding[0],
	}
}

func NewTextWidget(text string) *widget.Text {
	text_widget := widget.Text{}
	text_widget.SetValue(text)
	return &text_widget
}

type Alignment int8

const (
	Start Alignment = iota + 1
	Center
	End
)

func align(items []gui.LinearLayoutItem, alignment Alignment) []gui.LinearLayoutItem {
	switch alignment {
	case Start:
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	case Center:
		items = append([]gui.LinearLayoutItem{
			{
				Size: gui.FlexibleSize(1),
			},
		}, items...)
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	case End:
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	}
	return items
}

func getOppositeLayout(direction gui.LayoutDirection) gui.LayoutDirection {
	if direction == gui.LayoutDirectionHorizontal {
		return gui.LayoutDirectionVertical
	}
	return gui.LayoutDirectionHorizontal
}

func getLayoutAlignment(direction gui.LayoutDirection, horizontal, vertical Alignment) Alignment {
	if direction == gui.LayoutDirectionHorizontal {
		return horizontal
	} else {
		return vertical
	}
}

func Align(layout1 gui.LinearLayout, horizontal, vertical Alignment) gui.LinearLayout {
	layout1.Items = align(layout1.Items, getLayoutAlignment(layout1.Direction, horizontal, vertical))

	layout2 := gui.LinearLayout{
		Direction: getOppositeLayout(layout1.Direction),
		Items: []gui.LinearLayoutItem{
			{
				Layout: layout1,
			},
		},
	}
	layout2.Items = align(layout2.Items, getLayoutAlignment(layout2.Direction, horizontal, vertical))

	return layout2
}

func BorderRadius(ctx *gui.Context) int {
	return widget.RoundedCornerRadius(ctx)
}

func Gap(ctx *gui.Context) int {
	return BorderRadius(ctx)
}

func LineWidth(ctx *gui.Context) float32 {
	return 1 * float32(ctx.Scale())
}

func LineColor(ctx *gui.Context) color.Color {
	return draw_color.ScaleAlpha(draw_color.Color(ctx.ColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)
}
