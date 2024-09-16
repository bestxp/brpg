package game

import (
	"github.com/bestxp/brpg/internal/level"
	"github.com/bestxp/brpg/internal/level/levels"
	"log"
	"math/rand"
	"time"

	"github.com/bestxp/brpg/internal/actions"
	"github.com/bestxp/brpg/pkg"
	uuid "github.com/satori/go.uuid"
)

// World represents game state
type World struct {
	Replica bool
	Units   map[string]*pkg.Unit
	MyID    string
	Levels  map[levels.LevelName]*level.Level
}

func (world *World) Me() *pkg.Unit {
	return world.Units[world.MyID]
}

func (world *World) ActiveClientWorld() *level.Level {
	return world.Levels[levels.LevelName(world.Me().Pos.Level)]
}

func (world *World) AddPlayer() string {
	skins := []string{"big_demon", "big_zombie"}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.Levels[levels.Lobby]
	unit := &pkg.Unit{
		Id:     id,
		Frame:  int32(rnd.Intn(4)),
		Skin:   skins[rnd.Intn(len(skins))],
		Action: actions.UnitIdle.String(),
		Speed:  1,

		Pos: &pkg.Pos{
			X:     level.StartPos.X,
			Y:     level.StartPos.Y,
			Level: levels.Lobby.String(),
		},
	}
	world.Units[id] = unit

	return id
}

func (world *World) HandleEvent(event *pkg.Event) {
	log.Println(event.GetType())

	switch event.GetType() {
	case pkg.Event_type_connect:
		data := event.GetConnect()
		world.Units[data.Unit.Id] = data.Unit

	case pkg.Event_type_teleport:
		data := event.GetTeleport()
		world.Units[data.PlayerId].Pos = data.Pos

	case pkg.Event_type_init:
		data := event.GetInit()
		if world.Replica {
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
		unit := world.Units[data.PlayerId]
		unit.Action = actions.UnitIdle.String()

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (world *World) Evolve(lvl map[levels.LevelName]*level.Level) {
	world.Levels = lvl

	ticker := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-ticker.C:
			for _, unit := range world.Units {
				if world.Replica && world.Me().Pos.GetLevel() != unit.Pos.GetLevel() {
					continue
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
						log.Println("UNKNOWN DIRECTION: ", unit.Direction)
					}

					posTo = levels.Level(levels.LevelName(unit.Pos.Level)).WalkCalc(level.Vector{
						From: level.Pos{X: unit.Pos.X, Y: unit.Pos.Y},
						To:   posTo,
					})
					unit.Side = side
					unit.Pos.X = posTo.X
					unit.Pos.Y = posTo.Y
				}
			}
		}
	}
}
