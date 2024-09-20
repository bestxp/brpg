package gui

import (
	"image"
	"log"

	e "github.com/hajimehoshi/ebiten/v2"
)

type Entity struct {
	img  *e.Image
	rect image.Rectangle
	X, Y float64

	frame int

	hovered bool
}

// SetFrame sync frames for drawing
func (i *Entity) SetFrame(f int) {
	i.frame = f
}

func (i *Entity) TryHover(x int, y int) bool {
	i.hovered = i.rect.Overlaps(image.Rect(x, y, x+1, y+1))
	return i.hovered
}

func (i *Entity) HandleClick(button e.MouseButton) {
	log.Println("mouse ", button)
}

func (i *Entity) DrawAt(screen *e.Image, x, y float64) {
	if i.img == nil {
		return
	}
	w, h := i.Size()

	i.X = x
	i.Y = y

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	i.rect = image.Rect(int(x), int(y), int(x)+w, int(y)+h)
	screen.DrawImage(i.img, op)
}

func (i *Entity) Size() (int, int) {
	return i.img.Size()
}
