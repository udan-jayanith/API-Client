package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"

	message_model "Zbolt/message-model"
	home "Zbolt/pages/home"
	request_page "Zbolt/pages/request"
	//"runtime/pprof"
)

type Root struct {
	gui.DefaultWidget
	background     basicwidget.Background
	menubar_widget basicwidget.Menubar[struct{}]

	welcome_page_widget home.HomePage
	request_page_widget request_page.RequestPage
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.background)

	r.menubar_widget.SetItems([]basicwidget.MenubarItem{
		{
			Text: "Zbolt",
		},
		{
			Text: "Project",
		},
		{
			Text: "Logs",
		},
	})
	adder.AddWidget(&r.menubar_widget)
	adder.AddWidget(&r.request_page_widget)
	adder.AddWidget(&message_model.MessageModel)
	return nil
}

func (r *Root) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(&r.background, b)
	layouter.LayoutWidget(&message_model.MessageModel, widgetBounds.Bounds())

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &r.menubar_widget,
			},
			{
				Widget: &r.request_page_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

//go:embed icon.png
var zbolt_icon_bytes []byte

func main() {
	//os.Remove("./cpu.prof")
	//f, _ := os.Create("cpu.prof")
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()
	//go tool pprof -http=:8080 cpu.prof

	zebolt_icon, _, err := image.Decode(bytes.NewReader(zbolt_icon_bytes))
	if err != nil {
		log.Fatal(err.Error())
	}
	ebiten.SetWindowIcon([]image.Image{zebolt_icon})
	op := &gui.RunOptions{
		Title:         "Zbolt",
		WindowMinSize: image.Pt(800, 544),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}

	message_model.Show("Hello world", message_model.Alert, nil)
	if err := gui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
