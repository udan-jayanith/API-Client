package request_page

import (
	"Zbolt/basic"
	"image"

	requests_handler "Zbolt/pages/request/requests-handler"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SidebarItem[T any] struct {
	Text, IconName string
	Value          T
}

var (
	sidebar_item_contextmenu_env  gui.EnvKey = gui.GenerateEnvKey()
	sidebar_item_rename_popup_env            = gui.GenerateEnvKey()

	on_item_rename_env = gui.GenerateEnvKey()
	on_item_delete_env = gui.GenerateEnvKey()
	path_env           = gui.GenerateEnvKey()
)

type (
	OnItemRenameFunc[T any] = func(ctx *gui.Context, path string, item SidebarItem[T], new_name string)
	OnItemDeleteFunc[T any] = func(ctx *gui.Context, path string, item SidebarItem[T])
	OnFolderCreateFunc      = func(ctx *gui.Context, path string, folder_name string)
	OnRequestItemCreateFunc = func(ctx *gui.Context, path string, request_name string, request_type requests_handler.RequestType)
)

type Sidebar[T any] struct {
	gui.DefaultWidget

	add_folder_menu           widget.PopupMenu[struct{}]
	sidebar_item_context_menu widget.ContextMenuArea[struct{}]
	rename_item_menu          popup_skeleton_area

	header        sidebar_header_widget
	sidebar_items []widget.ListItem[struct{}]
	list_widget   widget.List[struct{}]
	item          sidebar_item_widget[T]

	on_item_rename OnItemRenameFunc[T]
	on_item_delete OnItemDeleteFunc[T]
}

func (sidebar *Sidebar[T]) Env(ctx *gui.Context, key gui.EnvKey, source *gui.EnvSource) (any, bool) {
	switch key {
	case sidebar_item_contextmenu_env:
		return &sidebar.sidebar_item_context_menu, true
	case sidebar_item_rename_popup_env:
		return &sidebar.rename_item_menu, true
	case on_item_rename_env:
		return sidebar.on_item_rename, true
	case on_item_delete_env:
		return sidebar.on_item_delete, true
	case path_env:
		return sidebar.header.Path(), true
	}
	return nil, false
}

func (sidebar *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sidebar.list_widget.SetStyle(widget.ListStyleSidebar)
	if len(sidebar.sidebar_items) == 0 {
		sidebar.item.SetSideBarItem(SidebarItem[T]{
			IconName: "search",
			Text:     "Item",
		})
		sidebar.list_widget.SetItems([]widget.ListItem[struct{}]{
			{
				Header:       true,
				Content:      &sidebar.header,
				Unselectable: true,
			},
			{
				Border: true,
			},
			{
				Content: &sidebar.item,
				Padding: basic.NewPadding(0),
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

func (sidebar *Sidebar[T]) OnItemRename(fn OnItemRenameFunc[T]) {
	sidebar.on_item_rename = fn
}

func (sidebar *Sidebar[T]) OnItemDelete(fn OnItemDeleteFunc[T]) {
	sidebar.on_item_delete = fn
}

func (sidebar *Sidebar[T]) OnFolderCreate(fn OnFolderCreateFunc) {
	sidebar.header.OnFolderCreate(fn)
}

func (sidebar *Sidebar[T]) OnRequestItemCreate(fn OnRequestItemCreateFunc) {
	sidebar.header.OnRequestItemCreate(fn)
}

// TODO: add OnSearch
