package main

import (
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	engine "github.com/bestxp/brpg/pkg"
	"github.com/gin-gonic/gin"
)

var world *game.World

func init() {
	world = &game.World{
		Replica: false,
		Units:   map[string]*engine.Unit{},
	}
}

func main() {
	go world.Evolve(levels.All())

	hub := newHub()
	go hub.run()

	r := gin.New()
	r.GET("/ws", ginWsServe(hub, world))
	r.Run(":3000")
}

func ginWsServe(hub *Hub, world *game.World) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		serveWs(hub, world, c.Writer, c.Request)
	})
}
