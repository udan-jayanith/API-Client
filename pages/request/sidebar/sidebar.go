package sidebar

import (
	"Zbolt/basic"
	"image"

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
	sidebar_item_rename_popup_env gui.EnvKey = gui.GenerateEnvKey()
)

type Sidebar[T any] struct {
	gui.DefaultWidget

	add_folder_menu           widget.PopupMenu[struct{}]
	sidebar_item_context_menu widget.ContextMenuArea[struct{}]
	rename_item_menu          popup_skeleton_area

	sidebar_header sidebar_header_widget
	sidebar_items  []widget.ListItem[struct{}]
	list_widget    widget.List[struct{}]
	item           sidebar_item_widget[struct{}]

	on_sidebar_item_clicked func(ctx *gui.Context, sidebar_item SidebarItem[T])
	on_sidebar_item_rename  func(ctx *gui.Context, sidebar_item SidebarItem[T], new_name string)
	on_sidebar_item_delete  func(ctx *gui.Context, sidebar_item SidebarItem[T])
}

func (sidebar *Sidebar[T]) Env(ctx *gui.Context, key gui.EnvKey, source *gui.EnvSource) (any, bool) {
	if key == sidebar_item_contextmenu_env {
		return &sidebar.sidebar_item_context_menu, true
	} else if key == sidebar_item_rename_popup_env {
		return &sidebar.rename_item_menu, true
	}
	return nil, false
}

func (sidebar *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sidebar.list_widget.SetStyle(widget.ListStyleSidebar)
	if len(sidebar.sidebar_items) == 0 {
		sidebar.item.SetSideBarItem(SidebarItem[struct{}]{
			IconName: "search",
			Text:     "Item",
		})
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

// OnItemRename
// OnItemDelete
// OnSearch
// OnFolderCreate
// on_variable_panel_click func(ctx *gui.Context)
// on_create_request_click func(ctx *gui.Context)
// on_folder_create
