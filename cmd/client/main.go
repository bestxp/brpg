package main

import (
	"os"

	"github.com/bestxp/brpg/internal/client/config"
	game2 "github.com/bestxp/brpg/internal/client/game"
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/infra/network"
	"github.com/bestxp/brpg/internal/level/levels"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

type Config struct {
	title string
}

func main() {
	conf := &Config{
		title: "Monster Dungeon",
	}

	cfg := config.GetConfig()

	w, h := e.Monitor().Size()
	log.Debug().Msgf("mon size %d x %d, requst size %s", w, h, cfg.Resolution)

	if w > cfg.Resolution.Width() {
		w = cfg.Resolution.Width()
	}
	if h > cfg.Resolution.Height() {
		h = cfg.Resolution.Height()
	} else {
		h -= 100
	}

	log.Debug().Msgf("win size %d x %d", w, h)

	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(w, h)
	e.SetWindowTitle(conf.title)
	e.SetWindowResizingMode(e.WindowResizingModeEnabled)
	e.SetVsyncEnabled(cfg.VsyncEnabled)

	world := game.NewWorld(true)
	var err error

	levels.All()
	go world.Evolve()

	n := network.FromHost(cfg.Host)
	if n == nil {
		log.Fatal().Msg("Can't connect to remote server")
		return
	}
	gg := game2.NewGame(n, world)

	if err != nil {
		log.Fatal().Err(err)
	}

	if err := e.RunGame(gg); err != nil {
		log.Fatal().Err(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
