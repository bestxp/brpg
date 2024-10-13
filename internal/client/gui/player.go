package gui

import (
	"image/color"

	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

var frames map[string]resources.Frames

func init() {
	var err error
	frames, err = resources.LoadResources()
	if err != nil {
		panic(err)
	}
}

type Player struct {
	*Entity

	Condition Condition

	Dx, Dy float64

	unit *pkg.Unit

	texture *resources.Entity
}

func NewPlayerFromServer(unit *pkg.Unit, e *resources.Textures) *Player {
	p := &Player{
		Entity: &Entity{},
	}
	p.X = unit.Pos.X
	p.Y = unit.Pos.Y
	p.unit = unit

	p.texture = e.Entities[unit.Skin]
	if p.texture == nil {
		log.Debug().Msg("skin not loaded")
	}
	return p
}

func (p *Player) updateImg() {
	if p.img == nil {
		p.img = e.NewImage(int(p.texture.W), int(p.texture.H))

		op := &e.DrawImageOptions{}

		var action = "idle"
		if p.unit.Action == "run" {
			switch p.unit.Direction {
			case pkg.Direction_down:
				action = "walk_down"
			case pkg.Direction_left:
				action = "walk_left"
			case pkg.Direction_right:
				action = "walk_right"
			}
		}

		frames := p.texture.Actions[action].Frames()

		skin := e.NewImageFromImage(frames[(p.frame/7)%len(frames)].Image())
		p.img.DrawImage(skin, op)

		var (
			percentBarWidth   = int(p.texture.W - 6)
			percentBarHeight  = 4
			percentBarPadding = 1
		)

		healthBar := e.NewImage(percentBarWidth, percentBarHeight)
		healthBar.Fill(color.RGBA{A: 255})

		percent := float64(p.unit.Info.CurrentHealth) / float64(p.unit.Info.MaxHealth)
		percent = float64(percentBarWidth-percentBarPadding*2) * percent

		currHealth := e.NewImage(int(percent), percentBarHeight-percentBarPadding*2)
		currHealth.Fill(color.RGBA{R: 255, A: 255})

		currHealthOp := &e.DrawImageOptions{}
		currHealthOp.GeoM.Translate(float64(percentBarPadding), float64(percentBarPadding))

		healthBar.DrawImage(currHealth, currHealthOp)

		hbOpt := &e.DrawImageOptions{}
		hbOpt.GeoM.Translate(float64(int(p.texture.W)-percentBarWidth)/2, 0)
		p.img.DrawImage(healthBar, hbOpt)
	}
}

func (p *Player) DrawAt(screen *e.Image, x, y float64) {
	p.updateImg()
	p.Entity.DrawAt(screen, x, y)
}

func (p *Player) TryHover(x int, y int) bool {
	v := p.Entity.TryHover(x, y)

	return v
}
