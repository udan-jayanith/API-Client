package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	"image"
	"image/color"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SidebarItem[T any] struct {
	Text, IconName string
	Value          T
}

var (
	sidebar_item_contextmenu_env gui.EnvKey = gui.GenerateEnvKey()
)

type sidebar_item_widget[T any] struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text
	val         T

	contextmenu_open        bool
	contextmenu_area_cached *widget.ContextMenuArea[struct{}]
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
	adder.AddWidget(contextmenu_area)
	return nil
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sd.icon_widget.SetSize(widget.LineHeight(ctx))
	adder.AddWidget(&sd.icon_widget)
	adder.AddWidget(&sd.text_widget)
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
		gui.RequestRebuild(sd)
		sd.contextmenu_open = true
	}
	return result
}

func (sd *sidebar_item_widget[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	if widgetBounds.IsHitAtCursor() {
		basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), color.Alpha16{2505}, basic.BorderRadius(ctx))
	}
}

type sidebar_header_widget struct {
	gui.DefaultWidget

	path_widget           CommonWidgets.WidgetWithTooltip[*CommonWidgets.Path]
	on_path_value_changed func(ctx *gui.Context, path string)
	options               struct {
		create_request_button, create_folder_button, variable_panel_button CommonWidgets.ButtonWithTooltip
		request_icon, folder_icon, variable_icon                           *ebiten.Image

		on_variable_panel_click func(ctx *gui.Context)
		on_create_request_click func(ctx *gui.Context)
		on_folder_create        func(ctx *gui.Context, folder_name, dir string)
	}

	search_bar                  CommonWidgets.WidgetWithTooltip[*CommonWidgets.TextInputWithContextMenu]
	search_icon                 *ebiten.Image
	on_search_bar_value_changed func(ctx *gui.Context, value string)
}

func (header *sidebar_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	header.path_widget.SetTooltip("Path")
	header.path_widget.Widget().SetPath("/home/udan/Documents")
	adder.AddWidget(&header.path_widget)

	if header.options.request_icon == nil {
		header.options.request_icon = icons.Store.Open("add-box")
		header.options.folder_icon = icons.Store.Open("create-new-folder")
		header.options.variable_icon = icons.Store.Open("variable")
	}

	header.options.create_request_button.SetTooltip("Create new request")
	header.options.create_request_button.SetIcon(header.options.request_icon)
	adder.AddWidget(&header.options.create_request_button)

	header.options.create_folder_button.SetTooltip("Create new folder")
	header.options.create_folder_button.SetIcon(header.options.folder_icon)
	adder.AddWidget(&header.options.create_folder_button)

	header.options.variable_panel_button.SetTooltip("Open variable panel")
	header.options.variable_panel_button.SetIcon(header.options.variable_icon)
	adder.AddWidget(&header.options.variable_panel_button)

	header.search_bar.SetTooltip("Search bar")
	if header.search_icon == nil {
		header.search_icon = icons.Store.Open("search")
	}
	header.search_bar.Widget().SetIcon(header.search_icon)
	adder.AddWidget(&header.search_bar)
	return nil
}

func (header *sidebar_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Gap:       basic.Gap(ctx),
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &header.path_widget,
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       basic.Gap(ctx),
					Items: []gui.LinearLayoutItem{
						{
							Widget: &header.options.create_request_button,
						},
						{
							Widget: &header.options.create_folder_button,
						},
						{
							Widget: &header.options.variable_panel_button,
						},
					},
				},
			},
			{
				Widget: &header.search_bar,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (header *sidebar_header_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var size image.Point

	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = widget.UnitSize(ctx) * 6
	}

	if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	} else {
		constaints := gui.FixedWidthConstraints(size.X)
		size.Y = header.path_widget.Measure(ctx, constaints).Y
		size.Y += header.search_bar.Measure(ctx, constaints).Y
		size.Y += widget.UnitSize(ctx)
		size.Y += basic.Gap(ctx) * 2
	}
	return size
}

func (header *sidebar_header_widget) SetPath(path string) {
	header.path_widget.Widget().SetPath(path)
}

func (header *sidebar_header_widget) Path() string {
	return header.path_widget.Widget().Path()
}

func (header *sidebar_header_widget) OnPathChanged(fn func(ctx *gui.Context, path string)) {
	header.path_widget.Widget().OnSelect(fn)
}

func (header *sidebar_header_widget) OnSearchBarValueChanged(fn func(context *gui.Context, text string, committed bool)) {
	header.search_bar.Widget().OnValueChanged(fn)
}

func (header *sidebar_header_widget) SetSearchBarValue(value string) {
	header.search_bar.Widget().SetValue(value)
}

func (header *sidebar_header_widget) SearchBarValue() string {
	return header.search_bar.Widget().Value()
}

func (header *sidebar_header_widget) OnRequestButtonClicked(fn func(context *gui.Context)) {
	header.options.create_request_button.OnDown(fn)
}

func (header *sidebar_header_widget) OnFolderButtonClicked(fn func(context *gui.Context)) {
	header.options.create_folder_button.OnDown(fn)
}

func (header *sidebar_header_widget) OnVariableButtonClicked(fn func(context *gui.Context)) {
	header.options.variable_panel_button.OnDown(fn)
}

type Sidebar[T any] struct {
	gui.DefaultWidget

	sidebar_header            sidebar_header_widget
	sidebar_items             []widget.ListItem[struct{}]
	sidebar_item_context_menu widget.ContextMenuArea[struct{}]
	list_widget               widget.List[struct{}]

	on_sidebar_item_clicked func(ctx *gui.Context, sidebar_item SidebarItem[T])
	on_sidebar_item_rename  func(ctx *gui.Context, sidebar_item SidebarItem[T], new_name string)
	on_sidebar_item_delete  func(ctx *gui.Context, sidebar_item SidebarItem[T])
}

func (sidebar *Sidebar[T]) Env(ctx *gui.Context, key gui.EnvKey, source *gui.EnvSource) (any, bool) {
	if key == sidebar_item_contextmenu_env {
		return &sidebar.sidebar_item_context_menu, true
	}
	return nil, false
}

func (sidebar *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sidebar.list_widget.SetStyle(widget.ListStyleSidebar)
	if len(sidebar.sidebar_items) == 0 {
		sidebar.list_widget.SetItems([]widget.ListItem[struct{}]{
			{
				Header:       true,
				Content:      &sidebar.sidebar_header,
				Unselectable: true,
			},
			{
				Border: true,
			},
			{
				Text: "Item",
			},
		})
	} else {
		sidebar.list_widget.SetItems(sidebar.sidebar_items)
	}
	adder.AddWidget(&sidebar.list_widget)
	return nil
}

func (sidebar *Sidebar[T]) Layout(context *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sidebar.list_widget, widgetBounds.Bounds())
}

func (sidebar *Sidebar[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sidebar.list_widget.Measure(ctx, constraints)
}

func (sidebar *Sidebar[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()

	vector.StrokeLine(dst, float32(b.Max.X), float32(b.Min.Y), float32(b.Max.X), float32(b.Max.Y), basic.LineWidth(ctx), basic.LineColor(ctx), false)
}

func (sidebar *Sidebar[T]) SetSidebarItems(items []SidebarItem[T]) {
	// TODO: Finish this
}
