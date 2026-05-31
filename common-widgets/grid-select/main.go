package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	CommonWidgets "Zbolt/common-widgets"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Root struct {
	guigui.DefaultWidget
	background  basicwidget.Background
	grid_select CommonWidgets.GridSelect[struct{}]
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&r.background)

	r.grid_select.SetItems([]CommonWidgets.GridSelectItem[struct{}]{
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
	})
	adder.AddWidget(&r.grid_select)
	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&r.background, widgetBounds.Bounds())
	layouter.LayoutWidget(&r.grid_select, widgetBounds.Bounds())
}

func main() {
	r := &Root{}
	op := &guigui.RunOptions{
		Title: "Grid-Select",
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(r, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
