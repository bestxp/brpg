package gui

import (
	"image"
	"image/color"

	e "github.com/hajimehoshi/ebiten/v2"
)

type Icon struct {
	*Entity
	icon *e.Image
}

func NewIcon(img image.Image) *Icon {
	i := &Icon{Entity: &Entity{
		img: nil,
	}}
	if img != nil {
		i.icon = e.NewImageFromImage(img)
	}
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

	i.img.DrawImage(i.icon, &e.DrawImageOptions{})

	return i.img
}

func (i *Icon) TryHover(x int, y int) bool {
	i.Entity.TryHover(x, y)
	i.updateIcon()
	return i.hovered
}
