package scene

import (
	"image/color"
	"sort"

	"github.com/bestxp/brpg/internal/client/camera"
	"github.com/bestxp/brpg/internal/client/gui"
	"github.com/bestxp/brpg/internal/game"

	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ Scene = &GameScene{}

type GameScene struct {
	events EventStore

	guiElements []*gui.Icon
	frame       int
	players     []*gui.Player

	camera *camera.Camera
	world  *game.World
}

func (scene *GameScene) Update() {
	// @todo calc by dx dy of sprite
	w, h := e.WindowSize()

	scene.camera.FollowTarget(scene.world.Me().Pos.X, scene.world.Me().Pos.Y, float64(w), float64(h-100), 0, 50)
	wmap, _ := scene.world.ActiveClientWorld().EImage()
	wmapW, wmapH := wmap.Size()
	ww, wh := e.WindowSize()
	scene.camera.Constrain(
		float64(wmapW),
		float64(wmapH),
		float64(ww),
		float64(wh-100),
	)

	scene.players = scene.players[0:0]
	for _, unit := range scene.world.Units {
		if unit.Pos.Level != scene.world.Me().Pos.Level {
			continue
		}
		scene.players = append(scene.players, gui.NewPlayerFromServer(unit))
	}
}

func (scene *GameScene) Frame(frame int) {
	scene.frame = frame
}

func NewGameScene(world *game.World) *GameScene {
	w := &GameScene{
		events: NewEventStore(),
		world:  world,
	}

	w.guiElements = make([]*gui.Icon, 0, 100)
	for i := 0; i < 10; i++ {
		w.guiElements = append(w.guiElements, gui.NewIcon())
	}

	w.players = make([]*gui.Player, 0)
	w.camera = camera.NewCamera(0, 0)

	return w
}

func (scene *GameScene) DrawAt(screen *e.Image) {
	w, h := e.WindowSize()
	img := e.NewImage(w, h)

	for i, p := range scene.guiElements {
		w, h := p.Size()

		x := (w+(50-w)/2)*i + (50-w)/2
		y := (50 - h) / 2
		p.SetFrame(scene.frame)
		p.DrawAt(screen, float64(x), float64(y))
	}

	mainFrame := e.NewImage(w, h-100)
	mainFrame.Fill(color.RGBA{
		R: 90,
		G: 90,
		B: 90,
		A: 255,
	})

	// render world tiles
	worldOp := &e.DrawImageOptions{}
	worldOp.GeoM.Translate(scene.camera.X, scene.camera.Y)
	worldMap, err := scene.world.ActiveClientWorld().EImage()
	if err != nil {
		panic(err)
	}
	mainFrame.DrawImage(worldMap, worldOp)
	scene.Players(mainFrame)

	guiOp := &e.DrawImageOptions{}
	guiOp.GeoM.Translate(0, 50)
	screen.DrawImage(mainFrame, guiOp)

	screen.DrawImage(img, &e.DrawImageOptions{})
}

func (scene *GameScene) OnClick(elemID string, event GUIEvent, fn EventFn) {
	scene.events.Add(elemID, event, fn)
}

func (scene *GameScene) MouseMove(x, y int) {
	for _, p := range scene.guiElements {
		if p.TryHover(x, y) {
			if inpututil.IsMouseButtonJustPressed(e.MouseButtonRight) {
				p.HandleClick(e.MouseButtonRight)
			}
			if inpututil.IsMouseButtonJustPressed(e.MouseButtonLeft) {
				p.HandleClick(e.MouseButtonLeft)
			}
			if inpututil.IsMouseButtonJustPressed(e.MouseButtonMiddle) {
				p.HandleClick(e.MouseButtonMiddle)
			}
		}
	}
}

func (scene *GameScene) Players(screen *e.Image) {
	for _, pl := range scene.players {
		pl.SetFrame(scene.frame)
		pl.DrawAt(screen, pl.X+scene.camera.X, pl.Y+scene.camera.Y)
	}

	sort.Slice(scene.players, func(i, j int) bool {
		_, h1 := scene.players[i].Size()
		depth1 := scene.players[i].Y + float64(h1)
		_, h2 := scene.players[j].Size()
		depth2 := scene.players[j].Y + float64(h2)
		return depth1 < depth2
	})

}
