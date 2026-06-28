package http_widget

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_file_content struct {
	gui.DefaultWidget

	close_icon    icons.Icon
	file_selected bool
	filename      CommonWidgets.TextWithTooltip
	select_btn    widget.Button
}

func (content *request_file_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	content.close_icon.SetIcon("close")
	adder.AddWidget(&content.close_icon)

	if content.file_selected {
		// TODO: when the filename is clicked open the directory where the file is at.
		content.filename.SetTooltip("File name")
		adder.AddWidget(&content.filename)
	}

	content.select_btn.SetText("Select file")
	adder.AddWidget(&content.select_btn)

	return nil
}

func (content *request_file_content) layout(ctx *gui.Context) gui.LinearLayout {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       basic.Gap(ctx),
		Padding:   basic.NewPadding(basic.BorderRadius(ctx)),
		Items:     make([]gui.LinearLayoutItem, 0, 5),
	}

	layout.Items = append(layout.Items, []gui.LinearLayoutItem{
		{
			Layout: gui.LinearLayout{
				Direction: gui.LayoutDirectionHorizontal,
				Items: []gui.LinearLayoutItem{
					{
						Size: gui.FlexibleSize(1),
					},
					{
						Widget: &content.close_icon,
					},
				},
			},
		},
		{
			Size: gui.FlexibleSize(1),
		},
	}...)

	if content.file_selected {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &content.filename,
		})
	}

	layout.Items = append(layout.Items, []gui.LinearLayoutItem{
		{
			Widget: &content.select_btn,
		},
		{
			Size: gui.FlexibleSize(1),
		},
	}...)

	return layout
}

func (content *request_file_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	content.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (content *request_file_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return content.layout(ctx).Measure(ctx, constraints)
}

func (content *request_file_content) OnFileSelect(func(ctx *gui.Context)) {

}

func (content *request_file_content) OnFileDirectoryOpen(func(ctx *gui.Context)) {

}

func (content *request_file_content) SetFileSelected(selected bool) {
	content.file_selected = true
}

func (content *request_file_content) FileSelected() bool {
	return content.file_selected
}

func (content *request_file_content) OnClose(func(ctx *gui.Context)) {
	
}
