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

	header      sidebar_header_widget
	items       []widget.ListItem[struct{}]
	list_widget widget.List[struct{}]

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

func (sidebar *Sidebar[T]) list_items() []widget.ListItem[struct{}] {
	sidebar.sidebar_items()
	return sidebar.items
}

func (sidebar *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sidebar.list_widget.SetStyle(widget.ListStyleSidebar)
	sidebar.list_widget.SetItems(sidebar.list_items())
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

func (sidebar *Sidebar[T]) sidebar_items() []widget.ListItem[struct{}] {
	if len(sidebar.items) == 0 {
		sidebar.items = []widget.ListItem[struct{}]{
			{
				Header:       true,
				Content:      &sidebar.header,
				Unselectable: true,
			},
			{
				Border: true,
			},
		}
	}
	return sidebar.items[2:]
}

func (sidebar *Sidebar[T]) update_sidebar_items(items []SidebarItem[T]) {
	sidebar_items := sidebar.sidebar_items()
	for i := range sidebar_items {
		if sidebar_items[i].Content == nil {
			sidebar_items[i].Content = &sidebar_item_widget[T]{}
		}

		item_widget := sidebar_items[i].Content.(*sidebar_item_widget[T])
		item_widget.SetSideBarItem(items[i])
	}
}

func (sidebar *Sidebar[T]) SetSidebarItems(items []SidebarItem[T]) {
	sidebar_items := sidebar.sidebar_items()
	l := len(sidebar_items)
	if l > len(items) {
		sidebar.items = sidebar.items[:2]
		sidebar.items = append(sidebar.items, sidebar_items[:len(items)]...)
	} else if l < len(items) {
		diff := len(items) - l
		list := make([]widget.ListItem[struct{}], diff)
		sidebar.items = append(sidebar.items, list...)
	}
	sidebar.update_sidebar_items(items)
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

// Path returns the current directory path.
func (sidebar *Sidebar[T]) Path() string {
	return sidebar.header.Path()
}

func (sidebar *Sidebar[T]) OnItemSelect(f func(ctx *gui.Context, index int)) {
	sidebar.list_widget.OnItemSelected(func(ctx *gui.Context, index int) {
		f(ctx, index-2)
	})
}

func (sidebar *Sidebar[T]) SelectItemByIndex(index int) {
	sidebar.list_widget.SelectItemByIndex(index + 2)
}