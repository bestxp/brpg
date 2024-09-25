package keyboard

import (
	"fmt"

	"github.com/bestxp/brpg/internal/infra/network"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
)

type GameKeyboard struct {
	conn    *network.Network
	lastkey e.Key
	prevKey e.Key

	unitid    string
	lastEvent pkg.Event_Type
}

func NewGameKeyboard(unitid string, conn *network.Network) *GameKeyboard {
	return &GameKeyboard{conn: conn, unitid: unitid}
}

func (k *GameKeyboard) Handle() error {
	event := &pkg.Event{
		Type: pkg.Event_type_idle,
	}

	if e.IsKeyPressed(e.KeyV) && k.lastkey != e.KeyV {
		e.SetVsyncEnabled(!e.IsVsyncEnabled())
		k.lastkey = e.KeyV
	}

	if e.IsKeyPressed(e.KeyA) || e.IsKeyPressed(e.KeyLeft) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  k.unitid,
					Direction: pkg.Direction_left,
				},
			},
		}
		if k.lastkey != e.KeyA {
			k.lastkey = e.KeyA
		}
	}

	if e.IsKeyPressed(e.KeyD) || e.IsKeyPressed(e.KeyRight) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  k.unitid,
					Direction: pkg.Direction_right,
				},
			},
		}
		if k.lastkey != e.KeyD {
			k.lastkey = e.KeyD
		}
	}

	if e.IsKeyPressed(e.KeyW) || e.IsKeyPressed(e.KeyUp) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  k.unitid,
					Direction: pkg.Direction_up,
				},
			},
		}
		if k.lastkey != e.KeyW {
			k.lastkey = e.KeyW
		}
	}

	if e.IsKeyPressed(e.KeyS) || e.IsKeyPressed(e.KeyDown) {
		event = &pkg.Event{
			Type: pkg.Event_type_move,
			Data: &pkg.Event_Move{
				Move: &pkg.EventMove{
					PlayerId:  k.unitid,
					Direction: pkg.Direction_down,
				},
			},
		}
		if k.lastkey != e.KeyS {
			k.lastkey = e.KeyS
		}
	}

	if event.Type == pkg.Event_type_move {
		if k.prevKey != k.lastkey {
			if err := k.conn.Send(event); err != nil {
				return fmt.Errorf("move: %w", err)
			}
		}
	} else {
		if k.lastEvent != pkg.Event_type_idle {
			event = &pkg.Event{
				Type: pkg.Event_type_idle,
				Data: &pkg.Event_Idle{
					Idle: &pkg.EventIdle{PlayerId: k.unitid},
				},
			}

			if err := k.conn.Send(event); err != nil {
				return fmt.Errorf("idle: %w", err)
			}
			k.lastkey = -1
		}
	}

	k.prevKey = k.lastkey
	k.lastEvent = event.Type

	return nil
}
