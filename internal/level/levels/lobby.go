package levels

import (
	"github.com/bestxp/brpg/internal/level"
	"github.com/bestxp/brpg/internal/resources"
)

func GetLobbyLevel() *level.Level {
	frames, err := resources.LoadResources()
	if err != nil {
		panic(err)
	}

	l := level.NewLevel("Trial's Cave", frames)
	l.StartPos.X = 100
	l.StartPos.Y = 100

	floor := level.Tile{
		CanWalk: true,
		Texture: "floor_1",
	}
	wall := level.Tile{
		CanWalk: false,
		Texture: "floor_4",
	}

	ladder := level.Tile{
		CanWalk: true,
		Texture: "floor_ladder",
	}

	l.SetMap([][]level.Tile{
		{wall, wall, wall, wall, wall, wall, wall},
		{wall, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, ladder, wall},
		{wall, wall, wall, wall, wall, wall, wall},
	})

	return l
}
