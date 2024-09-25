package game

import (
	"fmt"
	"log"

	"github.com/bestxp/brpg/internal/client/camera"
	"github.com/bestxp/brpg/internal/client/gui"
	"github.com/bestxp/brpg/internal/client/keyboard"
	"github.com/bestxp/brpg/internal/client/scene"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/infra/network"
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game implements ebiten.Game interface.
type Game struct {
	Conn *network.Network

	gui *e.Image

	Camera *camera.Camera
	frame  int
	world  *game.World

	audioContext *audio.Context

	showDebug bool
	players   []*gui.Player

	scene   scene.Scene
	welcome scene.Scene
}

func (g *Game) Run(afterConnectFn func()) {
	go func(c *network.Network) {
		defer c.Close()

		for {
			_, event, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			if event.Type == pkg.Event_type_connect {
				g.Camera.InitCoords(g.world.Me().Pos.X, g.world.Me().Pos.Y)
				afterConnectFn()
			}
			g.world.HandleEvent(event)
		}
	}(g.Conn)
}

func NewGame(c *network.Network, w *game.World) *Game {
	g := &Game{
		Conn:   c,
		world:  w,
		Camera: camera.NewCamera(0, 0),
	}
	au, err := resources.LoadAudios()
	if err != nil {
		panic(err)
	}
	g.audioContext = audio.NewContext(44100)
	track := au["/resources/audio/bg.mp3"]
	if track != nil {
		s, err := mp3.DecodeWithSampleRate(44100, track.Stream)
		if err != nil {
			panic(err)
		}
		loop := audio.NewInfiniteLoop(s, s.Length())
		p, err := g.audioContext.NewPlayer(loop)
		if err != nil {
			panic(err)
		}
		p.SetVolume(0.21)
		p.Play()

	} else {
		println("no track")
	}

	g.players = make([]*gui.Player, 0)

	s := scene.NewWelcomeScene()
	s.OnClick("login", scene.Click, func() error {
		g.Run(func() {
			g.scene = scene.NewGameScene(g.world, keyboard.NewGameKeyboard(g.world.MyID, g.Conn))
		})
		return nil
	})

	g.scene = s
	g.welcome = s

	return g
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	g.handleMouse()
	g.handleKeyboard()
	g.scene.Update()

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *e.Image) {
	g.frame++
	defer func() {
		if g.frame > 60 {
			g.frame = 0
		}
	}()

	g.scene.Frame(g.frame)
	g.scene.DrawAt(screen)

	if g.showDebug {
		ebitenutil.DebugPrint(screen,
			fmt.Sprintf("TPS: %0.2f \n VS: %v \n", e.ActualTPS(), e.IsVsyncEnabled())+
				g.UnitInfo()+
				PrintMemUsage())
	}
}

func (g *Game) UnitInfo() string {
	me := g.world.Me()
	out := fmt.Sprintf("\n X: %f", me.Pos.X)
	out += fmt.Sprintf("\n Y: %f", me.Pos.Y)
	out += fmt.Sprintf("\n Lvl: %s", me.Pos.Level)
	out += fmt.Sprintf("\n Speed: %f", me.Speed)

	return out
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return e.WindowSize()
}

func (g *Game) handleMouse() {
	x, y := e.CursorPosition()
	g.scene.MouseMove(x, y)

	for _, p := range g.players {
		// 50px shift of interfaces of game drawing
		if p.TryHover(x, y-50) {
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

func (g *Game) handleKeyboard() {
	if inpututil.IsKeyJustReleased(e.KeyF4) {
		g.showDebug = !g.showDebug
	}
}
