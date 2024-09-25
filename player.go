package main

import (
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleGamepadControls() {
	if g.gamepadIDs == nil {
		g.gamepadIDs = map[ebiten.GamepadID]struct{}{}
	}

	// Log the gamepad connection events.
	g.gamepadIDsBuf = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIDsBuf[:0])
	for _, id := range g.gamepadIDsBuf {
		log.Printf("gamepad connected: id: %d, SDL ID: %s", id, ebiten.GamepadSDLID(id))
		g.gamepadIDs[id] = struct{}{}
	}
	for id := range g.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("gamepad disconnected: id: %d", id)
			delete(g.gamepadIDs, id)
		}
	}

	for id := range g.gamepadIDs {
		// Handle gamepad axis movement
		axisX := ebiten.GamepadAxisValue(id, 0)
		axisY := ebiten.GamepadAxisValue(id, 1)

		op := &ebiten.VibrateGamepadOptions{
			Duration:        200 * time.Millisecond,
			StrongMagnitude: 0.5,
			WeakMagnitude:   0.5,
		}
		ebiten.VibrateGamepad(id, op)

		// Deadzone adjustment
		const deadzone = 0.1
		if math.Abs(axisX) < deadzone {
			axisX = 0
		}
		if math.Abs(axisY) < deadzone {
			axisY = 0
		}

		g.Move(axisX, axisY)
	}
}

func (g *Game) HandleWASDControls() {

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

func (g *Game) UpdatePlayer() {

	// VibrateControllerOnCollision(g.player, g.softColliders, g.gamepadIDs)

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
