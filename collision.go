package main

import (
	"go-game/entities"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// make soft colision
// make hard colision

// enum for axis x and y
const (
	X = "x"
	Y = "y"
)

type Axis string

func CheckSoftCollision(sprite *entities.Sprite, ownCollider entities.Collider, colliders []entities.Collider) {
	for _, collider := range colliders {
		// Verificar se é o próprio colisor
		if ownCollider == collider {
			continue
		}

		if collider.Rect.Overlaps(
			entities.FloatRect{
				MinX: sprite.X,
				MinY: sprite.Y,
				MaxX: sprite.X + 16.0,
				MaxY: sprite.Y + 16.0,
			},
		) {
			// Calcular a área de sobreposição
			overlapMinX := math.Max(sprite.X, float64(collider.Rect.MinX))
			overlapMinY := math.Max(sprite.Y, float64(collider.Rect.MinY))
			overlapMaxX := math.Min(sprite.X+16.0, float64(collider.Rect.MaxX))
			overlapMaxY := math.Min(sprite.Y+16.0, float64(collider.Rect.MaxY))

			overlapWidth := overlapMaxX - overlapMinX
			overlapHeight := overlapMaxY - overlapMinY
			overlapArea := overlapWidth * overlapHeight

			// Calcular o vetor de direção
			spriteCenterX := sprite.X + 8.0
			spriteCenterY := sprite.Y + 8.0
			colliderCenterX := float64(collider.Rect.MinX) + 8.0
			colliderCenterY := float64(collider.Rect.MinY) + 8.0

			directionX := colliderCenterX - spriteCenterX
			directionY := colliderCenterY - spriteCenterY
			magnitude := math.Sqrt(directionX*directionX + directionY*directionY)

			if magnitude != 0 {
				// Normalizar o vetor de direção
				directionX /= magnitude
				directionY /= magnitude

				// Aplicar uma força de repulsão proporcional à área de sobreposição
				repulsionStrength := overlapArea * (collider.Weight / 1000) // Ajuste este valor conforme necessário
				repulsionX := directionX * repulsionStrength
				repulsionY := directionY * repulsionStrength

				// Atualizar a posição do sprite para suavizar o movimento
				sprite.Dx -= repulsionX
				sprite.Dy -= repulsionY
			}
		}
	}
}

func (g *Game) IsTouchingPlayer(spriteCollider entities.Collider, player entities.Player) bool {
	if spriteCollider.Rect.Overlaps(
		entities.FloatRect{
			MinX: player.X,
			MinY: player.Y,
			MaxX: player.X + 16.0,
			MaxY: player.Y + 16.0,
		},
	) {
		if !player.Invencible {
			*g.player.Health -= 1
			g.player.Invencible = true
			PlayerHitSoundPlayer.Play()
		}
		return true
	}

	return false
}

func IsSurrounded(sprite *entities.Sprite, ownCollider entities.Collider, colliders []entities.Collider, margin float64) bool {
	collidedTop := false
	collidedBottom := false
	collidedLeft := false
	collidedRight := false

	for _, collider := range colliders {
		// Verificar se é o próprio colisor
		if ownCollider == collider {
			continue
		}

		if collider.Rect.Overlaps(
			entities.FloatRect{
				MinX: sprite.X,
				MinY: sprite.Y,
				MaxX: sprite.X + 16.0,
				MaxY: sprite.Y + 16.0,
			},
		) {
			// Calcular a direção do movimento em relação ao colisor
			spriteCenterX := sprite.X + 8.0
			spriteCenterY := sprite.Y + 8.0
			colliderCenterX := float64(collider.Rect.MinX) + 8.0
			colliderCenterY := float64(collider.Rect.MinY) + 8.0

			directionX := colliderCenterX - spriteCenterX
			directionY := colliderCenterY - spriteCenterY

			// Verificar colisão com margem de erro
			if math.Abs(directionY) < margin && directionX > 0 {
				collidedRight = true
			}
			if math.Abs(directionY) < margin && directionX < 0 {
				collidedLeft = true
			}
			if math.Abs(directionX) < margin && directionY > 0 {
				collidedBottom = true
			}
			if math.Abs(directionX) < margin && directionY < 0 {
				collidedTop = true
			}
		}
	}

	return collidedTop && collidedBottom && collidedLeft && collidedRight
}

func CheckHardCollision(sprite *entities.Sprite, colliders []image.Rectangle, axis Axis) {
	if axis == X {
		CheckHardCollisionHorizontal(sprite, colliders)
	} else if axis == Y {
		CheckHardCollisionVertical(sprite, colliders)
	}
}

func CheckHardCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckHardCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func (g *Game) DrawColliders(screen *ebiten.Image) {
	if g.showColliders {
		for _, collider := range g.hardColliders {
			vector.StrokeRect(
				screen,
				float32(collider.Min.X)+float32(g.cam.X),
				float32(collider.Min.Y)+float32(g.cam.Y),
				float32(collider.Dx()),
				float32(collider.Dy()),
				1.0,
				color.RGBA{255, 0, 0, 255},
				true,
			)
		}

		for _, collider := range g.softColliders {
			vector.StrokeRect(
				screen,
				float32(collider.Rect.MinX)+float32(g.cam.X),
				float32(collider.Rect.MinY)+float32(g.cam.Y),
				float32(float32(collider.Rect.MaxX-collider.Rect.MinX)),
				float32(float32(collider.Rect.MaxY-collider.Rect.MinY)),
				1.0,
				color.RGBA{0, 0, 255, 255},
				true,
			)
		}
	}
}
