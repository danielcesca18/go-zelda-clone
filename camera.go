package main

import "math"

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

// sets the position of the camera based on the position of the target and the size of the screen
func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.X = -targetX + screenWidth/2.0
	c.Y = -targetY + screenHeight/2.0
}

// stops the camera from showing past the boundaries of the tilemap
func (c *Camera) Constrain(
	tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight float64,
) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, screenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, screenHeight-tilemapHeightPixels)
}

// function that calls the other two functions
func (g *Game) UpdateCamera() {
	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Height)*16.0,
		320,
		240,
	)
	if g.CameraShakeCounter < 12 {
		g.cam.X += math.Sin(float64(g.Tick)*0.8) * 1.5
		g.cam.Y += math.Cos(float64(g.Tick)*0.8) * 1.5
	}
}
