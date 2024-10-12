package resources

import (
	"fmt"
	"image"
	"io/fs"
	"path"

	_ "image/png"

	types "github.com/bestxp/brpg/internal/resources/yaml"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	pathPart = "./resources/textures"
	fileExt  = ".yml"
)

type (
	Manager struct {
		fs               fs.FS
		withoutResources bool
		path             string
	}
)

func FromFS(fsystem fs.FS, withoutResources bool) *Manager {
	return &Manager{
		fs:               fsystem,
		withoutResources: withoutResources,
		path:             path.Clean(pathPart),
	}
}

func (m *Manager) Load(texture string) (*Textures, error) {
	file := path.Clean(fmt.Sprintf("%s/%s%s", m.path, texture, fileExt))

	f, err := m.fs.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed open %s:%w", file, err)
	}

	defer f.Close()
	t := types.Texture{}

	if err = yaml.NewDecoder(f).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed decode %s: %w", file, err)
	}

	return m.Collect(t)
}

func (m *Manager) Collect(t types.Texture) (*Textures, error) {
	texture := &Textures{}

	images := make(map[string]*Image, len(t.Images))
	for _, i := range t.Images {
		img, err := m.imageFromSource(i.Source)
		if err != nil {
			return nil, err
		}
		if i.Tiles != nil {
			img.SetTileInfo(i.Tiles.Tile.W, i.Tiles.Tile.H, i.Tiles.Cols, i.Tiles.Rows)
		}
		images[i.Name] = img
	}

	texture.Images = images

	return texture, nil
}

func (m *Manager) imageFromSource(fpath string) (*Image, error) {
	file := path.Clean(fmt.Sprintf("%s/%s", m.path, fpath))
	f, err := m.fs.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed read file %s: %w", file, err)
	}
	im, tt, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed decode image %s: %w", file, err)

	}

	if nr, ok := im.(*image.NRGBA); ok {
		return NewImage(nr), nil
	}
	return nil, fmt.Errorf("unknown type for %s (%s)", file, tt)
}

type Textures struct {
	Images map[string]*Image
}

func NewImage(i *image.NRGBA) *Image {
	return &Image{image: i, size: [2]uint16{uint16(i.Bounds().Dx()), uint16(i.Bounds().Dy())}}
}

type Image struct {
	image *image.NRGBA
	size  [2]uint16

	tiles Tiles
}

func (i *Image) Tile(row, col int) (*image.NRGBA, error) {
	if i.IsTiled() {
		return i.tiles[row][col], nil
	}
	return nil, fmt.Errorf("no tiled resource")
}

type Tiles map[int]map[int]*image.NRGBA

func (i *Image) Width() int {
	return int(i.size[0])
}

func (i *Image) Heigth() int {
	return int(i.size[1])
}

func (i *Image) IsTiled() bool {
	return len(i.tiles) > 0
}

func (i *Image) TilesRows() int {
	return len(i.tiles)
}

func (i *Image) TilesCount() int {
	if !i.IsTiled() {
		return 0
	}
	return len(i.tiles[0])
}

func (i *Image) SetTileInfo(w, h, cols, rows int) {
	i.tiles = make(Tiles, cols)
	if i.image == nil {
		log.Debug().Msgf("no image loaded")
		return
	}

	for j := 1; j <= rows; j++ {
		i.tiles[j] = make(map[int]*image.NRGBA, cols)
		for col := 1; col <= cols; col++ {
			rect := image.Rect(w*(col-1), h*(j-1), col*w, h*j)
			log.Debug().Str("rect", rect.String()).Msg("sub rect")
			i.tiles[j][col] = i.image.SubImage(rect).(*image.NRGBA)
		}
	}
}
