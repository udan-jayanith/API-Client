package sidebar

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type sidebar_header_widget struct {
	gui.DefaultWidget

	path_widget           CommonWidgets.WidgetWithTooltip[*CommonWidgets.Path]
	on_path_value_changed func(ctx *gui.Context, path string)
	options               struct {
		create_request_button, create_folder_button, variable_panel_button CommonWidgets.ButtonWithTooltip
		request_icon, folder_icon, variable_icon                           *ebiten.Image
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
