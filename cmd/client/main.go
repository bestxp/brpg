package main

import (
	"image"
	"log"
	"os"

	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	title  string
	width  int
	height int
	scale  int
}

type Sprite struct {
	Frames []image.Image
	Frame  int
	X      float64
	Y      float64
	Side   pkg.Direction
	Config image.Config
}

type Camera struct {
	X       float64
	Y       float64
	Padding float64
}

var world *game.World

var config *Config
var frames map[string]resources.Frames

func init() {
	config = &Config{
		title:  "Just Dungeon",
		width:  1024,
		height: 768,
		scale:  2,
	}

	world = &game.World{
		Replica: true,
		Units:   map[string]*pkg.Unit{},
	}

	var err error
	frames, err = resources.LoadResources()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	levels.All()
	go world.Evolve()
	var err error

	host := getEnv("HOST", "localhost")
	c, _, _ := websocket.DefaultDialer.Dial("ws://"+host+":3000/ws", nil)

	game := &Game{Conn: c, world: world}

	if err != nil {
		log.Fatal(err)
	}

	go func(c *websocket.Conn) {
		defer c.Close()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			event := &pkg.Event{}
			err = proto.Unmarshal(message, event)
			if err != nil {
				log.Fatal(err)
			}

			world.HandleEvent(event)

			if event.Type == pkg.Event_type_connect {
				me := world.Units[world.MyID]
				game.Camera = &Camera{
					X:       me.Pos.X,
					Y:       me.Pos.Y,
					Padding: 30,
				}
			}
		}
	}(c)

	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(config.width, config.height)
	e.SetWindowTitle(config.title)
	e.SetWindowResizable(true)

	if err := e.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
