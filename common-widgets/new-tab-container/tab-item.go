package CommonWidgets

import (
	"image"
	"image/color"

	"Zbolt/basic"
	"Zbolt/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type tab_item_widget struct {
	gui.DefaultWidget

	index int

	icon icons.Icon
	text widget.Text

	closable   bool
	close_icon icons.Icon

	movable  bool
	selected bool

	is_layout_cached bool
	cached_layout    gui.LinearLayout
}

func (tab_item *tab_item_widget) build_layout(ctx *gui.Context) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Padding:   basic.NewPadding(basic.BorderRadius(ctx)),
		Items:     make([]gui.LinearLayoutItem, 0, 3),
	}

	if tab_item.icon.IconName() != "" {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &tab_item.icon,
		})
	}

	text_w := tab_item.text.Measure(ctx, gui.Constraints{}).X
	var w gui.Size
	max_w := widget.UnitSize(ctx) * 6
	if text_w > max_w {
		w = gui.FixedSize(max_w)
	}
	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Widget: &tab_item.text,
		Size:   w,
	})

	if tab_item.closable {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &tab_item.close_icon,
		})
	}

	tab_item.cached_layout = layout
}

func (tab_item *tab_item_widget) layout(ctx *gui.Context) gui.LinearLayout {
	if !tab_item.is_layout_cached {
		tab_item.build_layout(ctx)
		tab_item.is_layout_cached = true
	}
	return tab_item.cached_layout
}

func (tab_item *tab_item_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if tab_item.icon.IconName() != "" {
		adder.AddWidget(&tab_item.icon)
	}

	tab_item.text.SetBold(tab_item.selected)
	tab_item.text.SetEllipsisString("...")
	adder.AddWidget(&tab_item.text)

	if tab_item.closable {
		tab_item.close_icon.SetIcon("close")
		adder.AddWidget(&tab_item.close_icon)
	}
	return nil
}

func (tab_item *tab_item_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	tab_item.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab_item *tab_item_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return tab_item.layout(ctx).Measure(ctx, constraints)
}

func (tab_item *tab_item_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	r := basic.BorderRadius(ctx)
	if tab_item.selected {
		bgc := basicwidgetdraw.BackgroundColorFromSemanticColor(ctx.ColorMode(), basicwidgetdraw.SemanticColorAccent)
		basicwidgetdraw.DrawRoundedRect(ctx, dst, b, bgc, r)
	}

	if widgetBounds.IsHitAtCursor() {
		basicwidgetdraw.DrawRoundedRect(ctx, dst, b, color.Alpha16{2505}, r)
	}
}

func (tab_item *tab_item_widget) CursorShape(ctx *gui.Context, widgetBounds *gui.WidgetBounds) (ebiten.CursorShapeType, bool) {
	if widgetBounds.IsHitAtCursor() {
		return ebiten.CursorShapePointer, true
	}
	return 0, false
}

func (tab_item *tab_item_widget) SetText(text string) {
	tab_item.is_layout_cached = false
	tab_item.text.SetValue(text)
}

func (tab_item *tab_item_widget) SetIconName(icon_name string) {
	tab_item.is_layout_cached = false
	tab_item.icon.SetIcon(icon_name)
}

func (tab_item *tab_item_widget) SetClosable(closable bool) {
	tab_item.is_layout_cached = false
	tab_item.closable = closable
}

func (tab_item *tab_item_widget) SetMovable(movable bool) {
	tab_item.movable = movable
}

func (tab_item *tab_item_widget) SetSelected(selected bool) {
	tab_item.selected = selected
}
