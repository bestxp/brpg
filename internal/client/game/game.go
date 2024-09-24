package game

import (
	"fmt"
	"log"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/internal/client/camera"
	"github.com/bestxp/brpg/internal/client/gui"
	"github.com/bestxp/brpg/internal/client/scene"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/infra/network"
	"github.com/bestxp/brpg/internal/level/levels"
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

	Camera           *camera.Camera
	frame            int
	world            *game.World
	lastKey, prevKey e.Key

	audioContext *audio.Context

	showDebug bool
	players   []*gui.Player

	scene scene.Scene
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
		p.SetVolume(0.33)
		p.Play()

	} else {
		println("no track")
	}

	g.players = make([]*gui.Player, 0)

	s := scene.NewWelcomeScene()
	s.OnClick("login", scene.Click, func() error {
		g.scene = scene.NewGameScene(g.world)
		return nil
	})

	g.scene = s

	return g
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	g.handleKeyboard()
	g.handleMouse()
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

func (g *Game) handleKeyboard() {
	event := &pkg.Event{}

	if e.IsKeyPressed(e.KeyV) && g.lastKey != e.KeyV {
		e.SetVsyncEnabled(!e.IsVsyncEnabled())
		g.lastKey = e.KeyV
	}

	if e.IsKeyPressed(e.KeyF4) && g.lastKey != e.KeyF4 {
		g.showDebug = !g.showDebug
		g.lastKey = e.KeyF4
	}

	if e.IsKeyPressed(e.KeyA) || e.IsKeyPressed(e.KeyLeft) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  g.world.MyID,
					Direction: pkg.Direction_left,
				},
			},
		}
		if g.lastKey != e.KeyA {
			g.lastKey = e.KeyA
		}
	}

	if e.IsKeyPressed(e.KeyQ) && g.lastKey != e.KeyQ {
		startPos := levels.GetSewageLevel().StartPos
		event = &pkg.Event{
			Type: pkg.Event_type_teleport,
			Data: &pkg.Event_Teleport{
				Teleport: &pkg.EventTeleport{
					PlayerId: g.world.MyID,
					Pos: &pkg.Pos{
						X:     startPos.X,
						Y:     startPos.Y,
						Level: levels.Sewage.BaseName().String(),
					},
				},
			},
		}
		g.lastKey = e.KeyQ
	}
	if e.IsKeyPressed(e.KeyE) && g.lastKey != e.KeyE {
		startPos := levels.GetLobbyLevel().StartPos
		event = &pkg.Event{
			Type: pkg.Event_type_teleport,
			Data: &pkg.Event_Teleport{
				Teleport: &pkg.EventTeleport{
					PlayerId: g.world.MyID,
					Pos: &pkg.Pos{
						X:     startPos.X,
						Y:     startPos.Y,
						Level: levels.Lobby.BaseName().String(),
					},
				},
			},
		}
		g.lastKey = e.KeyE
	}

	if e.IsKeyPressed(e.KeyD) || e.IsKeyPressed(e.KeyRight) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  g.world.MyID,
					Direction: pkg.Direction_right,
				},
			},
		}
		if g.lastKey != e.KeyD {
			g.lastKey = e.KeyD
		}
	}

	if e.IsKeyPressed(e.KeyW) || e.IsKeyPressed(e.KeyUp) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  g.world.MyID,
					Direction: pkg.Direction_up,
				},
			},
		}
		if g.lastKey != e.KeyW {
			g.lastKey = e.KeyW
		}
	}

	if e.IsKeyPressed(e.KeyS) || e.IsKeyPressed(e.KeyDown) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  g.world.MyID,
					Direction: pkg.Direction_down,
				},
			},
		}
		if g.lastKey != e.KeyS {
			g.lastKey = e.KeyS
		}
	}

	unit := g.world.Units[g.world.MyID]

	if event.Type == pkg.Event_type_move {
		if g.prevKey != g.lastKey {
			err := g.Conn.Send(event)
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else if event.Type == pkg.Event_type_teleport {
		err := g.Conn.Send(event)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		if unit.Action != actions.UnitIdle.String() {
			event = &pkg.Event{
				Type: pkg.Event_type_idle,
				Data: &pkg.Event_Idle{
					Idle: &pkg.EventIdle{PlayerId: g.world.MyID},
				},
			}
			err := g.Conn.Send(event)
			if err != nil {
				log.Println(err)
				return
			}
			g.lastKey = -1
		}
	}

	g.prevKey = g.lastKey
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
