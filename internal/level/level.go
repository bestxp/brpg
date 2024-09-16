package level

import (
	"errors"
	engine "github.com/bestxp/brpg"
	e "github.com/hajimehoshi/ebiten/v2"
	"math"
	"sync"
)

const baseTileSize int = 16
const baseUnitSize float64 = 16
const baseScale int = 2

var errOutOfBound = errors.New("out of map bound")

type Pos struct {
	X, Y float64
}

type Vector struct {
	From, To Pos
}

type Level struct {
	scale int

	Name string
	_map [][]Tile

	StartPos Pos // start location

	once  sync.Once
	image *e.Image
}

type Tile struct {
	CanWalk bool
	Texture string

	TopLeft Pos
}

func NewLevel(name string) *Level {
	return &Level{_map: [][]Tile{}, scale: baseScale, Name: name}
}

func (l *Level) Map() [][]Tile {
	return l._map
}

func (l *Level) SetMap(m [][]Tile) {
	for i, row := range m {
		for j, t := range row {
			w, h := t.Size()
			m[i][j].TopLeft.X = float64(w * i)
			m[i][j].TopLeft.Y = float64(h * j)
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
	topX, topY, bottomX, bottomY := l.LevelSize()

	from := vector.From
	toCoords := vector.To
	walkStopped := false

	walkX := false
	walkY := false

	// add shifts for calculate end coords
	if vector.From.Y < vector.To.Y {
		vector.From.Y += l.unitSize()
		walkY = true
	} else if vector.From.Y > vector.To.Y {
		walkY = true
	}

	if vector.From.X < vector.To.X {
		vector.From.X += l.unitSize()
		walkX = true
	} else if vector.From.X > vector.To.X {
		walkX = true
	}

	if float64(topX) >= vector.To.X && !walkY {
		toCoords.X = float64(topX)
		walkStopped = true
	}
	if float64(bottomX) <= vector.To.X && !walkY {
		toCoords.X = float64(bottomX)
		walkStopped = true
	}

	if float64(topY) >= vector.To.Y && !walkX {
		toCoords.Y = float64(topY)
		walkStopped = true
	}

	if float64(bottomY) <= vector.To.Y && !walkX {
		toCoords.Y = float64(bottomY)
		walkStopped = true
	}

	if walkStopped {
		return toCoords
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

func (l *Level) EImage(frames map[string]engine.Frames) (*e.Image, error) {
	_, _, width, height := l.LevelSize()

	levelImage := e.NewImage(width, height)

	for _, row := range l._map {
		for _, tile := range row {
			op := &e.DrawImageOptions{}
			op.GeoM.Translate(tile.TopLeft.X, tile.TopLeft.Y)
			op.GeoM.Scale(float64(l.Scale()), float64(l.Scale()))

			img := e.NewImageFromImage(frames[tile.Texture].Frames[0])
			levelImage.DrawImage(img, op)
		}
	}
	l.image = levelImage

	return l.image, nil
}
