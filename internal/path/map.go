package path

import "github.com/bestxp/brpg/internal/level"

type Path struct {
	level *level.Level
}

func New(l *level.Level) *Path {
	return &Path{level: l}
}
