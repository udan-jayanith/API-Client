package http_widget

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"image"

	"Zbolt/pages/request/http/internal/tooltip"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	httpheaders "gitlab.com/j.udanjayanith/http-headers"
)

type http_header_description struct {
	gui.DefaultWidget

	heading           basicwidget.Text
	hr                CommonWidgets.HorizontalLine
	description_panel basicwidget.Panel
	description       gui.WidgetWithPadding[*basicwidget.Text]
	refrence          CommonWidgets.Hyperlink
}

func (w *http_header_description) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.heading.SetBold(true)
	w.heading.SetSelectable(true)
	w.heading.SetEllipsisString("...")
	adder.AddWidget(&w.heading)
	adder.AddWidget(&w.hr)

	w.description.SetPadding(basic.NewPadding(0, basic.BorderRadius(ctx), 0, 0))
	description := w.description.Widget()
	description.SetMultiline(true)
	description.SetWrapMode(basicwidget.WrapModeNormal)
	description.SetSelectable(true)
	w.description_panel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)
	w.description_panel.SetStyle(basicwidget.PanelStyleDefault)
	w.description_panel.SetContent(&w.description)
	adder.AddWidget(&w.description_panel)

	w.refrence.SetMultiline(true)
	w.refrence.SetWrapMode(basicwidget.WrapModeAnywhere)
	adder.AddWidget(&w.refrence)
	return nil
}

func (w *http_header_description) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(basic.BorderRadius(ctx))
}

func (w *http_header_description) gap(ctx *gui.Context) int {
	return basic.Gap(ctx) / 2
}

func (w *http_header_description) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   w.padding(ctx),
		Gap:       w.gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.heading,
			},
			{
				Widget: &w.hr,
			},
			{
				Widget: &w.description_panel,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &w.refrence,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *http_header_description) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := basicwidget.UnitSize(ctx)
	max_size := image.Pt(u*12, u*8)

	size := image.Point{}
	size.X = max_size.X
	constraints = gui.FixedWidthConstraints(max_size.X)
	size.Y += w.heading.Measure(ctx, constraints).Y
	size.Y += w.hr.Measure(ctx, constraints).Y
	size.Y += w.description.Measure(ctx, constraints).Y
	size.Y += w.refrence.Measure(ctx, constraints).Y

	padding := w.padding(ctx)
	size.Y += padding.Top + padding.Bottom
	size.Y += w.gap(ctx) * 3

	if size.Y >= max_size.Y {
		return max_size
	}
	return size
}

func (w *http_header_description) Set(heading, description, refrence string) {
	w.heading.SetValue(heading)
	w.description.Widget().SetValue(description)
	w.refrence.SetRefrence(refrence)
}

type HttpHeaderTable struct {
	tooltip tooltip.TooltipArea
	header  *httpheaders.Header

	tooltip_content http_header_description
	hover_bounds    image.Rectangle
	CommonWidgets.AttributeTable
}

func (w *HttpHeaderTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.AttributeTable.OnHover(func(ctx *gui.Context, t string, widget *CommonWidgets.EditableText, widget_bounds *gui.WidgetBounds) {
		if t == "value" {
			return
		}

		h_name := widget.Value()
		h := httpheaders.Search(h_name)
		if h == nil || h.HeaderName != h_name {
			w.header = nil
			return
		}
		w.header = h
		w.hover_bounds = widget_bounds.Bounds()
		size := widget.Measure(ctx, gui.FixedWidthConstraints(w.hover_bounds.Dx()))
		w.hover_bounds.Max.X = w.hover_bounds.Min.X + size.X
		w.hover_bounds.Max.Y = w.hover_bounds.Min.Y + size.Y
	})
	w.AttributeTable.Build(ctx, adder)
	if w.header != nil {
		w.tooltip_content.Set(w.header.HeaderName, string(w.header.Description), w.header.Reference)
		w.tooltip.SetContent(&w.tooltip_content)
		adder.AddWidget(&w.tooltip)
	}
	return nil
}

func (w *HttpHeaderTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	w.AttributeTable.Layout(ctx, widgetBounds, layouter)
	if w.header != nil {
		layouter.LayoutWidget(&w.tooltip, w.hover_bounds)
	}
}
