package main

import (
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/pkg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	Conn *websocket.Conn

	gui *e.Image

	Camera           *Camera
	frame            int
	world            *game.World
	lastKey, prevKey e.Key
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	g.handleKeyboard(g.Conn)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *e.Image) {
	w, h := e.WindowSize()
	g.gui = e.NewImage(w, h-100)
	g.gui.Fill(color.RGBA{
		R: 90,
		G: 90,
		B: 90,
		A: 255,
	})
	guiOp := &e.DrawImageOptions{}
	guiOp.GeoM.Translate(0, 50)

	// Write your game's rendering.
	g.handleCamera(g.gui, w, h-100, 0, 50)

	screen.DrawImage(g.gui, guiOp)

	for i := 0; i < 10; i++ {
		icon := e.NewImage(32, 32)
		icon.Fill(color.White)
		icOp := &e.DrawImageOptions{}
		icOp.GeoM.Translate((32+(50-32)/2)*float64(i)+(50-32)/2, (50-32)/2)

		screen.DrawImage(icon, icOp)
	}

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

		op.GeoM.Translate(sprite.X-g.Camera.X, sprite.Y-g.Camera.Y)

		img := e.NewImageFromImage(sprite.Frames[(g.frame/7+sprite.Frame)%4])
		screen.DrawImage(img, op)
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("TPS: %0.2f \n VS: %v \n", e.CurrentTPS(), e.IsVsyncEnabled())+
			g.UnitInfo()+
			PrintMemUsage())
}

func (f *Game) UnitInfo() string {
	me := world.Me()
	out := fmt.Sprintf("\n X: %f", me.Pos.X)
	out += fmt.Sprintf("\n Y: %f", me.Pos.Y)
	out += fmt.Sprintf("\n Lvl: %s", me.Pos.Level)
	out += fmt.Sprintf("\n Speed: %f", me.Speed)

	return out
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.WindowSize()
}

func (g *Game) handleCamera(screen *e.Image, w, h int, shiftX, shiftY float64) {
	if g.Camera == nil {
		return
	}

	player := g.world.Me()
	frame := frames[player.Skin+"_"+player.Action]

	g.Camera.X = player.Pos.X - float64(w-frame.Config.Width)/2
	g.Camera.Y = player.Pos.Y - float64(h-frame.Config.Height)/2

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(-g.Camera.X-shiftX, -g.Camera.Y-shiftY)
	img, err := g.world.ActiveClientWorld().EImage()
	if err != nil {
		panic(err)
	}
	screen.DrawImage(img, op)
}

func (g *Game) handleKeyboard(c *websocket.Conn) {
	event := &pkg.Event{}

	if e.IsKeyPressed(e.KeyV) {
		e.SetVsyncEnabled(!e.IsVsyncEnabled())
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
