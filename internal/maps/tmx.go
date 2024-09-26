package maps

import (
	_ "image/gif"
	_ "image/png"

	"encoding/xml"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MapLoader struct {
	fs fs.FS
}

func FromFs(fs fs.FS) *MapLoader {
	return &MapLoader{
		fs: fs,
	}
}

func (f *MapLoader) Load(file string) (*Map, error) {
	dirname := filepath.Dir(file)
	fi, err := f.fs.Open(file)
	if err != nil {
		return nil, fmt.Errorf("read map: %w", err)
	}
	defer fi.Close()

	reader := xml.NewDecoder(fi)

	var m Map
	if err = reader.Decode(&m); err != nil {
		return nil, fmt.Errorf("decode map: %w", err)
	}

	for idx, tset := range m.Tileset {
		if !strings.HasPrefix(tset.Source, "/") {
			tset.Source = fmt.Sprintf("%s/%s", dirname, tset.Source)
		}
		tset.Source = filepath.Clean(tset.Source)
		tileFi, err := f.fs.Open(tset.Source)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", tset.Source, err)
		}

		reader := xml.NewDecoder(tileFi)
		tset.Data = &TilesetData{}
		if err = reader.Decode(tset.Data); err != nil {
			tileFi.Close()
			return nil, fmt.Errorf("decode %s: %w", tset.Source, err)
		}
		tileFi.Close()

		if tset.Data.Image != nil {
			if tset.Data.Image.Image, err = f.loadImage(dirname, tset.Data.Image.Source); err != nil {
				return nil, err
			}
		}

		for idx, tile := range tset.Data.Tile {
			if tile.Image != nil {
				if tile.Image.Image, err = f.loadImage(dirname, tile.Image.Source); err != nil {
					return nil, err
				}
				tset.Data.Tile[idx] = tile
			}
		}

		m.Tileset[idx] = tset
	}

	return &m, nil
}
func (f *MapLoader) loadImage(path, imagePath string) (*TilesetDataImageImage, error) {
	if !strings.HasPrefix(imagePath, "/") {
		imagePath = fmt.Sprintf("%s/%s", path, imagePath)
	}
	imagePath = filepath.Clean(imagePath)
	var (
		err    error
		tImage = &TilesetDataImageImage{}
	)

	if tImage.EImage, tImage.Img, err = ebitenutil.NewImageFromFileSystem(f.fs, imagePath); err != nil {
		return nil, fmt.Errorf("failed load image %s: %w", imagePath, err)
	}

	return tImage, nil
}
