package gui

import (
	e "github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Icon struct {
	*Entity
}

func NewIcon() *Icon {
	i := &Icon{Entity: &Entity{
		img: nil,
	}}
	i.updateIcon()
	return i
}

func (i *Icon) updateIcon() *e.Image {
	if i.img == nil {
		i.img = e.NewImage(32, 32)
	}
	if !i.Entity.hovered {
		i.img.Fill(color.White)
	} else {
		i.img.Fill(color.RGBA{
			R: 100 - uint8(i.frame)*3,
			G: 22 + uint8(i.frame)*3,
			B: 22 + uint8(i.frame)*3,
			A: 255,
		})
	}

	return i.img
}

func (i *Icon) TryHover(x int, y int) bool {
	i.Entity.TryHover(x, y)
	i.updateIcon()
	return i.hovered
}
