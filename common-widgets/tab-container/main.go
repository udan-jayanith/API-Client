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

	tab_container CommonWidgets.TabContainer[string]

	text   basicwidget.Text
	button basicwidget.Button
	toggle basicwidget.Toggle

	is_set bool
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	context.SetPreferredColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&r.background)

	r.text.SetValue("Hello world")
	r.button.SetText("Button")
	if !r.is_set {
		r.tab_container.SetItems([]CommonWidgets.TabContainerItem[string]{
			{
				TabItem: CommonWidgets.TabItem[string]{
					Text:  "text",
					Value: "text",
				},
				Widget: &r.text,
			},
			{
				TabItem: CommonWidgets.TabItem[string]{
					Text:  "button",
					Value: "btn",
				},
				Widget: &r.button,
			},
			{
				TabItem: CommonWidgets.TabItem[string]{
					Text:  "toggle",
					Value: "toggle",
				},
				Widget: &r.toggle,
			},
		})
		r.is_set = true
	}
	r.tab_container.SetClosable(true)
	r.tab_container.OnClose(func(closed CommonWidgets.TabItemContainer[string]) {
		println("closed", closed.Item.Text)
	})
	r.tab_container.OnSwap(func(from, to CommonWidgets.TabItemContainer[string]) {
		println("swaped")
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
