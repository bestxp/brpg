package main

import (
	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/level/levels"
	"github.com/gin-gonic/gin"
)

func main() {
	var world = game.NewWorld(false)

	levels.All()
	go world.Evolve()

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
