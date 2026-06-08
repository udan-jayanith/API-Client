package request_page

import (
	"Zbolt/basic"
	CommonWidgets "Zbolt/common-widgets"
	requests_handler "Zbolt/pages/request/requests-handler"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_create_panel struct {
	gui.DefaultWidget

	select_request_type widget.Text
	grid_selector       CommonWidgets.GridSelect[requests_handler.RequestType]
	is_grid_items_set   bool

	hr                 CommonWidgets.HorizontalLine
	enter_request_name widget.Text
	request_name_input CommonWidgets.TextInputWithContextMenu
	create             widget.Button
	on_item_create     OnRequestItemCreateFunc
}

func (panel *request_create_panel) set_grid_items() {
	if panel.is_grid_items_set {
		return
	}
	panel.is_grid_items_set = true
	index, _ := panel.grid_selector.SelectedItemIndex()
	panel.grid_selector.SetItems([]CommonWidgets.GridSelectItem[requests_handler.RequestType]{
		{
			Title:    "HTTP",
			IconName: "large-icons/http",
			Value:    requests_handler.HTTP,
		},
		{
			Title:    "Websocket",
			IconName: "large-icons/websocket",
			Value:    requests_handler.Websocket,
		},
		{
			Title:    "GraphQL",
			IconName: "large-icons/graphql",
			Value:    requests_handler.GraphQL,
		},
		{
			Title:    "Grpc",
			IconName: "large-icons/grpc",
			Value:    requests_handler.Grpc,
		},
	})
	panel.grid_selector.SelectItemByIndex(index)
}

func (panel *request_create_panel) path(ctx *gui.Context) string {
	path, ok := basic.GetEnvMust[string](ctx, panel, path_env)
	if !ok {
		panic("path_env failed")
	}
	return path
}

func (panel *request_create_panel) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	panel.select_request_type.SetValue("Select request type")
	adder.AddWidget(&panel.select_request_type)

	panel.set_grid_items()
	adder.AddWidget(&panel.grid_selector)

	adder.AddWidget(&panel.hr)
	panel.enter_request_name.SetValue("Enter request name")
	adder.AddWidget(&panel.enter_request_name)
	adder.AddWidget(&panel.request_name_input)

	panel.create.SetText("Create")
	panel.create.SetType(widget.ButtonTypePrimary)
	panel.create.OnDown(func(context *gui.Context) {
		defer panel.Clear()
		if panel.on_item_create == nil {
			return
		}

		item, ok := panel.grid_selector.SelectedItem()
		if !ok {
			panic("No items in grid select")
		}
		panel.on_item_create(ctx, panel.path(ctx), panel.request_name_input.Value(), item)
	})
	adder.AddWidget(&panel.create)
	return nil
}

func (panel *request_create_panel) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 2)
}

func (panel *request_create_panel) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	grid_selecttor_size := panel.grid_selector.RecommendedSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       basic.Gap(ctx),
		Padding:   panel.padding(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &panel.select_request_type,
			},
			{
				Widget: &panel.grid_selector,
				Size:   gui.FixedSize(grid_selecttor_size.Y),
			},
			{
				Widget: &panel.hr,
			},
			{
				Widget: &panel.enter_request_name,
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       basic.Gap(ctx),
					Items: []gui.LinearLayoutItem{
						{
							Widget: &panel.request_name_input,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &panel.create,
						},
					},
				},
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (panel *request_create_panel) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	panel.set_grid_items()

	var size image.Point
	grid_selecttor_size := panel.grid_selector.RecommendedSize(ctx)
	padding := panel.padding(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		size.X = w
	} else {
		size.X = grid_selecttor_size.X + padding.End + padding.Start
	}

	if h, ok := constraints.FixedHeight(); ok {
		size.Y = h
	} else {
		w := size.X - (padding.End + padding.Start)
		constraints := gui.FixedWidthConstraints(w)
		size.Y += panel.select_request_type.Measure(ctx, constraints).Y
		size.Y += grid_selecttor_size.Y
		size.Y += int(basic.LineWidth(ctx))
		size.Y += panel.enter_request_name.Measure(ctx, constraints).Y
		w -= panel.create.Measure(ctx, gui.Constraints{}).X
		size.Y += panel.request_name_input.Measure(ctx, gui.FixedWidthConstraints(w)).Y
		size.Y += basic.Gap(ctx) * 4
		size.Y += padding.Bottom + padding.Top
	}

	return size
}

func (panel *request_create_panel) Clear() {
	panel.set_grid_items()
	panel.grid_selector.SelectItemByIndex(0)
	panel.request_name_input.SetValue("")
}

func (panel *request_create_panel) OnCreate(fn OnRequestItemCreateFunc) {
	panel.on_item_create = fn
}
