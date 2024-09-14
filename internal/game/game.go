package game

import (
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
}

func (world *World) AddPlayer() string {
	skins := []string{"big_demon", "big_zombie", "elf_f"}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	unit := &pkg.Unit{
		Id:     id,
		X:      100,
		Y:      100,
		Frame:  int32(rnd.Intn(4)),
		Skin:   skins[rnd.Intn(len(skins))],
		Action: actions.UnitIdle.String(),
		Speed:  1,
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
		if unit.X < 0 {
			unit.X = 0
		}
		if unit.Y < 0 {
			unit.Y = 0
		}

	case pkg.Event_type_idle:
		data := event.GetIdle()
		unit := world.Units[data.PlayerId]
		unit.Action = actions.UnitIdle.String()

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (world *World) Evolve() {
	ticker := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-ticker.C:
			for _, unit := range world.Units {
				if unit.Action == actions.UnitMove.String() {
					switch unit.Direction {
					case pkg.Direction_left:
						unit.X -= unit.Speed
						unit.Side = pkg.Direction_left
					case pkg.Direction_right:
						unit.X += unit.Speed
						unit.Side = pkg.Direction_right
					case pkg.Direction_up:
						unit.Y -= unit.Speed
					case pkg.Direction_down:
						unit.Y += unit.Speed
					default:
						log.Println("UNKNOWN DIRECTION: ", unit.Direction)
					}
				}
			}
		}
	}
}
