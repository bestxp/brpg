package resources

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
)

type Tiles map[int]map[int]*image.NRGBA
type Image struct {
	image *image.NRGBA
	size  [2]uint16

	tiles Tiles
}

func NewImage(i *image.NRGBA) *Image {
	return &Image{image: i, size: [2]uint16{uint16(i.Bounds().Dx()), uint16(i.Bounds().Dy())}}
}

func (i *Image) Tile(row, col int) (*Image, error) {
	if i.IsTiled() {
		return NewImage(i.tiles[row][col]), nil
	}
	return nil, fmt.Errorf("no tiled resource")
}

func (i *Image) Image() image.Image {
	return i.image
}

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

func (i *Image) TilesCols() int {
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
			i.tiles[j][col] = i.image.SubImage(rect).(*image.NRGBA)
		}
	}
}

func (i *Image) Frames() []*Image {
	return []*Image{i}
}

func (i *Image) Copy() *Image {
	c := *i

	c.image = imaging.Clone(i.image)
	return &c
}

type Animation struct {
	frames   []*Image
	Duration uint16
	Loop     bool
}

func NewAnimation(frames []*Image, duration uint16, opts ...AnimationOpt) *Animation {
	a := &Animation{
		frames:   frames,
		Duration: duration,
		Loop:     true,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

type AnimationOpt func(a *Animation)

func WithFlip(flip string) func(a *Animation) {
	return func(a *Animation) {
		for idx := range a.frames {
			frame := a.frames[idx].Copy()
			if flip == "horizontal" {
				frame.image = imaging.FlipH(a.frames[idx].image)
			}
			if flip == "vertical" {
				frame.image = imaging.FlipV(a.frames[idx].image)
			}
			a.frames[idx] = frame
		}
	}
}

func WithLoop(loop bool) AnimationOpt {
	return func(a *Animation) {
		a.Loop = loop
	}
}

func (a *Animation) Frames() []*Image {
	return a.frames
}
