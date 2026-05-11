package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	"Zbolt/icons"
	"Zbolt/pages/request/requests-handler"
	"image"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type sidebar_item_type_card struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text

	Text, Icon_name string
}

func (sitc *sidebar_item_type_card) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sitc.icon_widget.SetIcon("large-icons/" + sitc.Icon_name)
	u := widget.UnitSize(ctx)
	sitc.icon_widget.SetSize(u*2)
	adder.AddWidget(&sitc.icon_widget)

	sitc.text_widget.SetValue(sitc.Text)
	sitc.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	sitc.text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	adder.AddWidget(&sitc.text_widget)

	return nil
}

func (sitc *sidebar_item_type_card) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Gap:       size,
		Padding:   basic.NewPadding(size),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{
				Widget: &sitc.icon_widget,
			},
			{
				Widget: &sitc.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sitc *sidebar_item_type_card) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	size := widget.UnitSize(ctx) / 4

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = sitc.icon_widget.Measure(ctx, constraints).Y
		point.Y += sitc.text_widget.Measure(ctx, constraints).Y
		point.Y += size * 3
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = point.Y
	}

	return point
}

type request_name_inputs_widget struct {
	gui.DefaultWidget

	text_widget              widget.Text
	input_widget             CommonWidgets.TextInputWithContextMenu
	create_button_widget     widget.Button
	on_create_button_clicked func(name string)
}

func (rniw *request_name_inputs_widget) OnCreateButtonClicked(fn func(name string)) {
	rniw.on_create_button_clicked = fn
}

func (rniw *request_name_inputs_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rniw.text_widget.SetValue("Enter request name")
	rniw.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	rniw.text_widget.SetHorizontalAlign(widget.HorizontalAlignStart)
	rniw.text_widget.SetScale(1)
	adder.AddWidget(&rniw.text_widget)

	adder.AddWidget(&rniw.input_widget)

	rniw.create_button_widget.SetText("Create")
	rniw.create_button_widget.SetType(widget.ButtonTypePrimary)
	if rniw.on_create_button_clicked != nil {
		rniw.create_button_widget.OnDown(func(context *gui.Context) {
			rniw.on_create_button_clicked(rniw.input_widget.Value())
		})
	}
	adder.AddWidget(&rniw.create_button_widget)
	return nil
}

func (rniw *request_name_inputs_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rniw.text_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       widget.UnitSize(ctx) / 4,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &rniw.input_widget,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &rniw.create_button_widget,
						},
					},
				},
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rniw *request_name_inputs_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = rniw.text_widget.Measure(ctx, constraints).Y + u
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = rniw.input_widget.Measure(ctx, constraints).X + rniw.create_button_widget.Measure(ctx, constraints).X
	}

	return point
}

type sidebar_item_types_panel struct {
	gui.DefaultWidget

	http, websocket, graphql, grpc gui.WidgetWithSize[*widget.Button]
	selected_request_type          requests_handler.RequestType

	line CommonWidgets.HorizontalLine

	select_type_text_widget widget.Text
	request_name_input      request_name_inputs_widget
	on_create_clicked       func(request *requests_handler.Request)
}

func (sitp *sidebar_item_types_panel) Clear() {
	sitp.selected_request_type = requests_handler.HTTP
	sitp.request_name_input.input_widget.SetValue("")
}

func (sitp *sidebar_item_types_panel) OnCreateButtonClicked(fn func(request *requests_handler.Request)) {
	sitp.on_create_clicked = fn
}

func (sitp *sidebar_item_types_panel) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sitp.select_type_text_widget.SetValue("Select request type")
	sitp.select_type_text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	sitp.select_type_text_widget.SetHorizontalAlign(widget.HorizontalAlignStart)
	sitp.select_type_text_widget.SetScale(1)
	adder.AddWidget(&sitp.select_type_text_widget)

	u := widget.UnitSize(ctx)
	item_size := image.Pt(u*4, u*4)

	sitp.http.SetFixedSize(item_size)
	http := sitp.http.Widget()
	http.SetContent(&sidebar_item_type_card{
		Text:      "HTTP",
		Icon_name: "http",
	})
	http.OnDown(func(context *gui.Context) {
		sitp.selected_request_type = requests_handler.HTTP
	})
	adder.AddWidget(&sitp.http)

	sitp.websocket.SetFixedSize(item_size)
	websocket := sitp.websocket.Widget()
	websocket.SetContent(&sidebar_item_type_card{
		Text:      "Websocket",
		Icon_name: "websocket",
	})
	websocket.OnDown(func(context *gui.Context) {
		sitp.selected_request_type = requests_handler.Websocket
	})

	adder.AddWidget(&sitp.websocket)

	sitp.graphql.SetFixedSize(item_size)
	graphql := sitp.graphql.Widget()
	graphql.SetContent(&sidebar_item_type_card{
		Text:      "GraphQL",
		Icon_name: "graphql",
	})
	graphql.OnDown(func(context *gui.Context) {
		sitp.selected_request_type = requests_handler.GraphQL
	})
	adder.AddWidget(&sitp.graphql)

	sitp.grpc.SetFixedSize(item_size)
	grpc := sitp.grpc.Widget()
	grpc.SetContent(&sidebar_item_type_card{
		Text:      "gRPC",
		Icon_name: "grpc",
	})
	grpc.OnDown(func(context *gui.Context) {
		sitp.selected_request_type = requests_handler.Grpc
	})
	adder.AddWidget(&sitp.grpc)

	http.SetType(widget.ButtonTypeNormal)
	websocket.SetType(widget.ButtonTypeNormal)
	graphql.SetType(widget.ButtonTypeNormal)
	grpc.SetType(widget.ButtonTypeNormal)

	switch sitp.selected_request_type {
	case requests_handler.HTTP:
		http.SetType(widget.ButtonTypePrimary)
	case requests_handler.Websocket:
		websocket.SetType(widget.ButtonTypePrimary)
	case requests_handler.GraphQL:
		graphql.SetType(widget.ButtonTypePrimary)
	case requests_handler.Grpc:
		grpc.SetType(widget.ButtonTypePrimary)
	default:
		log.Fatal("Unknown request type selected")
	}

	adder.AddWidget(&sitp.line)

	sitp.request_name_input.OnCreateButtonClicked(func(name string) {
		req := requests_handler.NewRequest(sitp.selected_request_type, name)
		sitp.on_create_clicked(&req)
	})
	adder.AddWidget(&sitp.request_name_input)
	return nil
}

func (sitp *sidebar_item_types_panel) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Padding:   basic.NewPadding(size * 2),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{
				Widget: &sitp.select_type_text_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Gap:       size,
					Direction: gui.LayoutDirectionHorizontal,

					Items: []gui.LinearLayoutItem{
						{
							Widget: &sitp.http,
						},
						{
							Widget: &sitp.websocket,
						},
						{
							Widget: &sitp.graphql,
						},
						{
							Widget: &sitp.grpc,
						},
					},
				},
			},
			{
				Size: gui.FixedSize(size * 2),
			},
			{
				Widget: &sitp.line,
			},
			{
				Size: gui.FixedSize(size),
			},
			{
				Widget: &sitp.request_name_input,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sitp *sidebar_item_types_panel) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	var point image.Point
	gap, padding := u/4, u/2

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = sitp.select_type_text_widget.Measure(ctx, constraints).Y
		point.Y += sitp.request_name_input.Measure(ctx, constraints).Y
		point.Y += (u * 4) + gap + gap*2 + padding*3 + sitp.line.Measure(ctx, gui.Constraints{}).Y
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = (u*4)*4 + gap*3 + padding*2
	}

	return point
}
