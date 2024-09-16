package levels

import "github.com/bestxp/brpg/internal/level"

type LevelName string

const (
	Lobby  LevelName = "lobby"
	Sewage LevelName = "sewage"
)

func (l LevelName) String() string {
	return string(l)
}

var levels = map[LevelName]*level.Level{}

func init() {
	levels[Lobby] = GetLobbyLevel()
	levels[Sewage] = GetSewageLevel()
}

func All() map[LevelName]*level.Level {
	return levels
}

func Level(l LevelName) *level.Level {
	return levels[l]
}
