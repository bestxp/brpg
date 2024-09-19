package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/internal/client/camera"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	Conn *websocket.Conn

	gui *e.Image

	Camera           *camera.Camera
	frame            int
	world            *game.World
	lastKey, prevKey e.Key

	audioContext *audio.Context

	showDebug bool

	polygons []*Poly
}

type Poly struct {
	rect image.Rectangle
	icon *e.Image

	hovered bool
}

func (p *Poly) Icon() *e.Image {
	if p.icon == nil {
		p.icon = e.NewImage(32, 32)
		p.icon.Fill(color.White)
	}

	return p.icon
}

func (p *Poly) Hover() {
	if !p.hovered {
		p.icon.Fill(color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		})
		p.hovered = true
	}

}

func (p *Poly) Over() {
	if p.hovered {
		p.icon.Fill(color.White)
		p.hovered = false
	}
}

func (p *Poly) SetRect(r image.Rectangle) {
	p.rect = r
}

func NewGame(c *websocket.Conn, w *game.World) *Game {
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

		s, err := mp3.Decode(g.audioContext, track.Stream)
		if err != nil {
			panic(err)
		}
		loop := audio.NewInfiniteLoop(s, s.Length())

		p, err := audio.NewPlayer(g.audioContext, loop)
		if err != nil {
			panic(err)
		}
		p.SetVolume(0.5)
		p.Play()

	} else {
		println("no track")
	}

	g.polygons = make([]*Poly, 0, 100)
	for i := 0; i < 10; i++ {
		g.polygons = append(g.polygons, &Poly{})
	}

	return g
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	g.handleKeyboard(g.Conn)
	g.handleMouse(g.Conn)
	w, h := e.WindowSize()
	// @todo calc by dx dy of sprite
	g.Camera.FollowTarget(world.Me().Pos.X-8, world.Me().Pos.Y-8, float64(w), float64(h-100), 0, 50)
	wmap, _ := world.ActiveClientWorld().EImage()
	wmapW, wmapH := wmap.Size()
	ww, wh := e.WindowSize()
	g.Camera.Constrain(
		float64(wmapW),
		float64(wmapH),
		float64(ww),
		float64(wh-100),
	)
	return nil
}

func (g *Game) DrawGUI(screen *e.Image) {
	g.frame++

	sprites := []Sprite{}
	for _, unit := range g.world.Units {
		if unit.Pos.GetLevel() == g.world.Me().Pos.GetLevel() {
			sprites = append(sprites, Sprite{
				Frames: frames[unit.Skin+"_"+unit.Action].Frames,
				Frame:  int(unit.Frame),
				X:      unit.Pos.X,
				Y:      unit.Pos.Y,
				Side:   unit.Side,
				Config: frames[unit.Skin+"_"+unit.Action].Config,
			})
		}
	}

	sort.Slice(sprites, func(i, j int) bool {
		depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
		depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
		return depth1 < depth2
	})

	for _, sprite := range sprites {
		op := &e.DrawImageOptions{}

		if sprite.Side == pkg.Direction_left {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(sprite.Config.Width), 0)
		}

		op.GeoM.Translate(sprite.X+g.Camera.X, sprite.Y+g.Camera.Y)

		img := e.NewImageFromImage(sprite.Frames[(g.frame/7+sprite.Frame)%4])
		screen.DrawImage(img, op)
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *e.Image) {
	for i, p := range g.polygons {
		icon := p.Icon()
		w, h := icon.Size()
		x := (w+(50-w)/2)*i + (50-w)/2
		y := (50 - h) / 2
		p.SetRect(image.Rect(x, y, x+w, y+h))
		icOp := &e.DrawImageOptions{}
		icOp.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(icon, icOp)
	}

	w, h := e.WindowSize()
	g.gui = e.NewImage(w, h-100)
	g.gui.Fill(color.RGBA{
		R: 90,
		G: 90,
		B: 90,
		A: 255,
	})

	// render world tiles
	worldOp := &e.DrawImageOptions{}
	worldOp.GeoM.Translate(g.Camera.X, g.Camera.Y)
	worldMap, err := g.world.ActiveClientWorld().EImage()
	if err != nil {
		panic(err)
	}
	g.gui.DrawImage(worldMap, worldOp)
	g.DrawGUI(g.gui)

	guiOp := &e.DrawImageOptions{}
	guiOp.GeoM.Translate(0, 50)
	screen.DrawImage(g.gui, guiOp)

	if g.showDebug {
		ebitenutil.DebugPrint(screen,
			fmt.Sprintf("TPS: %0.2f \n VS: %v \n", e.CurrentTPS(), e.IsVsyncEnabled())+
				g.UnitInfo()+
				PrintMemUsage())
	}
}

func (g *Game) UnitInfo() string {
	me := world.Me()
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

func (g *Game) handleKeyboard(c *websocket.Conn) {
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
					PlayerId:  world.MyID,
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
					PlayerId: world.MyID,
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
					PlayerId: world.MyID,
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
					PlayerId:  world.MyID,
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
					PlayerId:  world.MyID,
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
					PlayerId:  world.MyID,
					Direction: pkg.Direction_down,
				},
			},
		}
		if g.lastKey != e.KeyS {
			g.lastKey = e.KeyS
		}
	}

	unit := world.Units[world.MyID]

	if event.Type == pkg.Event_type_move {
		if g.prevKey != g.lastKey {
			message, err := proto.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, message)
		}
	} else if event.Type == pkg.Event_type_teleport {
		message, err := proto.Marshal(event)
		if err != nil {
			log.Println(err)
			return
		}
		c.WriteMessage(websocket.BinaryMessage, message)
	} else {
		if unit.Action != actions.UnitIdle.String() {
			event = &pkg.Event{
				Type: pkg.Event_type_idle,
				Data: &pkg.Event_Idle{
					Idle: &pkg.EventIdle{PlayerId: world.MyID},
				},
			}
			message, err := proto.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, message)
			g.lastKey = -1
		}
	}

	g.prevKey = g.lastKey
}

func (g *Game) handleMouse(conn *websocket.Conn) {
	x, y := e.CursorPosition()
	for _, p := range g.polygons {
		if p.rect.Overlaps(image.Rect(x, y, x+1, y+1)) {
			p.Hover()
		} else {
			p.Over()
		}
	}
}
