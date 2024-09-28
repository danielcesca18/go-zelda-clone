package main

import (
	"go-game/entities"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) spawnEnemy() {

	if g.Tick%60 == 0 && g.spawnEnemies {
		g.newEnemy(g.player.X+100, g.player.Y+100)
		g.newEnemy(g.player.X-100, g.player.Y-+100)
		g.newEnemy(g.player.X-100, g.player.Y+100)
		g.newEnemy(g.player.X+100, g.player.Y-100)
	}
}

func (g *Game) newEnemy(x, y float64) {
	enemyImg, _, err := ebitenutil.NewImageFromFile("assets/images/enemy.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	newCollider := entities.Collider{
		Rect: &entities.FloatRect{
			MinX: x,
			MinY: y,
			MaxX: x + 16,
			MaxY: y + 16,
		},
		Weight: 20,
	}

	g.softColliders = append(g.softColliders, newCollider)

	health := uint(20)
	status := "CHASING"

	g.enemies = append(g.enemies, &entities.Enemy{
		Sprite: &entities.Sprite{
			Img: enemyImg,
			X:   x,
			Y:   y,
		},
		Collider: g.softColliders[len(g.softColliders)-1],
		Health:   &health,
		Status:   &status,
		Knockback: entities.Knockback{
			DirectionX: 0.0,
			DirectionY: 0.0,
			Velocity:   1.5,
		},
	})
}

func (g *Game) killEnemy(enemy *entities.Enemy) {
	for i, e := range g.enemies {
		if e == enemy {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			for j, c := range g.softColliders {
				if c == enemy.Collider {
					g.softColliders = append(g.softColliders[:j], g.softColliders[j+1:]...)
				}
			}
			break
		}
	}

}

func (g *Game) updateEnemies() {
	for _, enemy := range g.enemies {

		if g.killEnemies {
			g.killEnemy(enemy)
			continue
		}

		if *enemy.Health <= 0 {
			g.killEnemy(enemy)
			KillSoundPlayer.Play()
			g.Points++
			continue
		}

		enemy.Dx = 0.0
		enemy.Dy = 0.0

		if *enemy.Status == "CHASING" && g.enemiesFollowsPlayer {
			// Calcular a direção do movimento em relação ao jogador
			directionX := (g.player.X + 8) - (enemy.X + 8)
			directionY := (g.player.Y + 8) - (enemy.Y + 8)
			magnitude := math.Sqrt(directionX*directionX + directionY*directionY)

			if magnitude != 0 {
				// Normalizar o vetor de direção
				directionX /= magnitude
				directionY /= magnitude

				// Aplicar a velocidade do inimigo ao vetor de direção normalizado
				velocity := 0.5 // Ajuste este valor conforme necessário
				enemy.Dx = directionX * velocity
				enemy.Dy = directionY * velocity
			}
		}

		if *enemy.Status == "HIT" {

			if enemy.Knockback.DirectionX == 0 && enemy.Knockback.DirectionY == 0 {
				enemy.Knockback.DirectionX = (g.player.X + 8) - (enemy.X + 8)
				enemy.Knockback.DirectionY = (g.player.Y + 8) - (enemy.Y + 8)
			}
			// Calcular a direção do movimento em relação ao jogador
			directionX := enemy.Knockback.DirectionX
			directionY := enemy.Knockback.DirectionY

			magnitude := math.Sqrt(directionX*directionX + directionY*directionY)

			if magnitude != 0 {
				// Normalizar o vetor de direção
				directionX /= magnitude
				directionY /= magnitude

				// Aplicar a velocidade do inimigo ao vetor de direção normalizado
				enemy.Dx = -directionX * enemy.Knockback.Velocity
				enemy.Dy = -directionY * enemy.Knockback.Velocity
			}

			enemy.HitCounter++
			if enemy.HitCounter >= 15 {
				*enemy.Status = "CHASING"
				enemy.HitCounter = 0
				enemy.Knockback.DirectionX = 0
				enemy.Knockback.DirectionY = 0
			}
		}

		CheckSoftCollision(enemy.Sprite, enemy.Collider, g.softColliders)

		enemy.X += enemy.Dx
		CheckHardCollision(enemy.Sprite, g.hardColliders, X)

		enemy.Y += enemy.Dy
		CheckHardCollision(enemy.Sprite, g.hardColliders, Y)

		enemy.Collider.Rect.MaxX = enemy.X + 16
		enemy.Collider.Rect.MinX = enemy.X
		enemy.Collider.Rect.MaxY = enemy.Y + 16
		enemy.Collider.Rect.MinY = enemy.Y
	}
	if len(g.enemies) == 0 {
		g.killEnemies = false
	}
}

func (g *Game) DrawEnemies(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	for _, sprite := range g.enemies {
		opts.GeoM.Reset()
		opts.ColorScale.Reset()

		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		if *sprite.Status == "HIT" {
			opts.ColorScale.Scale(128, 0, 0, 1)
		}

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}
}
