package gui

import e "github.com/hajimehoshi/ebiten/v2"

type Condition uint8

const (
	Idle Condition = iota
	Run
)

type Gui interface {
	DrawAt(image *e.Image, x, y float64)
	SetFrame(int)
}
