package home

import (
	"API-Client/basic"
	"API-Client/icons"
	"image"

	opener "codeberg.org/udan-jayanith/Opener"
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type footer_widget struct {
	gui.DefaultWidget

	update_icon       *ebiten.Image
	check_for_updates widget.Button

	repo_link icons.Icon
	stars     struct {
		star_icon  icons.Icon
		star_count widget.Text
	}
	app_version widget.Text
}

func (fw *footer_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if fw.update_icon == nil {
		fw.update_icon = icons.Store.Open("update")
	}
	fw.check_for_updates.SetIcon(fw.update_icon)
	fw.check_for_updates.SetText("check for updates")
	adder.AddWidget(&fw.check_for_updates)

	icon_size := widget.LineHeight(ctx)
	icon_size -= (icon_size / 3)

	fw.repo_link.SetIcon("github")
	fw.repo_link.SetSize(icon_size)
	fw.repo_link.OnClick(func() {
		opener.Open("https://github.com/udan-jayanith/API-Client")
	})
	adder.AddWidget(&fw.repo_link)

	fw.stars.star_icon.SetIcon("star")
	fw.stars.star_icon.SetSize(icon_size)
	adder.AddWidget(&fw.stars.star_icon)

	fw.stars.star_count.SetValue("0")
	fw.stars.star_count.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&fw.stars.star_count)

	fw.app_version.SetValue("v0.0")
	fw.app_version.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&fw.app_version)
	return nil
}

func (fw *footer_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4

	footer_end_layout := gui.LinearLayoutItem{
		Layout: gui.LinearLayout{
			Direction: gui.LayoutDirectionHorizontal,
			Gap:       size * 2,
			Items: []gui.LinearLayoutItem{
				{
					Widget: &fw.repo_link,
				},
				{
					Layout: gui.LinearLayout{
						Direction: gui.LayoutDirectionHorizontal,
						Gap:       size / 2,
						Items: []gui.LinearLayoutItem{
							{
								Widget: &fw.stars.star_icon,
							},
							{
								Widget: &fw.stars.star_count,
							},
						},
					},
				},
				{
					Widget: &fw.app_version,
				},
			},
		},
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       size * 2,
		Padding:   basic.NewPadding(size),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &fw.check_for_updates,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						footer_end_layout,
						{
							Size: gui.FlexibleSize(1),
						},
					},
				},
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (fw *footer_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	gap := u / 4
	padding := basic.NewPadding(u / 4)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += fw.check_for_updates.Measure(ctx, gui.Constraints{}).X
		point.X += fw.repo_link.Measure(ctx, gui.Constraints{}).X
		point.X += fw.app_version.Measure(ctx, gui.Constraints{}).X

		point.X += fw.stars.star_icon.Measure(ctx, gui.Constraints{}).X
		point.X += fw.stars.star_count.Measure(ctx, gui.Constraints{}).X
		point.X += gap / 2
		point.X += gap * 3
		point.X += padding.Start + padding.End
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += padding.Top + padding.Bottom + widget.UnitSize(ctx)
	}

	return point
}
