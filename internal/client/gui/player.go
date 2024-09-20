package gui

import (
	"github.com/bestxp/brpg/internal/resources"
	"github.com/bestxp/brpg/pkg"
	e "github.com/hajimehoshi/ebiten/v2"
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
		p.img.DrawImage(skin, op)
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
