package camera

import (
	"math"
)

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) InitCoords(x, y float64) {
	c.X = x
	c.Y = y
}

// FollowTarget sets the position of the camera based on the position of the target and the size of the screen
func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight, shiftx, shifty float64) {
	c.X = -targetX + screenWidth/2.0 - shiftx
	c.Y = -targetY + screenHeight/2.0 - shifty
}

// Constrain stops the camera from showing past the boundaries of the tilemap
func (c *Camera) Constrain(
	tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight float64,
) {
	c.X = math.Max(c.X, screenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, screenHeight-tilemapHeightPixels)

	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)
}
