package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	"image"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type SidebarItem[T any] struct {
	Text, IconName string
	Value          T
}

type sidebar_item_widget[T any] struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text
}

func (sd *sidebar_item_widget[T]) SetSideBarItem(sidebar_item SidebarItem[T]) {
	if sidebar_item.IconName == "" {
		log.Fatalln("Sidebar item dosen't have a icon")
	} else {
		sd.icon_widget.SetIcon(sidebar_item.IconName)
	}

	sd.text_widget.SetEllipsisString("...")
	sd.text_widget.SetValue(sidebar_item.Text)
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sd.icon_widget.SetSize(widget.LineHeight(ctx))
	adder.AddWidget(&sd.icon_widget)
	adder.AddWidget(&sd.text_widget)

	return nil
}

func (sd *sidebar_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	gap := u / 6
	b := widgetBounds.Bounds()

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

	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = u * 6
	}

	if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	} else {
		size.Y = widget.LineHeight(ctx)
	}
	return size
}

/*
type Sidebar[T comparable] struct {
	gui.DefaultWidget

	options struct {
		search_icon   *ebiten.Image
		search_widget CommonWidgets.WidgetWithTooltip[*CommonWidgets.TextInputWithContextMenu]
		add           struct {
			create_request_button, create_folder_button, variable_button CommonWidgets.ButtonWithTooltip
			create_request_icon, create_folder_icon, variable_icon       *ebiten.Image

			on_variable_clicked func(ctx *gui.Context)
			on_request_create   func(ctx *gui.Context)

			folder_popup          folder_create_popup
			folder_popup_position image.Point
			on_folder_create      func(ctx *gui.Context, folder_name string)
		}
	}

	list struct {
		path            CommonWidgets.WidgetWithTooltip[*CommonWidgets.Path]
		widget          widget.List[T]
		items           []widget.ListItem[T]
		on_item_clicked func(value T)

		// TODO: use context menu area
		contextmenu struct {
			menu     widget.PopupMenu[struct{}]
			position image.Point

			rename_popup_widget folder_create_popup
			right_clicked_item  *sidebar_item_widget[T]
		}
	}

	/*

} */

type sidebar_header_widget struct {
	gui.DefaultWidget

	path_widget           CommonWidgets.WidgetWithTooltip[*CommonWidgets.Path]
	on_path_value_changed func(ctx *gui.Context, path string)
	options               struct {
		create_request_button, create_folder_button, variable_panel_button CommonWidgets.ButtonWithTooltip
		request_icon, folder_icon, variable_icon                           *ebiten.Image

		on_variable_panel_click func(ctx *gui.Context)
		on_create_request_click func(ctx *gui.Context)
		on_folder_create_click  func(ctx *gui.Context, folder_name, dir string)
	}

	search_bar                  CommonWidgets.WidgetWithTooltip[*CommonWidgets.TextInputWithContextMenu]
	search_icon                 *ebiten.Image
	on_search_bar_value_changed func(ctx *gui.Context, value string)
}

func (header *sidebar_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
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

	header.options.variable_panel_button.SetTooltip("Open variables panel")
	header.options.variable_panel_button.SetIcon(header.options.variable_icon)
	adder.AddWidget(&header.options.variable_panel_button)

	header.search_bar.SetTooltip("Search bar")
	if header.search_icon == nil {
		header.options.variable_icon = icons.Store.Open("search")
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
		size.Y= h
	} else {
		// TODO:
	}
	return size
}

func (header *sidebar_header_widget) SetPath(path string) {
	header.path_widget.Widget().SetPath(path)
}

func (header *sidebar_header_widget) Path() string {
	return header.path_widget.Widget().Path()
}

type Sidebar[T any] struct {
	gui.DefaultWidget

	sidebar_items           gui.WidgetSlice[*sidebar_item_widget[T]]
	on_sidebar_item_clicked func(ctx *gui.Context, sidebar_item SidebarItem[T])
	on_sidebar_item_rename  func(ctx *gui.Context, sidebar_item SidebarItem[T], new_name string)
	on_sidebar_item_delete  func(ctx *gui.Context, sidebar_item SidebarItem[T])
}

func (sidebar *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	return nil
}

func (sidebar *Sidebar[T]) Layout(context *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
}

func (sidebar *Sidebar[T]) SetSidebarItems(items SidebarItem[T]) {
	// TODO: Finish this
}

func (sidebar *Sidebar[T]) SetPath(path string) {
	sidebar.path_widget.SetTooltip("Path")
	sidebar.path_widget.Widget().SetPath(path)
}
