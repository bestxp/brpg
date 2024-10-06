package maps

import (
	"embed"
	"image"
	_ "image/gif"
	_ "image/png"
	"path"

	"encoding/xml"
	"fmt"
	"strings"
)

type MapLoader struct {
	fs embed.FS
}

func FromFs(fs embed.FS) *MapLoader {
	return &MapLoader{
		fs: fs,
	}
}

func (f *MapLoader) Load(file string) (*Map, error) {
	file = path.Clean(file)
	dirname := path.Dir(file)
	fi, err := f.fs.Open(file)
	if err != nil {
		return nil, fmt.Errorf("read map: %w", err)
	}
	defer fi.Close()

	reader := xml.NewDecoder(fi)

	var m TiledMap
	if err = reader.Decode(&m); err != nil {
		return nil, fmt.Errorf("decode map: %w", err)
	}

	for idx, tset := range m.Tileset {
		if !strings.HasPrefix(tset.Source, "/") {
			tset.Source = path.Join(dirname, tset.Source)
		}
		tset.Source = path.Clean(tset.Source)
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

	return NewMap(m), nil
}
func (f *MapLoader) loadImage(fpath, imagePath string) (*TilesetDataImageImage, error) {
	if !strings.HasPrefix(imagePath, "/") {
		imagePath = path.Join(fpath, imagePath)
	}
	tImage := &TilesetDataImageImage{}

	file, err := f.fs.Open(path.Clean(imagePath))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	tImage.Img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}

	return tImage, nil
}
