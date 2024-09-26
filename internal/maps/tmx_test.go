package maps

import (
	"testing"

	"github.com/bestxp/brpg"
	"github.com/stretchr/testify/assert"
)

const baseMap = "resources/Tiled/Tilemaps/Beginning Fields.tmx"

func TestMapLoader(t *testing.T) {
	loader := FromFs(brpg.FS())

	m, err := loader.Load(baseMap)
	assert.NotNilf(t, m, "map not parced")
	assert.NoErrorf(t, err, "file exists")
}
