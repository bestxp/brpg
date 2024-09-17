package levels

import (
	"github.com/bestxp/brpg/internal/level"
	"github.com/bestxp/brpg/internal/resources"
)

func GetSewageLevel() *level.Level {
	frames, err := resources.LoadResources()
	if err != nil {
		panic(err)
	}

	l := level.NewLevel(Sewage.BaseName(), "Sewage Cave", frames)
	l.StartPos.X = 32
	l.StartPos.Y = 32

	floor := level.Tile{
		CanWalk: true,
		Texture: "floor_3",
	}
	wall := level.Tile{
		CanWalk: false,
		Texture: "floor_2",
	}

	ld := level.Tile{
		CanWalk: true,
		Texture: "floor_ladder",
	}

	l.SetMap([][]level.Tile{
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, ld, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, ld, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, ld, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, floor, floor, floor, floor, floor, floor, floor, floor, floor, floor, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
	})

	return l
}
