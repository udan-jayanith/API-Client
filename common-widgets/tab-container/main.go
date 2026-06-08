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
	background basicwidget.Background

	tab_container CommonWidgets.TabContainer

	text   basicwidget.Text
	button basicwidget.Button
	toggle basicwidget.Toggle
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	context.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&r.background)

	r.text.SetValue("Hello world")
	r.button.SetText("Button")
	r.tab_container.SetItems([]CommonWidgets.TabContainerItem{
		{
			TabItem: CommonWidgets.TabItem{
				Text:  "text",
				Value: "text",
			},
			Widget: &r.text,
		},
		{
			TabItem: CommonWidgets.TabItem{
				Text:  "button",
				Value: "btn",
			},
			Widget: &r.button,
		},
		{
			TabItem: CommonWidgets.TabItem{
				Text:  "toggle",
				Value: "toggle",
			},
			Widget: &r.toggle,
		},
	})
	adder.AddWidget(&r.tab_container)
	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&r.background, widgetBounds.Bounds())
	layouter.LayoutWidget(&r.tab_container, widgetBounds.Bounds())
}

func main() {
	r := &Root{}
	op := &guigui.RunOptions{
		Title: "Tab-Container",
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(r, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
