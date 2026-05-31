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
	context.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&r.background)

	index, _ := r.grid_select.SelectedItemIndex()
	r.grid_select.SetItems([]CommonWidgets.GridSelectItem[struct{}]{
		{
			IconName: "copy",
			Title:    "Item lore lorem f gag ag a a h hah aha h aahahahahahahaha",
		},
		{
			IconName: "add-box",
			Title:    "Item",
		},
		{
			IconName: "large-icons/http",
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
	r.grid_select.SelectItemByIndex(index)
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
