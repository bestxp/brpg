package level

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestLevel_WalcCalc(t *testing.T) {
	wall := Tile{
		CanWalk: false,
		Texture: "wall",
	}
	floor := Tile{
		CanWalk: true,
		Texture: "wall",
	}

	base := NewLevel(LevelName("test"), "test", nil)
	base.SetMap([][]Tile{
		{wall, wall, wall, wall},
		{wall, floor, floor, wall},
		{wall, floor, floor, wall},
		{wall, wall, wall, wall},
	})

	v := Vector{
		From: Pos{140, 0},
		To:   Pos{140, 1},
	}
	posTo := base.WalkCalc(v)
	assert.Equal(t, posTo, v.To)
}
