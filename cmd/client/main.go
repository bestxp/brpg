package main

import (
	"log"
	"os"

	game2 "github.com/bestxp/brpg/internal/client/game"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/infra/network"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	title  string
	width  int
	height int
}

func main() {
	config := &Config{
		title:  "Just Dungeon",
		width:  1024,
		height: 768,
	}
	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(config.width, config.height)
	e.SetWindowTitle(config.title)
	e.SetWindowResizingMode(e.WindowResizingModeEnabled)

	world := game.NewWorld(true)
	var err error

	levels.All()
	go world.Evolve()

	host := getEnv("HOST", "localhost")

	n := network.FromHost(host)
	gg := game2.NewGame(n, world)

	if err != nil {
		log.Fatal(err)
	}

	go func(c *network.Network) {
		defer c.Close()

		for {
			_, event, err := n.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			world.HandleEvent(event)

			if event.Type == pkg.Event_type_connect {
				gg.Camera.InitCoords(world.Me().Pos.X, world.Me().Pos.Y)
			}
		}
	}(n)

	if err := e.RunGame(gg); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
