package home

import (
	"API-Client/basic"
	"API-Client/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type HomePage struct {
	gui.DefaultWidget

	background      widget.Background
	recently_opened gui.WidgetWithPadding[*widget.Text]
	sidebar         widget.List[struct{}]

	zbolt_icon                                            icons.Icon
	open_button, new_project_button, quick_request_button widget.Button
	footer_widget                                         footer_widget
}

func (wp *HomePage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&wp.background)

	wp.recently_opened.Widget().SetValue("Recently opened")
	wp.recently_opened.Widget().SetScale(1.2)
	wp.recently_opened.SetPadding(gui.Padding{
		Start: widget.UnitSize(ctx) / 2,
	})
	adder.AddWidget(&wp.recently_opened)

	wp.sidebar.SetItemsByStrings([]string{"Download manager", "Retailers site", "Update manager", "Video platform"})
	wp.sidebar.SetStyle(widget.ListStyleSidebar)
	adder.AddWidget(&wp.sidebar)

	wp.zbolt_icon.SetIcon("large-icons/zbolt-passtrough")
	size := widget.UnitSize(ctx) * 14
	wp.zbolt_icon.SetSize(size)
	adder.AddWidget(&wp.zbolt_icon)

	wp.open_button.SetText("Open project")
	adder.AddWidget(&wp.open_button)

	wp.new_project_button.SetText("New project")
	adder.AddWidget(&wp.new_project_button)

	wp.quick_request_button.SetText("Quick request")
	adder.AddWidget(&wp.quick_request_button)

	adder.AddWidget(&wp.footer_widget)
	return nil
}

func (wp *HomePage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	bounds := widgetBounds.Bounds()
	layouter.LayoutWidget(&wp.background, bounds)

	u := widget.UnitSize(ctx)
	padding := u / 4
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Padding:   basic.NewPadding(padding, padding, padding, 0),
		Items: []gui.LinearLayoutItem{
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionVertical,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &wp.recently_opened,
						},
						{
							Widget: &wp.sidebar,
							Size:   gui.FlexibleSize(1),
						},
					},
				},
				Size: gui.FlexibleSize(1),
			},
		},
	}

	main_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Size: gui.FlexibleSize(1),
			},
			// zbolt icon
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &wp.zbolt_icon,
						},
						{
							Size: gui.FlexibleSize(1),
						},
					},
				},
			},
			{
				Size: gui.FlexibleSize(1),
			},
			// buttons
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       u / 2,
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &wp.open_button,
						},
						{
							Widget: &wp.new_project_button,
						},
						{
							Widget: &wp.quick_request_button,
						},
						{
							Size: gui.FlexibleSize(1),
						},
					},
				},
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &wp.footer_widget,
			},
		},
	}

	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Size:   gui.FlexibleSize(4),
		Layout: main_layout,
	})
	layout.LayoutWidgets(ctx, bounds, layouter)
}
