// icon transparency is “CF“ if icon is too bright.
package icons

import (
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"sync"
	"weak"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var store fs.FS = func() fs.FS {
	_, err := os.Stat("./icons/icons")
	if err == nil {
		return os.DirFS("./icons/icons")
	}

	_, err = os.Stat("./icons")
	if err != nil {
		log.Fatal(err.Error())

	}
	return os.DirFS("./icons")
}()

type cache_store struct {
	mutex  sync.Mutex
	images map[string]*weak.Pointer[ebiten.Image]
}

func (cs *cache_store) Open(icon_name string) *ebiten.Image {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	icon_name = icon_name + ".png"
	ebiten_image := cs.images[icon_name]
	if ebiten_image == nil || ebiten_image.Value() == nil {
		file, err := store.Open(icon_name)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err.Error())
		}

		eg := weak.Make(ebiten.NewImageFromImage(img))
		ebiten_image = &eg
		cs.images[icon_name] = ebiten_image
	}

	return ebiten_image.Value()
}

func new() *cache_store {
	cs := cache_store{
		images: make(map[string]*weak.Pointer[ebiten.Image], 10),
	}

	return &cs
}

var (
	Store *cache_store = new()
)

type Icon struct {
	gui.DefaultWidget

	image_widget widget.Image
	icon_name    string

	size_specified bool
	size           image.Point

	on_click func()
}

func (icon *Icon) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&icon.image_widget)
	return nil
}

func (icon *Icon) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&icon.image_widget, widgetBounds.Bounds())
}

func (icon *Icon) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if !icon.size_specified {
		u := widget.UnitSize(ctx)
		icon.SetSize(u)
	}
	return icon.size
}

func (icon *Icon) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && icon.on_click != nil {
		icon.on_click()
	}
	return gui.HandleInputResult{}
}

func (icon *Icon) OnClick(fn func()) {
	icon.on_click = fn
}

func (icon *Icon) SetSize(size ...int) {
	var pt image.Point
	if len(size) == 0 {
		return
	} else if len(size) == 1 {
		point := image.Pt(size[0], size[0])
		pt = point
	} else if len(size) == 2 {
		point := image.Pt(size[0], size[1])
		pt = point
	} else if len(size) > 2 {
		log.Fatal("Extra arguments to NewIcon")
	}
	icon.size_specified = true
	icon.size = pt
	//gui.RequestRebuild(icon)
}

func (icon *Icon) Size() image.Point {
	return icon.size
}

func (icon *Icon) SetIcon(icon_name string) {
	if icon.icon_name == icon_name {
		return
	}
	icon.icon_name = icon_name
	icon.image_widget.SetImage(Store.Open(icon_name))
	//gui.RequestRedraw(&icon.image_widget)
}

func (icon *Icon) IconName() string {
	return icon.icon_name
}

// TODO: remove this function
func NewIcon(icon_name string, size ...int) *Icon {
	icon := Icon{}
	if len(size) > 0 {
		icon.SetSize(size...)
	}
	icon.SetIcon(icon_name)
	return &icon
}
