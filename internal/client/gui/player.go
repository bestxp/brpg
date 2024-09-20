package gui

import (
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
	"image/color"
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
	skin resources.Frames
}

func NewPlayerFromServer(unit *pkg.Unit) *Player {
	p := &Player{
		Entity: &Entity{},
	}
	p.X = unit.Pos.X
	p.Y = unit.Pos.Y
	p.unit = unit
	p.skin = frames[unit.Skin+"_"+unit.Action]

	return p
}

func (p *Player) updateImg() {
	if p.Entity.img == nil {
		p.Entity.img = e.NewImage(p.skin.Config.Width, p.skin.Config.Height)

		op := &e.DrawImageOptions{}
		if p.unit.Direction == pkg.Direction_left {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(p.skin.Config.Width), 0)
		}
		skin := e.NewImageFromImage(p.skin.Frames[(p.frame/7+int(p.unit.Frame))%4])
		p.Entity.img.DrawImage(skin, op)

		var (
			percentBarWidth   = p.skin.Config.Width - 6
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
		hbOpt.GeoM.Translate(float64(p.skin.Config.Width-percentBarWidth)/2, 0)
		p.Entity.img.DrawImage(healthBar, hbOpt)
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
