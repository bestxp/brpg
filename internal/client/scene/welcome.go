package scene

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"

	"github.com/bestxp/brpg/internal/client/gui"
	e "github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var _ Scene = &WelcomeScene{}

type WelcomeScene struct {
	events EventStore

	loginButton *gui.Button
	frame       int
}

func (scene *WelcomeScene) Update() {
}

func NewWelcomeScene() *WelcomeScene {
	w := &WelcomeScene{
		events:      NewEventStore(),
		loginButton: gui.NewButton("Enter", 250, 70),
	}

	w.loginButton.TextOps.ColorScale.ScaleWithColor(colornames.Black)

	return w
}

func (scene *WelcomeScene) DrawAt(screen *e.Image) {
	w, h := e.WindowSize()
	img := e.NewImage(w, h)

	img.Fill(color.RGBA{R: 124, G: 102, B: 135, A: 255})

	scene.loginButton.DrawAt(img, 0, 0)

	screen.DrawImage(img, &e.DrawImageOptions{})
}

func (scene *WelcomeScene) OnClick(elemID string, event GUIEvent, fn EventFn) {
	scene.events.Add(elemID, event, fn)
}

func (scene *WelcomeScene) MouseMove(x, y int) {
	if scene.loginButton.TryHover(x, y) {
		if inpututil.IsMouseButtonJustPressed(e.MouseButtonLeft) {
			ev := scene.events.Get("login", Click)
			for _, fn := range ev {
				fn()
			}
		}
	}
}

func (scene *WelcomeScene) Frame(frame int) {
	scene.frame = frame
}
