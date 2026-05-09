package http_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_header_widget struct {
	gui.DefaultWidget

	status, response_time, content_lenght, proto CommonWidgets.TextWithTooltip
}

func (w *response_header_widget) clear() {
	w.status.SetValue("")
	w.response_time.SetValue("")
	w.content_lenght.SetValue("")
	w.proto.SetValue("")
}

func (rhw *response_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rhw.status.SetTabular(true)
	rhw.status.SetTooltip("HTTP response status")
	if rhw.status.Value() == "" {
		rhw.status.SetValue("None")
	}
	adder.AddWidget(&rhw.status)

	rhw.response_time.SetTabular(true)
	rhw.response_time.SetTooltip("Response time")
	if rhw.response_time.Value() == "" {
		rhw.response_time.SetValue("0ms")
	}
	adder.AddWidget(&rhw.response_time)

	rhw.content_lenght.SetTabular(true)
	rhw.content_lenght.SetTooltip("Content length")
	if rhw.content_lenght.Value() == "" {
		rhw.content_lenght.SetValue("0B")
	}
	adder.AddWidget(&rhw.content_lenght)

	rhw.proto.SetTabular(true)
	rhw.proto.SetTooltip("Used HTTP protocol version")
	if rhw.proto.Value() == "" {
		rhw.proto.SetValue("HTTP vX.X")
	}
	adder.AddWidget(&rhw.proto)
	return nil
}

func (rhw *response_header_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 4)
}

func (rhw *response_header_widget) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (rhw *response_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       rhw.gap(ctx),
		Padding:   rhw.padding(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rhw.status,
			},
			{
				Widget: &rhw.response_time,
			},
			{
				Widget: &rhw.content_lenght,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rhw.proto,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rhw *response_header_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	padding := u / 2

	point := image.Pt(padding+gap*4, padding)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		status_measurements := rhw.status.Measure(ctx, constraints)
		point.X += status_measurements.X

		response_time_measurements := rhw.response_time.Measure(ctx, constraints)
		point.X += response_time_measurements.X

		size := rhw.content_lenght.Measure(ctx, constraints)
		point.X += size.X

		proto_measurements := rhw.proto.Measure(ctx, constraints)
		point.X += proto_measurements.X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += widget.LineHeight(ctx)
	}

	return point
}
