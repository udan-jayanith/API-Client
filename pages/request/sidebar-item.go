package request_page

import (
	"Zbolt/basic"
	"Zbolt/icons"
	"fmt"
	"image"
	"image/color"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type sidebar_item_widget[T any] struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text
	val         T

	contextmenu_open        bool
	contextmenu_area_cached *widget.ContextMenuArea[struct{}]

	rename_popup_open   bool
	rename_popup_cached *popup_skeleton_area
}

func (sd *sidebar_item_widget[T]) SetSideBarItem(sidebar_item SidebarItem[T]) {
	if sidebar_item.IconName == "" {
		log.Fatalln("Sidebar item dosen't have a icon")
	} else {
		sd.icon_widget.SetIcon(sidebar_item.IconName)
	}

	sd.text_widget.SetEllipsisString("...")
	sd.text_widget.SetValue(sidebar_item.Text)

	sd.val = sidebar_item.Value
}

func (sd *sidebar_item_widget[T]) contextmenu_area(ctx *gui.Context) *widget.ContextMenuArea[struct{}] {
	if sd.contextmenu_area_cached != nil {
		return sd.contextmenu_area_cached
	}
	val, ok := ctx.Env(sd, sidebar_item_contextmenu_env)
	if !ok {
		panic("Sitebar item context menu not found")
	}
	contextmenu_area, ok := val.(*widget.ContextMenuArea[struct{}])
	if !ok {
		panic("Expected *widget.ContextMenuArea[string]\nBut got somthing else")
	}
	sd.contextmenu_area_cached = contextmenu_area
	return contextmenu_area
}

func (sd *sidebar_item_widget[T]) build_context_menu(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !sd.contextmenu_open {
		return nil
	}

	contextmenu_area := sd.contextmenu_area(ctx)
	contextmenu_area.PopupMenu().SetItemsByStrings([]string{"Rename", "Delete"})
	contextmenu_area.PopupMenu().OnClose(func(context *gui.Context, reason widget.PopupCloseReason) {
		sd.contextmenu_open = false
	})
	contextmenu_area.PopupMenu().OnItemSelected(func(ctx *gui.Context, index int) {
		switch index {
		case 0:
			sd.rename_popup_open = true
			popup := sd.rename_popup(ctx)
			popup.SetOpen(true)
			popup.OnClose(func(context *gui.Context, reason widget.PopupCloseReason) {
				sd.rename_popup_open = false
			})
			popup.FocusInput(ctx)
			gui.RequestRedraw(sd)
		case 1:
			on_delete, ok := basic.GetEnvMust[OnItemDeleteFunc[T]](ctx, sd, on_item_delete_env)
			if ok && on_delete != nil {
				on_delete(ctx, sd.path(ctx), sd.item())
			}
		}
	})
	adder.AddWidget(contextmenu_area)
	return nil
}

func (sd *sidebar_item_widget[T]) rename_popup(ctx *gui.Context) *popup_skeleton_area {
	if sd.rename_popup_cached != nil {
		return sd.rename_popup_cached
	}
	val, ok := ctx.Env(sd, sidebar_item_rename_popup_env)
	if !ok {
		panic("Sitebar item rename popup not found")
	}
	popup, ok := val.(*popup_skeleton_area)
	if !ok {
		panic("Expected *widget.Popup\nBut got somthing else")
	}
	sd.rename_popup_cached = popup
	return popup
}

func (sd *sidebar_item_widget[T]) path(ctx *gui.Context) string {
	path, ok := basic.GetEnvMust[string](ctx, sd, path_env)
	if !ok {
		panic("path_env failed")
	}
	return path
}

func (sd *sidebar_item_widget[T]) item() SidebarItem[T] {
	return SidebarItem[T]{
		Text:     sd.text_widget.Value(),
		IconName: sd.icon_widget.IconName(),
		Value:    sd.val,
	}
}

func (sd *sidebar_item_widget[T]) build_rename_menu(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !sd.rename_popup_open {
		return nil
	}

	popup := sd.rename_popup(ctx)
	popup.SetButtonText("Rename")
	popup.SetHeading(fmt.Sprintf("Rename %s to:", sd.text_widget.Value()))
	popup.OnResult(func(new_name string) {
		on_item_rename, ok := basic.GetEnvMust[OnItemRenameFunc[T]](ctx, sd, on_item_rename_env)
		if ok && on_item_rename != nil {
			on_item_rename(ctx, sd.path(ctx), sd.item(), new_name)
		}
	})
	adder.AddWidget(popup)
	return nil
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sd.icon_widget.SetSize(widget.LineHeight(ctx))
	adder.AddWidget(&sd.icon_widget)
	adder.AddWidget(&sd.text_widget)

	if err := sd.build_rename_menu(ctx, adder); err != nil {
		return err
	}
	return sd.build_context_menu(ctx, adder)
}

func (sd *sidebar_item_widget[T]) padding(ctx *gui.Context) gui.Padding {
	u := widget.UnitSize(ctx)
	padding := basic.NewPadding(u/16, u/8)
	return padding
}

func (sd *sidebar_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	gap := u / 6
	b := widgetBounds.Bounds()
	if sd.contextmenu_open {
		layouter.LayoutWidget(sd.contextmenu_area(ctx), b)
	}
	if sd.rename_popup_open {
		layouter.LayoutWidget(sd.rename_popup(ctx), b.Add(image.Pt(widget.UnitSize(ctx), 0)))
	}

	padding := sd.padding(ctx)
	b.Min.X += padding.Start
	b.Max.X -= padding.End
	b.Min.Y += padding.Top
	b.Max.Y -= padding.Bottom

	icon_size := sd.icon_widget.Measure(ctx, gui.Constraints{})
	icon_bounds := b
	icon_bounds.Max.X = b.Min.X + icon_size.X
	layouter.LayoutWidget(&sd.icon_widget, icon_bounds)

	text_bounds := b
	text_bounds.Min.X = icon_bounds.Max.X + gap
	layouter.LayoutWidget(&sd.text_widget, text_bounds)

}

func (sd *sidebar_item_widget[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var size image.Point
	u := widget.UnitSize(ctx)
	padding := sd.padding(ctx)

	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = u*6 + padding.End + padding.Start
	}

	if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	} else {
		size.Y = widget.LineHeight(ctx) + padding.Top + padding.Bottom
	}
	return size
}

func (sd *sidebar_item_widget[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	result := sd.contextmenu_area(ctx).HandlePointingInput(ctx, widgetBounds)
	if result.IsHandled() {
		sd.contextmenu_open = true
		gui.RequestRebuild(sd)
	}
	return result
}

func (sd *sidebar_item_widget[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	if widgetBounds.IsHitAtCursor() {
		basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), color.Alpha16{2505}, basic.BorderRadius(ctx))
	}
}

// SidebarItem
// call OnRename
// call OnDelete
