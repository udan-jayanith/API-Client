package request_page

import (
	CommonWidgets "API-Client/common-widgets"
	"API-Client/icons"
	"fmt"
	"image"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SidebarItem[T comparable] struct {
	Text, IconName string
	Value          T
}

type sidebar_item_widget[T comparable] struct {
	gui.DefaultWidget

	icon_widget icons.Icon

	text_widget    widget.Text
	sidebar_item   SidebarItem[T]
	sidebar_widget *Sidebar[T]
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	line_height := widget.LineHeight(ctx)
	if sd.sidebar_item.IconName == "" {
		log.Fatalln("Sidebar item dosen't have a icon")
	} else {
		sd.icon_widget.SetIcon(sd.sidebar_item.IconName)
		sd.icon_widget.SetSize(line_height)
	}
	adder.AddWidget(&sd.icon_widget)

	sd.text_widget.SetValue(sd.sidebar_item.Text)
	adder.AddWidget(&sd.text_widget)

	return nil
}

func (sd *sidebar_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 6,
		Items: []gui.LinearLayoutItem{
			{},
			{
				Widget: &sd.icon_widget,
			},
			{
				Widget: &sd.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sd *sidebar_item_widget[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := sd.text_widget.Measure(ctx, constraints)
	point.X += widget.UnitSize(ctx) / 4
	point.X += sd.icon_widget.Measure(ctx, constraints).X
	return point
}
func (sd *sidebar_item_widget[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		sd.sidebar_widget.list.contextmenu.right_clicked_item = sd
	}

	return gui.HandleInputResult{}
}

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
}

func (sd *Sidebar[T]) SetItems(items []SidebarItem[T]) {
	sd.list.items = make([]widget.ListItem[T], 0, len(items))
	for _, item := range items {
		content_widget := sidebar_item_widget[T]{
			sidebar_widget: sd,
			sidebar_item:   item,
		}
		sd.list.items = append(sd.list.items, widget.ListItem[T]{
			Content: &content_widget,
			Value:   item.Value,
			Movable: false,
		})
	}
}

func (sd *Sidebar[T]) OnRequestCreate(callback func(ctx *gui.Context)) {
	sd.options.add.on_request_create = callback
}

func (sd *Sidebar[T]) OnFolderCreate(fn func(ctx *gui.Context, folder_name string)) {
	sd.options.add.on_folder_create = fn
}

func (sd *Sidebar[T]) OnItemClicked(fn func(item T)) {
	sd.list.on_item_clicked = fn
}

func (sd *Sidebar[T]) OnVariableClicked(fn func(ctx *gui.Context)) {
	sd.options.add.on_variable_clicked = fn
}

func (sd *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if sd.options.add.create_request_icon == nil {
		sd.options.add.create_request_icon = icons.Store.Open("add-box")
	}
	if sd.options.add.create_folder_icon == nil {
		sd.options.add.create_folder_icon = icons.Store.Open("create-new-folder")
	}
	if sd.options.add.variable_icon == nil {
		sd.options.add.variable_icon = icons.Store.Open("variable")
	}

	sd.options.add.create_request_button.SetTooltip("Create a new request")
	sd.options.add.create_request_button.SetIcon(sd.options.add.create_request_icon)
	sd.options.add.create_request_button.OnDown(func(context *gui.Context) {
		sd.options.add.on_request_create(ctx)
	})
	adder.AddWidget(&sd.options.add.create_request_button)

	sd.options.add.create_folder_button.SetTooltip("Create a new folder")
	sd.options.add.create_folder_button.SetIcon(sd.options.add.create_folder_icon)
	sd.options.add.create_folder_button.OnDown(func(context *gui.Context) {
		sd.options.add.folder_popup_position = image.Pt(ebiten.CursorPosition())
		sd.options.add.folder_popup.SetOpen(true)
	})
	adder.AddWidget(&sd.options.add.create_folder_button)

	sd.options.add.variable_button.SetTooltip("Open project variables panel")
	sd.options.add.variable_button.SetIcon(sd.options.add.variable_icon)
	if sd.options.add.on_variable_clicked != nil {
		sd.options.add.variable_button.OnDown(sd.options.add.on_variable_clicked)
	}
	adder.AddWidget(&sd.options.add.variable_button)

	sd.options.add.folder_popup.SetButtonText("Create")
	sd.options.add.folder_popup.SetFieldValue("Enter folder name")
	sd.options.add.folder_popup.OnButtonClicked(func(ctx *gui.Context, value string) {
		if sd.options.add.on_request_create != nil {
			sd.options.add.on_folder_create(ctx, value)
		}
		sd.options.add.folder_popup.ClearInput()
	})

	if sd.options.search_icon == nil {
		sd.options.search_icon = icons.Store.Open("search")
	}
	sd.options.search_widget.SetTooltip("Search bar")
	search_widget := sd.options.search_widget.Widget()
	search_widget.SetIcon(sd.options.search_icon)
	search_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&sd.options.search_widget)

	sd.list.widget.SetItems(sd.list.items)
	sd.list.widget.OnItemsSelected(func(context *gui.Context, indices []int) {
		if sd.list.on_item_clicked != nil {
			sd.list.on_item_clicked(sd.list.items[indices[0]].Value)
		}
	})
	adder.AddWidget(&sd.list.widget)

	sd.list.path.SetTooltip("Path relative to project path")
	path_wdiget := sd.list.path.Widget()
	path_wdiget.SetPath(`Root/zed/extensions/work/codebook`)
	path_wdiget.OnSelect(func(ctx *gui.Context, path string) {
		println(path)
	})
	adder.AddWidget(&sd.list.path)

	sd.list.contextmenu.menu.SetItemsByStrings([]string{"Rename", "Delete"})
	sd.list.contextmenu.menu.OnItemSelected(func(context *gui.Context, index int) {
		if index == 0 {
			sd.list.contextmenu.rename_popup_widget.SetOpen(true)
		}
	})

	sd.list.contextmenu.rename_popup_widget.SetButtonText("Rename")
	sd.list.contextmenu.rename_popup_widget.SetFieldValue("Set new name")

	adder.AddWidget(&sd.list.contextmenu.rename_popup_widget)
	adder.AddWidget(&sd.list.contextmenu.menu)
	adder.AddWidget(&sd.options.add.folder_popup)
	return nil
}

func (sd *Sidebar[T]) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (sd *Sidebar[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sd.list.contextmenu.menu, image.Rectangle{
		Min: sd.list.contextmenu.position,
	})

	folder_popup_measurements := sd.options.add.folder_popup.Measure(ctx, gui.Constraints{})
	layouter.LayoutWidget(&sd.options.add.folder_popup, image.Rectangle{
		Min: sd.options.add.folder_popup_position,
		Max: sd.options.add.folder_popup_position.Add(folder_popup_measurements),
	})

	rename_measurements := sd.list.contextmenu.rename_popup_widget.Measure(ctx, gui.Constraints{})
	layouter.LayoutWidget(&sd.list.contextmenu.rename_popup_widget, image.Rectangle{
		Min: sd.list.contextmenu.position,
		Max: sd.list.contextmenu.position.Add(rename_measurements),
	})

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       sd.gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &sd.list.path,
				//Size: gui.FixedSize(widget.UnitSize(ctx)),
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       sd.gap(ctx),
					Items: []gui.LinearLayoutItem{
						{
							Widget: &sd.options.add.create_request_button,
						},
						{
							Widget: &sd.options.add.create_folder_button,
						},
						{
							Widget: &sd.options.add.variable_button,
						},
					},
				},
			},
			{
				Widget: &sd.options.search_widget,
			},
			{
				Widget: &sd.list.widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sd *Sidebar[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) && sd.list.contextmenu.right_clicked_item != nil {
		sd.list.contextmenu.position = image.Pt(ebiten.CursorPosition())
		sd.list.contextmenu.menu.SetOpen(true)

		fmt.Println("right clicked", sd.list.contextmenu.right_clicked_item.text_widget.Value())

		sd.list.contextmenu.right_clicked_item = nil
	}
	return gui.HandleInputResult{}
}

func (sd *Sidebar[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sd.list.widget.Measure(ctx, constraints)
}
