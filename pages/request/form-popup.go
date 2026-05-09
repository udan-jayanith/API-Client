package request_page

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type popup_form_content struct {
	gui.DefaultWidget

	field_widget  CommonWidgets.Description
	input_widget  CommonWidgets.TextInputWithContextMenu
	button_widget widget.Button
}

func (content *popup_form_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if content.field_widget.Description() != "" {
		adder.AddWidget(&content.field_widget)
	}

	adder.AddWidget(&content.input_widget)

	content.button_widget.SetType(widget.ButtonTypePrimary)
	adder.AddWidget(&content.button_widget)
	return nil
}

func (content *popup_form_content) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (content *popup_form_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := content.gap(ctx)

	horizontal_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       gap,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &content.input_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &content.button_widget,
			},
		},
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       gap / 2,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &content.field_widget,
			},
			{
				Layout: horizontal_layout,
			},
		},
	}

	if content.field_widget.Description() == "" {
		layout = horizontal_layout
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (content *popup_form_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	gap := content.gap(ctx)

	point.X = u * 10
	if content.field_widget.Description() != "" {
		point.Y += content.field_widget.Measure(ctx, gui.Constraints{}).Y + gap/2
	}

	point.Y += max(content.button_widget.Measure(ctx, gui.Constraints{}).Y, content.input_widget.Measure(ctx, gui.Constraints{}).Y)
	return point
}

type folder_create_popup struct {
	gui.DefaultWidget

	popup_widget      widget.Popup
	popup_content     popup_form_content
	padding_widget    CommonWidgets.WidgetWithPadding[*popup_form_content]
	on_button_clicked func(ctx *gui.Context, value string)
}

func (sfp *folder_create_popup) Build(ctx *gui.Context, adder *gui.ChildAdder) error {

	sfp.popup_content.button_widget.OnUp(func(_ *gui.Context) {
		sfp.popup_widget.SetOpen(false)
		if sfp.on_button_clicked != nil {
			sfp.on_button_clicked(ctx, sfp.popup_content.input_widget.Value())
		}
	})

	sfp.padding_widget.SetWidget(&sfp.popup_content)
	sfp.padding_widget.SetPadding(basic.NewPadding(widget.UnitSize(ctx) / 3))

	sfp.popup_widget.SetContent(&sfp.padding_widget)
	sfp.popup_widget.SetAnimated(true)
	sfp.popup_widget.SetCloseByClickingOutside(true)
	adder.AddWidget(&sfp.popup_widget)
	return nil
}

func (sfp *folder_create_popup) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sfp.popup_widget, widgetBounds.Bounds())
}

func (sfp *folder_create_popup) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sfp.padding_widget.Measure(ctx, constraints)
}

func (sfp *folder_create_popup) SetButtonText(text string) {
	sfp.popup_content.button_widget.SetText(text)
}

func (sfp *folder_create_popup) SetFieldValue(text string) {
	sfp.popup_content.field_widget.SetDescription(text)
}

func (sfp *folder_create_popup) ClearInput() {
	sfp.popup_content.input_widget.SetValue("")
}

func (sfp *folder_create_popup) SetOpen(open bool) {
	sfp.popup_widget.SetOpen(open)
}

func (sfp *folder_create_popup) OnButtonClicked(fn func(ctx *gui.Context, value string)) {
	sfp.on_button_clicked = fn
}
