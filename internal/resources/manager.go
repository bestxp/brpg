package resources

import (
	"fmt"
	"image"
	"io/fs"
	"path"
	"strconv"
	"strings"

	_ "image/png"

	types "github.com/bestxp/brpg/internal/resources/yaml"
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

	texture.Images = make(map[string]*Image, len(t.Images))
	for _, i := range t.Images {
		img, err := m.imageFromSource(i.Source)
		if err != nil {
			return nil, err
		}
		if i.Tiles != nil {
			img.SetTileInfo(i.Tiles.Tile.W, i.Tiles.Tile.H, i.Tiles.Cols, i.Tiles.Rows)
		}
		texture.Images[i.Name] = img
	}

	texture.Animatons = make(map[string]*Animation, len(t.Animations))

	for _, a := range t.Animations {
		frames := make([]*Image, 0, len(a.Frames))
		for _, fr := range a.Frames {
			inf, err := m.parseLink(fr)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", a.Name, err)
			}
			if inf.Section != "images" {
				continue
			}
			i := texture.Images[inf.Name]
			if inf.Tiled && i.IsTiled() {
				if i, err = i.Tile(int(inf.Row), int(inf.Col)); err != nil {
					return nil, fmt.Errorf("tiled %s: %w", a.Name, err)
				}
			}
			frames = append(frames, i)
		}
		texture.Animatons[a.Name] = NewAnimation(frames, uint16(a.Duration), WithFlip(a.Flip), WithLoop(a.Loop))
	}

	texture.Entities = make(map[string]*Entity, len(t.Entities))

	for _, e := range t.Entities {
		ent := NewEntity(uint16(e.Size.W), uint16(e.Size.H))

		for _, a := range e.Actions {
			inf, err := m.parseLink(a.Sprite)
			if err != nil {
				return nil, fmt.Errorf("failed load sprite entity %s, %s: %w", e.Name, a.Name, err)
			}
			var act SpriteFramed
			if inf.Section == "image" {
				if texture.Images[inf.Name].IsTiled() && inf.Tiled {
					if act, err = texture.Images[inf.Name].Tile(int(inf.Row), int(inf.Col)); err != nil {
						return nil, fmt.Errorf("failed get sprite entity tile %s, %s: %w", e.Name, a.Name, err)
					}
				} else {
					act = texture.Images[inf.Name]
				}
			}
			if inf.Section == "animations" {
				act = texture.Animatons[inf.Name]
			}
			ent.Actions[a.Name] = act
		}

		texture.Entities[e.Name] = ent
	}

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

func (m *Manager) parseLink(source string) (*imageInfo, error) {
	if !strings.HasPrefix(source, "#") {
		return nil, nil
	}
	source = strings.TrimPrefix(source, "#")
	parts := strings.Split(source, ".")

	inf := &imageInfo{}

	inf.Section = parts[0]
	inf.Name = parts[1]

	var row, col int
	var err error
	if strings.Contains(inf.Name, "[") {
		inf.Tiled = true
		index := strings.Index(inf.Name, "[")
		tile := inf.Name[index+1 : len(inf.Name)-1]
		inf.Name = inf.Name[:index]

		tileCoors := strings.Split(tile, ",")
		if row, err = strconv.Atoi(tileCoors[0]); err != nil {
			return nil, fmt.Errorf("failed parce row %s: %w", tileCoors[0], err)
		}
		if col, err = strconv.Atoi(tileCoors[1]); err != nil {
			return nil, fmt.Errorf("failed parce cols %s: %w", tileCoors[1], err)
		}
		inf.Row = uint16(row)
		inf.Col = uint16(col)
	}

	return inf, nil
}

type imageInfo struct {
	Section, Name string
	Tiled         bool
	Col, Row      uint16
}

type Textures struct {
	Images    map[string]*Image
	Animatons map[string]*Animation
	Entities  map[string]*Entity
}

type Entity struct {
	W, H    uint16
	Actions map[string]SpriteFramed
}

type SpriteFramed interface {
	Frames() []*Image
}

func NewEntity(w, h uint16) *Entity {
	e := &Entity{
		W:       w,
		H:       h,
		Actions: map[string]SpriteFramed{},
	}

	return e
}
