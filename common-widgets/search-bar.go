package CommonWidgets

import (
	"Zbolt/icons"
	"image"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

/*
type inline_popup_suggestions struct {
	gui.DefaultWidget

	suggestions widget.SegmentedControl[struct{}]
	x           int
}

func (sug *inline_popup_suggestions) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sug.suggestions.SetItems([]widget.SegmentedControlItem[struct{}]{
		{
			Text: "Hello",
		},
		{
			Text: "World",
		},
	})
	adder.AddWidget(&sug.suggestions)
	return nil
}

func (sug *inline_popup_suggestions) layout(ctx *gui.Context) gui.LinearLayout {
	// TODO: Improve mem allocation for this.
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		//Padding:   basic.NewPadding(basic.BorderRadius(ctx)),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &sug.suggestions,
			},
		},
	}

	return layout
}

func (sug *inline_popup_suggestions) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	content_size := sug.suggestions.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()
	gap := basic.Gap(ctx)

	b.Min.Y -= content_size.Y + gap
	b.Max.Y = b.Min.Y + content_size.Y
	b.Min.X = sug.x - (content_size.X / 2)
	b.Max.X = sug.x + (content_size.X / 2)
	sug.layout(ctx).LayoutWidgets(ctx, b, layouter)
}

func (sug *inline_popup_suggestions) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return image.Point{}
}

func (sug *inline_popup_suggestions) SetSuggestions(suggestions []string) {
}

func (sug *inline_popup_suggestions) SetPosition(x int) {
	sug.x = x
}
*/

type SearchBar struct {
	gui.DefaultWidget

	search_bar      widget.TextInput
	search_icon     *ebiten.Image
	commited, typed bool
	t               time.Time

	on_search_fn func(context *gui.Context, query string)

	//suggestions inline_popup_suggestions
}

func (search_bar *SearchBar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if search_bar.search_icon == nil {
		search_bar.search_icon = icons.Store.Open("search")
	}
	search_bar.search_bar.SetIcon(search_bar.search_icon)
	search_bar.search_bar.SetPlaceholder("Search")
	search_bar.search_bar.OnValueChangedWithoutText(func(context *gui.Context, committed bool) {
		if committed {
			search_bar.commited = committed
		} else {
			search_bar.typed = true
		}
	})

	if search_bar.typed && (time.Since(search_bar.t).Milliseconds() > 500 || search_bar.t.IsZero()) || search_bar.commited {
		if search_bar.on_search_fn != nil {
			search_bar.on_search_fn(ctx, search_bar.search_bar.Value())
		}
		search_bar.t = time.Now()
		search_bar.typed = false
		search_bar.commited = false
	}

	adder.AddWidget(&search_bar.search_bar)
	//adder.AddWidget(&search_bar.suggestions)
	return nil
}

func (search_bar *SearchBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(&search_bar.search_bar, b)
	//	t, _, _ := search_bar.search_bar.CaretPositionAtTextIndexInBytes(ctx, len(search_bar.search_bar.Value())-1)
	//
	// search_bar.suggestions.SetPosition(t.X)
	// layouter.LayoutWidget(&search_bar.suggestions, b)
}

func (search_bar *SearchBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return search_bar.search_bar.Measure(ctx, constraints)
}

func (search_bar *SearchBar) OnSearch(fn func(context *gui.Context, query string)) {
	search_bar.on_search_fn = fn
}

func (search_bar *SearchBar) SetQuery(query string) {
	search_bar.search_bar.SetValue(query)
}

func (search_bar *SearchBar) Query() string {
	return search_bar.search_bar.Value()
}
