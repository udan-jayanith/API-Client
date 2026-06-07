package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type popup_content struct {
	gui.DefaultWidget

	heading      CommonWidgets.Description
	input_widget CommonWidgets.TextInputWithContextMenu
	button       widget.Button
}

func (content *popup_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&content.heading)
	adder.AddWidget(&content.input_widget)
	content.button.SetType(widget.ButtonTypePrimary)
	adder.AddWidget(&content.button)
	return nil
}

func (content *popup_content) layout(ctx *gui.Context) gui.LinearLayout {
	gap := basic.Gap(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       gap,
		Padding:   basic.NewPadding(widget.UnitSize(ctx) / 3),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &content.heading,
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       gap,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &content.input_widget,
							Size:   gui.FixedSize(widget.UnitSize(ctx) * 8),
						},
						{
							Widget: &content.button,
						},
					},
				},
			},
		},
	}
	return layout
}

func (content *popup_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	content.layout(ctx).LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (content *popup_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return content.layout(ctx).Measure(ctx, constraints)
}

func (content *popup_content) SetButtonText(text string) {
	content.button.SetText(text)
}

func (content *popup_content) OnButtonClick(f func(context *gui.Context)) {
	content.button.OnDown(f)
}

func (content *popup_content) SetHeading(heading string) {
	content.heading.SetDescription(heading)
}

func (content *popup_content) Value() string {
	return content.input_widget.Value()
}
func (content *popup_content) FocusInput(ctx *gui.Context) {
	ctx.SetFocused(&content.input_widget, true)
}

// Adapted from gui ContextMenuArea
type popup_skeleton_area struct {
	gui.DefaultWidget

	content   popup_content
	popup     widget.Popup
	on_result func(value string)
}

func (c *popup_skeleton_area) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if c.popup.IsOpen() {
		c.content.OnButtonClick(func(context *gui.Context) {
			c.popup.SetOpen(false)
			if c.on_result != nil {
				c.on_result(c.content.Value())
			}
		})
		c.popup.SetCloseByClickingOutside(true)
		c.popup.SetContent(&c.content)
		adder.AddWidget(&c.popup)
	}
	c.popup.SetModal(false)
	ctx.SetButtonInputReceptive(c, c.popup.IsOpen())
	return nil
}

func (c *popup_skeleton_area) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	popup_size := c.content.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()
	position := image.Pt(b.Min.X, b.Max.Y+basic.Gap(ctx))
	layouter.LayoutWidget(&c.popup, image.Rectangle{
		Min: position,
		Max: position.Add(popup_size),
	})
}

func (c *popup_skeleton_area) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	// Returning zero keeps a ContextMenuArea from contributing an unexpected size when used as an item
	// in a layout such as LinearLayout, which would otherwise pick up the inherited DefaultWidget size.
	return image.Point{}
}

func (c *popup_skeleton_area) FocusInput(ctx *gui.Context) {
	c.content.FocusInput(ctx)
}

func (c *popup_skeleton_area) SetOpen(open bool) {
	c.popup.SetOpen(open)
}

func (c *popup_skeleton_area) IsOpen() bool {
	return c.popup.IsOpen()
}

func (c *popup_skeleton_area) SetHeading(heading string) {
	c.content.SetHeading(heading)
}

func (c *popup_skeleton_area) SetButtonText(text string) {
	c.content.SetButtonText(text)
}

func (c *popup_skeleton_area) OnClose(f func(context *gui.Context, reason widget.PopupCloseReason)) {
	c.popup.OnClose(f)
}

func (c *popup_skeleton_area) OnResult(fn func(value string)) {
	c.on_result = fn
}
