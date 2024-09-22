package main

import (
	"github.com/bestxp/brpg/internal/infra/network"
	"log"
	"os"

	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	title  string
	width  int
	height int
}

var world *game.World

var config *Config
var frames map[string]resources.Frames

func init() {
	config = &Config{
		title:  "Just Dungeon",
		width:  1024,
		height: 768,
	}

	world = &game.World{
		IsClient: true,
		Units:    map[string]*pkg.Unit{},
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
	n := network.NewNetwork(c)
	game := NewGame(n, world)

	if err != nil {
		log.Fatal(err)
	}

	go func(c *websocket.Conn) {
		defer c.Close()

		for {
			_, event, err := n.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			world.HandleEvent(event)

			if event.Type == pkg.Event_type_connect {
				game.Camera.InitCoords(world.Me().Pos.X, world.Me().Pos.Y)
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
