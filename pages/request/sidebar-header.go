package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	requests_handler "Zbolt/pages/request/requests-handler"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type sidebar_header_widget struct {
	gui.DefaultWidget

	path_widget           CommonWidgets.WidgetWithTooltip[*CommonWidgets.Path]
	on_path_value_changed func(ctx *gui.Context, path string)
	options               struct {
		create_request_button, create_folder_button, variable_panel_button CommonWidgets.ButtonWithTooltip
		request_icon, folder_icon, variable_icon                           *ebiten.Image

		folder_create_popup popup_skeleton_area
		on_folder_create    OnFolderCreateFunc

		on_request_item_create OnRequestItemCreateFunc
		request_create_panel   request_create_panel

		variable_panel variable_panel_widget
	}

	search_bar                  CommonWidgets.TextInputWithContextMenu
	search_icon                 *ebiten.Image
	on_search_bar_value_changed func(ctx *gui.Context, value string)

	popup         widget.Popup
	popup_content gui.Widget
}

func (header *sidebar_header_widget) on_folder_create_btn_clicked(ctx *gui.Context) {
	header.options.folder_create_popup.SetOpen(true)
	header.options.folder_create_popup.OnResult(func(folder_name string) {
		if header.options.on_folder_create != nil {
			header.options.on_folder_create(ctx, header.Path(), folder_name)
		}
		header.options.folder_create_popup.SetOpen(false)
	})
	header.options.folder_create_popup.SetButtonText("Create")
	header.options.folder_create_popup.SetHeading("Enter the folder name:")
}

func (header *sidebar_header_widget) setup_popup() {
	header.popup.SetBackgroundSemanticColor(basicwidgetdraw.SemanticColorBase)
	header.popup.SetModal(true)
	header.popup.SetCloseByClickingOutside(true)
	header.popup.SetBackgroundBlurred(true)
	header.popup.SetBackgroundDark(true)
	header.popup.SetAnimated(false)
}

func (header *sidebar_header_widget) set_popup_content(content gui.Widget) {
	header.popup.SetContent(content)
	header.popup_content = content
}

func (header *sidebar_header_widget) on_request_panel_btn_clicked(ctx *gui.Context) {
	header.setup_popup()
	header.popup.OnClose(func(context *gui.Context, reason widget.PopupCloseReason) {
		header.options.request_create_panel.Clear()
	})
	header.options.request_create_panel.OnCreate(func(ctx *gui.Context, path string, request_name string, request_type requests_handler.RequestType) {
		println(request_name, path)
		if header.options.on_request_item_create != nil {
			header.options.on_request_item_create(ctx, path, request_name, request_type)
		}
		header.popup.SetOpen(false)
	})
	header.set_popup_content(&header.options.request_create_panel)
	header.popup.SetOpen(true)
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
	header.options.create_request_button.OnDown(header.on_request_panel_btn_clicked)
	adder.AddWidget(&header.options.create_request_button)

	header.options.create_folder_button.SetTooltip("Create new folder")
	header.options.create_folder_button.SetIcon(header.options.folder_icon)
	header.options.create_folder_button.OnDown(header.on_folder_create_btn_clicked)
	if header.options.folder_create_popup.IsOpen() {
		adder.AddWidget(&header.options.folder_create_popup)
	}
	adder.AddWidget(&header.options.create_folder_button)

	header.options.variable_panel_button.SetTooltip("Open variable panel")
	header.options.variable_panel_button.SetIcon(header.options.variable_icon)
	header.options.variable_panel_button.OnDown(func(context *gui.Context) {
		header.setup_popup()
		header.set_popup_content(&header.options.variable_panel)
		header.popup.SetOpen(true)
	})
	adder.AddWidget(&header.options.variable_panel_button)

	if header.search_icon == nil {
		header.search_icon = icons.Store.Open("search")
	}
	header.search_bar.SetIcon(header.search_icon)
	header.search_bar.SetPlaceholder("Search")
	adder.AddWidget(&header.search_bar)

	if header.popup.IsOpen() {
		adder.AddWidget(&header.popup)
	}
	return nil
}

func (header *sidebar_header_widget) layout_popup(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !header.popup.IsOpen() {
		return
	}

	content_size := header.popup_content.Measure(ctx, gui.Constraints{})
	b := ctx.AppBounds().Bounds()
	b.Min.X += b.Dx()/2 - content_size.X/2
	b.Max.X = b.Min.X + content_size.X

	b.Min.Y += b.Dy()/2 - content_size.Y/2
	b.Max.Y = b.Min.Y + content_size.Y
	layouter.LayoutWidget(&header.popup, b)
}

func (header *sidebar_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	header.layout_popup(ctx, widgetBounds, layouter)

	buttons_layout := gui.LinearLayout{
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
	}

	layout := gui.LinearLayout{
		Gap:       basic.Gap(ctx),
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &header.path_widget,
			},
			{
				Layout: buttons_layout,
			},
			{
				Widget: &header.search_bar,
			},
		},
	}

	b := widgetBounds.Bounds()
	layout.LayoutWidgets(ctx, b, layouter)

	if header.options.folder_create_popup.IsOpen() {
		b = layout.ItemBoundsAt(1, ctx, b)
		b = buttons_layout.ItemBoundsAt(1, ctx, b)
		layouter.LayoutWidget(&header.options.folder_create_popup, b)
	}
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
	header.search_bar.OnValueChanged(fn)
}

func (header *sidebar_header_widget) SetSearchBarValue(value string) {
	header.search_bar.SetValue(value)
}

func (header *sidebar_header_widget) SearchBarValue() string {
	return header.search_bar.Value()
}

func (header *sidebar_header_widget) OnRequestItemCreate(fn OnRequestItemCreateFunc) {
	header.options.on_request_item_create = fn
}

func (header *sidebar_header_widget) OnFolderCreate(fn OnFolderCreateFunc) {
	header.options.on_folder_create = fn
}
