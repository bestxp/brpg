package gui

import (
	"bytes"
	"image"

	"github.com/bestxp/brpg"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/colornames"
	"golang.org/x/text/language"
)

type Button struct {
	*Entity

	TextOps *text.DrawOptions
	w       int
	h       int
	label   string
}

func (b *Button) updateImg() {
	button := e.NewImage(b.w, b.h)

	// todo color or img
	if !b.Entity.hovered {
		button.Fill(colornames.Burlywood)
	} else {
		clr := colornames.Burlywood
		clr.R += 25
		clr.G += 10
		clr.B += 25
		button.Fill(clr)
	}

	f, err := text.NewGoTextFaceSource(bytes.NewBuffer(brpg.MainFont))
	if err != nil {
		log.Error().Err(err)
	}

	font := &text.GoTextFace{
		Source:   f,
		Size:     30,
		Language: language.AmericanEnglish,
	}
	wf, hf := text.Measure(b.label, font, font.Size*1.5)

	op := *b.TextOps

	op.GeoM.Translate(float64(b.w/2)-wf/2, float64(b.h/2)-hf/2)
	text.Draw(button, b.label, font, &op)
	b.Entity.img = button
}

func NewButton(value string, w, h int) *Button {
	b := &Button{
		Entity: &Entity{
			img:     nil,
			rect:    image.Rectangle{},
			X:       0,
			Y:       0,
			frame:   0,
			hovered: false,
		},
		w: w,
		h: h,

		label: value,

		TextOps: &text.DrawOptions{},
	}
	b.updateImg()

	return b
}

// DrawAt draws button at center of scene
func (b *Button) DrawAt(screen *e.Image, x, y float64) {
	b.Entity.DrawAt(screen, float64(screen.Bounds().Dx()/2-b.w/2), float64(screen.Bounds().Dy()/2-b.h/2))
}

func (b *Button) TryHover(x, y int) bool {
	defer b.updateImg()
	return b.Entity.TryHover(x, y)
}
