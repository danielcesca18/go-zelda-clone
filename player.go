package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleControls() {

	// Player attack
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Attack()
	}

	// Player movement
	var directionX, directionY float64
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		directionX = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		directionX = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		directionY = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		directionY = 1
	}
	g.Move(directionX, directionY)

	// Misc

	// music volume
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.SetVolume(false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.SetVolume(true)
	}

	// enemies follows player
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		if !g.enemiesFollowsPlayer {
			g.enemiesFollowsPlayer = true
		} else {
			g.enemiesFollowsPlayer = false
		}
	}

	// auto spawn enemies
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if !g.spawnEnemies {
			g.spawnEnemies = true
		} else {
			g.spawnEnemies = false
		}
	}

	// spawn enemy
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		randomOffsetX := float64(rand.Intn(201) - 100)
		randomOffsetY := float64(rand.Intn(201) - 100)

		g.newEnemy(g.player.X+randomOffsetX, g.player.Y+randomOffsetY)
	}

	// kill all enemies
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.killEnemies = true
	}

	// show colliders
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if !g.showColliders {
			g.showColliders = true
		} else {
			g.showColliders = false
		}
	}

}

func (g *Game) Move(directionX, directionY float64) {
	magnitude := math.Sqrt(directionX*directionX + directionY*directionY)

	if magnitude != 0 {
		// Normalizar o vetor de direção
		directionX /= magnitude
		directionY /= magnitude

		// Aplicar a velocidade do jogador ao vetor de direção normalizado
		velocity := 2.0 // Ajuste este valor conforme necessário
		g.player.Dx = directionX * velocity
		g.player.Dy = directionY * velocity
	} else {
		g.player.Dx = 0
		g.player.Dy = 0
	}
}

func (g *Game) Attack() {
	if g.player.Status == "IDLE" {
		HitSoundPlayer.Play()

		g.player.Status = "ATTACK"
		g.attackCounter = 0 // Reset the tick counter to start the animation from the beginning

		for _, enemy := range g.enemies {
			if g.player.Hitbox.Overlaps(enemy.Sprite) {
				*enemy.Health -= g.player.Attack.Damage
				*enemy.Status = "HIT"
			}
		}
	}
}

func (g *Game) UpdatePlayer() {

	CheckSoftCollision(g.player.Sprite, g.player.Collider, g.softColliders)

	g.player.X += g.player.Dx
	CheckHardCollision(g.player.Sprite, g.hardColliders, X)
	g.player.Y += g.player.Dy
	CheckHardCollision(g.player.Sprite, g.hardColliders, Y)

	g.player.Collider.Rect.MaxX = g.player.X + 16
	g.player.Collider.Rect.MinX = g.player.X
	g.player.Collider.Rect.MaxY = g.player.Y + 16
	g.player.Collider.Rect.MinY = g.player.Y
}

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 4
)

func (g *Game) DrawPlayer(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	// set the translation of our drawImageOptions to the player's position
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	// draw the player
	screen.DrawImage(
		// grab a subimage of the spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()
}

func (g *Game) DrawAttack(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	// Draw attack image
	if g.player.Status == "ATTACK" {

		// Reset GeoM
		opts.GeoM.Reset()

		// Translate to the hitbox position
		opts.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)                 // Center the image
		opts.GeoM.Rotate(math.Atan2(g.player.Hitbox.DirectionY, g.player.Hitbox.DirectionX)) // Rotate around the center
		opts.GeoM.Translate(g.player.Hitbox.X, g.player.Hitbox.Y)                            // Move to the hitbox position
		opts.GeoM.Translate(g.cam.X, g.cam.Y)                                                // Adjust for camera position

		// Calculate the frame to draw
		i := (g.attackCounter / 5) % frameCount
		sx, sy := frameOX+i*frameWidth, frameOY

		// Draw the attack image
		screen.DrawImage(g.player.Attack.Img.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), &opts)

		// Reset GeoM
		opts.GeoM.Reset()

		// Reset player status to IDLE after the last frame
		if i == frameCount-1 {
			g.player.Status = "IDLE"
		}
	}
}
