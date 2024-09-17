package level

import (
	"errors"
	"github.com/bestxp/brpg/internal/resources"
	e "github.com/hajimehoshi/ebiten/v2"
	"math"
	"sync"
)

const baseTileSize int = 16
const baseUnitSize float64 = 16
const baseScale int = 2

var errOutOfBound = errors.New("out of map bound")

type LevelName string

type Pos struct {
	X, Y float64
}

type Vector struct {
	From, To Pos
}

type Level struct {
	scale int

	ID   LevelName
	Name string
	_map [][]Tile

	StartPos Pos // start location

	once  sync.Once
	image *e.Image

	frames map[string]resources.Frames
}

type Tile struct {
	CanWalk bool
	Texture string

	TopLeft Pos
}

func NewLevel(
	ID LevelName,
	name string,
	frames map[string]resources.Frames,
) *Level {
	return &Level{
		ID:     ID,
		_map:   [][]Tile{},
		scale:  baseScale,
		Name:   name,
		frames: frames,
	}
}

func (l LevelName) String() string {
	return string(l)
}

func (l *Level) Map() [][]Tile {
	return l._map
}

func (l *Level) SetMap(m [][]Tile) {
	for i, row := range m {
		for j, t := range row {
			w, h := t.Size()
			m[i][j].TopLeft.X = float64(w * j)
			m[i][j].TopLeft.Y = float64(h * i)
		}
	}

	l._map = m
}

func (l *Level) unitSize() float64 {
	return baseUnitSize * float64(baseScale)
}

func (l *Level) getTileByCoords(coors Pos) (int, int, Tile, error) {
	y := int(math.Floor(coors.Y / float64(baseTileSize*baseScale)))
	x := int(math.Floor(coors.X / float64(baseTileSize*baseScale)))

	if len(l._map) <= y {
		return 0, 0, Tile{}, errOutOfBound
	}
	if len(l._map[y]) <= x {
		return 0, 0, Tile{}, errOutOfBound
	}

	return y, x, l._map[y][x], nil
}

func (l *Level) WalkCalc(vector Vector) Pos {

	from := vector.From
	toCoords := vector.To

	// add shifts for calculate end coords
	if vector.From.Y < vector.To.Y {
		vector.From.Y += l.unitSize()
	}

	if vector.From.X < vector.To.X {
		vector.From.X += l.unitSize()
	}

	_, _, tileFrom, err := l.getTileByCoords(vector.From)
	if err != nil {
		return from
	}

	_, _, tileTo, err := l.getTileByCoords(vector.To)
	if err != nil {
		return from
	}

	if !tileFrom.CanWalk {
		toCoords = from
	}

	if !tileTo.CanWalk {
		toCoords = from
	}

	return toCoords
}

func (l *Level) LevelSize() (int, int, int, int) {
	if len(l._map) == 0 {
		return 0, 0, baseTileSize * l.scale, baseTileSize * l.scale
	}

	yBound := len(l._map) * baseTileSize * l.scale
	xBound := len(l._map[0]) * baseTileSize * l.scale

	return 0, 0, xBound, yBound
}

func (l *Level) Scale() int {
	return l.scale
}

// Size returns width, height
func (t Tile) Size() (int, int) {
	return baseTileSize, baseTileSize
}

func (l *Level) EImage() (*e.Image, error) {
	if l.image != nil {
		return l.image, nil
	}
	_, _, width, height := l.LevelSize()

	levelImage := e.NewImage(width, height)

	for _, row := range l._map {
		for _, tile := range row {
			op := &e.DrawImageOptions{}
			op.GeoM.Translate(tile.TopLeft.X, tile.TopLeft.Y)
			op.GeoM.Scale(float64(l.Scale()), float64(l.Scale()))

			img := e.NewImageFromImage(l.frames[tile.Texture].Frames[0])
			levelImage.DrawImage(img, op)
		}
	}
	l.image = levelImage

	return l.image, nil
}
