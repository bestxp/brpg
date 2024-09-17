package levels

import (
	"fmt"

	"github.com/bestxp/brpg/internal/level"
)

type LevelName level.LevelName

const (
	Lobby  LevelName = "lobby"
	Sewage LevelName = "sewage"
)

func (l LevelName) BaseName() level.LevelName {
	return level.LevelName(l)
}

func (l LevelName) Level() (*level.Level, error) {
	if lvl, ok := levels[l]; ok {
		return lvl, nil
	}
	return nil, fmt.Errorf("%s not exists", l)
}

var levels = map[LevelName]*level.Level{}

func init() {
	levels[Lobby] = GetLobbyLevel()
	levels[Sewage] = GetSewageLevel()
}

func All() map[LevelName]*level.Level {
	return levels
}
