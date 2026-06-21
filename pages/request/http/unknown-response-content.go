package http_widget

import (
	"Zbolt/basic"
	"Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type unknow_response_content struct {
	gui.DefaultWidget

	msg                              widget.Text
	save_as_btn, open_externally_btn widget.Button
	save_as_ico, open_externally_ico *ebiten.Image
}

func (content *unknow_response_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	content.msg.SetMultiline(true)
	content.msg.SetEllipsisString("...")
	content.msg.SetHorizontalAlign(widget.HorizontalAlignCenter)
	content.msg.SetWrapMode(widget.WrapModeNormal)
	adder.AddWidget(&content.msg)

	if content.open_externally_ico == nil {
		content.open_externally_ico = icons.Store.Open("open-in")
	}
	content.open_externally_btn.SetText("Open Externally")
	content.open_externally_btn.SetIcon(content.open_externally_ico)
	adder.AddWidget(&content.open_externally_btn)

	if content.save_as_ico == nil {
		content.save_as_ico = icons.Store.Open("save-as")
	}
	content.save_as_btn.SetText("Save As")
	content.save_as_btn.SetIcon(content.save_as_ico)
	adder.AddWidget(&content.save_as_btn)
	return nil
}

func (content *unknow_response_content) layout(ctx *gui.Context) gui.LinearLayout {
	btn_w := max(content.open_externally_btn.Measure(ctx, gui.Constraints{}).X, content.save_as_btn.Measure(ctx, gui.Constraints{}).X)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       basic.Gap(ctx),
		Padding:   basic.NewPadding(basic.BorderRadius(ctx)),
		Items: []gui.LinearLayoutItem{
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &content.msg,
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       basic.Gap(ctx),
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &content.save_as_btn,
							Size:   gui.FixedSize(btn_w),
						},
						{
							Widget: &content.open_externally_btn,
							Size:   gui.FixedSize(btn_w),
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
		},
	}
	return layout
}

func (content *unknow_response_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	content.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (content *unknow_response_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return content.layout(ctx).Measure(ctx, constraints)
}

func (content *unknow_response_content) OnOpenExternally(f func(context *gui.Context)) {
	content.open_externally_btn.OnDown(f)
}

func (content *unknow_response_content) OnSaveAs(f func(context *gui.Context)) {
	content.save_as_btn.OnDown(f)
}
