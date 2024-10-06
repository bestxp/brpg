package game

import (
	"math/rand"
	"time"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/internal/level"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/pkg"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

// World represents game state
type World struct {
	IsClient bool
	Units    map[string]*pkg.Unit
	MyID     string
	Levels   map[levels.LevelName]*level.Level

	tick uint8
}

func NewWorld(client bool) *World {
	w := &World{
		IsClient: client,
		Units:    make(map[string]*pkg.Unit, 0),
	}

	return w
}

func (world *World) Me() *pkg.Unit {
	return world.Units[world.MyID]
}

func (world *World) ActiveClientWorld() *level.Level {
	lvl, err := levels.LevelName(world.Me().Pos.Level).Level()
	if err != nil {
		panic(err)
	}
	return lvl
}

func (world *World) AddPlayer() string {
	skins := []string{"big_demon", "big_zombie"}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	level, err := levels.Lobby.Level()
	if err != nil {
		panic(err)
	}

	unit := &pkg.Unit{
		Id:     id,
		Frame:  int32(rnd.Intn(4)),
		Skin:   skins[rnd.Intn(len(skins))],
		Action: actions.UnitIdle.String(),
		Speed:  1,
		Info: &pkg.CharInfo{
			MaxHealth:     100,
			CurrentHealth: 85,
		},
		Pos: &pkg.Pos{
			X:     level.StartPos.X,
			Y:     level.StartPos.Y,
			Level: levels.Lobby.BaseName().String(),
		},
	}
	world.Units[id] = unit

	return id
}

func (world *World) HandleEvent(event *pkg.Event) {
	log.Debug().Msg(event.GetType().String())

	switch event.GetType() {
	case pkg.Event_type_connect:
		data := event.GetConnect()
		world.Units[data.Unit.Id] = data.Unit

	case pkg.Event_type_teleport:
		data := event.GetTeleport()
		world.Units[data.PlayerId].Pos = data.Pos

	case pkg.Event_type_init:
		data := event.GetInit()
		if world.IsClient {
			world.MyID = data.PlayerId
			world.Units = data.Units
		}

	case pkg.Event_type_exit:
		data := event.GetExit()
		delete(world.Units, data.PlayerId)

	case pkg.Event_type_move:
		data := event.GetMove()
		unit := world.Units[data.PlayerId]
		unit.Action = actions.UnitMove.String()
		unit.Direction = data.Direction

	case pkg.Event_type_idle:
		data := event.GetIdle()
		if unit := world.Units[data.PlayerId]; unit != nil {
			unit.Action = actions.UnitIdle.String()
		}

	default:
		log.Error().Msgf("UNKNOWN EVENT: %s", event)
	}
}

func (world *World) tickUnits() {
	for _, unit := range world.Units {
		if world.IsClient && world.Me().Pos.GetLevel() != unit.Pos.GetLevel() {
			continue
		}
		// todo server only logic
		if world.tick == 60 {
			unit.Info.CurrentHealth++
			if unit.Info.CurrentHealth > unit.Info.MaxHealth {
				unit.Info.CurrentHealth = unit.Info.MaxHealth
			}
		}

		if unit.Action == actions.UnitMove.String() {
			side := unit.Side
			posTo := level.Pos{X: unit.Pos.X, Y: unit.Pos.Y}

			switch unit.Direction {
			case pkg.Direction_left:
				posTo.X -= unit.Speed
				side = pkg.Direction_left
			case pkg.Direction_right:
				posTo.X += unit.Speed
				side = pkg.Direction_right
			case pkg.Direction_up:
				posTo.Y -= unit.Speed
			case pkg.Direction_down:
				posTo.Y += unit.Speed
			default:
				log.Error().Msgf("UNKNOWN DIRECTION: %v", unit.Direction)
			}

			lvl, err := levels.LevelName(unit.Pos.Level).Level()
			if err != nil {
				log.Debug().Msg("UNKNOWN level")
				return
			}

			posTo = lvl.WalkCalc(level.Vector{
				From: level.Pos{X: unit.Pos.X, Y: unit.Pos.Y},
				To:   posTo,
			})
			unit.Side = side
			unit.Pos.X = posTo.X
			unit.Pos.Y = posTo.Y
		}
	}
}

func (world *World) Evolve() {
	ticker := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-ticker.C:
			world.tick++
			world.tickUnits()
			if world.tick > 60 {
				world.tick = 0
			}
		}
	}
}
