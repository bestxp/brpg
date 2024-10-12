package scene

import (
	"image/color"
	"sort"

	"github.com/bestxp/brpg"
	"github.com/bestxp/brpg/internal/client/camera"
	"github.com/bestxp/brpg/internal/client/gui"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/resources"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
)

var _ Scene = &GameScene{}

type GameScene struct {
	events EventStore

	guiElements []*gui.Icon
	frame       int
	players     []*gui.Player

	camera *camera.Camera
	world  *game.World

	keyboards []KeyboardInterface
}

type KeyboardInterface interface {
	Handle() error
}

func NewGameScene(world *game.World, k ...KeyboardInterface) *GameScene {
	w := &GameScene{
		events: NewEventStore(),
		world:  world,
	}

	gg := resources.FromFS(brpg.FS(), !world.IsClient)

	txt, err := gg.Load("gui")
	if err != nil {
		log.Error().Err(err)
	}
	if txt != nil {
		icons := txt.Images["tiled-icons"]
		if icons != nil && icons.IsTiled() {
			w.guiElements = make([]*gui.Icon, 0, 100)
			for i := 0; i < 10; i++ {
				img, err := icons.Tile(5, i+1)
				if err != nil {
					log.Error().Err(err)
					continue
				}
				if img == nil {
					log.Error().Msg("nil img")
				}
				w.guiElements = append(w.guiElements, gui.NewIcon(img))
			}
		} else {
			log.Error().Msg("No gui elements")
		}
	}

	w.players = make([]*gui.Player, 0)
	w.camera = camera.NewCamera(0, 0)
	w.keyboards = k

	return w
}

func (scene *GameScene) Update() {
	lenKey := len(scene.keyboards)
	for i := 0; i < lenKey; i++ {
		if err := scene.keyboards[i].Handle(); err != nil {
			log.Debug().Msgf("keyboard", err.Error())
		}
	}

	// @todo calc by dx dy of sprite
	w, h := e.WindowSize()

	if scene.world.Me() == nil {
		return
	}

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
