package CommonWidgets

import (
	"image"
	"image/color"

	"Zbolt/basic"
	"Zbolt/icons"

	draw_color "Zbolt/internal/draw"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TabItem struct {
	Text  string
	Value string
	Icon  *icons.Icon
}

type tab_item struct {
	gui.DefaultWidget

	index          int
	tab_item       TabItem
	tabs_container *tabs_container // tabs_container is the tab container

	text_widget widget.Text
	close_icon  icons.Icon

	relative_cursor_axis int
}

func (item *tab_item) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if item.tab_item.Icon != nil {
		adder.AddWidget(item.tab_item.Icon)
	}

	text_widget := &item.text_widget
	text_widget.SetValue(item.tab_item.Text)
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)

	adder.AddWidget(&item.text_widget)

	if item.tabs_container.closable {
		if item.close_icon.IconName() == "" {
			line_height := widget.LineHeight(ctx)
			size := line_height - line_height/4
			item.close_icon.SetSize(size)
		}

		if item.tabs_container.selected_item_index == item.index {
			item.close_icon.SetIcon("close")
		} else {
			item.close_icon.SetIcon("close-grey")
		}

		adder.AddWidget(&item.close_icon)
	}

	return nil
}

func (item *tab_item) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(0, widget.LineHeight(ctx)/2)
}

func (item *tab_item) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (item *tab_item) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Padding:   item.padding(ctx),
		Gap:       item.gap(ctx),
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, 3),
	}

	if item.tab_item.Icon != nil {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: item.tab_item.Icon,
		})
	}

	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Widget: &item.text_widget,
		Size:   gui.FlexibleSize(1),
	})

	if item.tabs_container.closable {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &item.close_icon,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (item *tab_item) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := item.text_widget.Measure(ctx, constraints)
	padding := item.padding(ctx)

	gap := item.gap(ctx)
	if item.tab_item.Icon != nil {
		icon_measurement := item.tab_item.Icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	if item.tabs_container.closable {
		icon_measurement := item.close_icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	point.X += padding.End + padding.Start
	point.Y = widget.UnitSize(ctx)
	return point
}

func (item *tab_item) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	var background_color color.Color
	var border_type basicwidgetdraw.RoundedRectBorderType

	cm := ctx.ColorMode()
	if item.tabs_container.selected_item_index == item.index {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.2, 0.2)
		border_type = basicwidgetdraw.RoundedRectBorderTypeInset
	} else if inpututil.MouseButtonPressDuration(ebiten.MouseButton0) >= 10 && widgetBounds.IsHitAtCursor() && item.tabs_container.holding.is_holding {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.4, 0.4)
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	} else if widgetBounds.IsHitAtCursor() {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.2, 0.2)
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	} else {
		background_color = basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	}

	r := basic.BorderRadius(ctx)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 1, border_type)
}

func (item *tab_item) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	b := widgetBounds.Bounds()
	is_hovering := widgetBounds.IsHitAtCursor()

	if item.tabs_container.closable && is_hovering {
		padding := item.padding(ctx)
		tab_item_size := widgetBounds.Bounds()
		tab_item_size.Max.X -= padding.End

		close_icon_size := item.close_icon.Measure(ctx, gui.Constraints{})
		cursor_x, _ := ebiten.CursorPosition()

		if cursor_x <= tab_item_size.Max.X && cursor_x >= tab_item_size.Max.X-close_icon_size.X && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			item.tabs_container.on_close(item.index, item.tab_item)
			return gui.HandleInputResult{}
		}
	}

	if is_hovering && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		item.tabs_container.on_select(item.index, item.tab_item, true)
	} else if is_hovering && inpututil.MouseButtonPressDuration(ebiten.MouseButton0) >= 10 && item.tabs_container.selected_item_index == item.index {
		cursor_axis, _ := ebiten.CursorPosition()
		item.tabs_container.on_holding(item.index, cursor_axis-b.Min.X)
		gui.RequestRebuild(item)
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && item.tabs_container.selected_item_index == item.index {
		item.tabs_container.on_mouse_up(item.index)
		gui.RequestRebuild(item)
	}

	return gui.HandleInputResult{}
}
