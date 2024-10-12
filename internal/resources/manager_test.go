package resources

import (
	"testing"

	"github.com/bestxp/brpg"
	"github.com/go-playground/assert/v2"
)

func TestManager_Load(t *testing.T) {

	fs := brpg.FS()

	manager := FromFS(fs, false)
	textures, err := manager.Load("gui")

	assert.Equal(t, nil, err)
	assert.Equal(t, len(textures.Images), 2)
}
