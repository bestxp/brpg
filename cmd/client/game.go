package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/pkg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	Conn *websocket.Conn

	Camera     *Camera
	frame      int
	levelImage *e.Image
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
	// Write your game's rendering.
	g.handleCamera(screen)

	g.frame++

	sprites := []Sprite{}
	for _, unit := range world.Units {
		sprites = append(sprites, Sprite{
			Frames: frames[unit.Skin+"_"+unit.Action].Frames,
			Frame:  int(unit.Frame),
			X:      unit.X,
			Y:      unit.Y,
			Side:   unit.Side,
			Config: frames[unit.Skin+"_"+unit.Action].Config,
		})
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
	out := fmt.Sprintf("\n X: %f", me.X)
	out += fmt.Sprintf("\n Y: %f", me.Y)
	out += fmt.Sprintf("\n Speed: %f", me.Speed)

	return out
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.WindowSize()
}

func (g *Game) handleCamera(screen *e.Image) {
	if g.Camera == nil {
		return
	}

	player := world.Units[world.MyID]
	frame := frames[player.Skin+"_"+player.Action]
	g.Camera.X = player.X - float64(config.width-frame.Config.Width)/2
	g.Camera.Y = player.Y - float64(config.height-frame.Config.Height)/2

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(-g.Camera.X, -g.Camera.Y)
	screen.DrawImage(g.levelImage, op)
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
		if lastKey != e.KeyA {
			lastKey = e.KeyA
		}
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
		if lastKey != e.KeyD {
			lastKey = e.KeyD
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
		if lastKey != e.KeyW {
			lastKey = e.KeyW
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
		if lastKey != e.KeyS {
			lastKey = e.KeyS
		}
	}

	unit := world.Units[world.MyID]

	if event.Type == pkg.Event_type_move {
		if prevKey != lastKey {
			message, err := proto.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, message)
		}
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
			lastKey = -1
		}
	}

	prevKey = lastKey
}
