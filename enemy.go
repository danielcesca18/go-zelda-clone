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

	g.enemies = append(g.enemies, &entities.Enemy{
		Sprite: &entities.Sprite{
			Img: enemyImg,
			X:   x,
			Y:   y,
		},
		FollowsPlayer: true,
		Collider:      g.softColliders[len(g.softColliders)-1],
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
	for _, sprite := range g.enemies {

		if g.killEnemies {
			g.killEnemy(sprite)
			continue
		}

		sprite.Dx = 0.0
		sprite.Dy = 0.0

		if sprite.FollowsPlayer {
			// if sprite.FollowsPlayer {
			// Calcular a direção do movimento em relação ao jogador
			directionX := (g.player.X + 8) - (sprite.X + 8)
			directionY := (g.player.Y + 8) - (sprite.Y + 8)
			magnitude := math.Sqrt(directionX*directionX + directionY*directionY)

			if magnitude != 0 {
				// Normalizar o vetor de direção
				directionX /= magnitude
				directionY /= magnitude

				// Aplicar a velocidade do inimigo ao vetor de direção normalizado
				velocity := 0.5 // Ajuste este valor conforme necessário
				sprite.Dx = directionX * velocity
				sprite.Dy = directionY * velocity
			}
		}

		CheckSoftCollision(sprite.Sprite, sprite.Collider, g.softColliders)

		sprite.X += sprite.Dx
		CheckHardCollision(sprite.Sprite, g.hardColliders, X)

		sprite.Y += sprite.Dy
		CheckHardCollision(sprite.Sprite, g.hardColliders, Y)

		sprite.Collider.Rect.MaxX = sprite.X + 16
		sprite.Collider.Rect.MinX = sprite.X
		sprite.Collider.Rect.MaxY = sprite.Y + 16
		sprite.Collider.Rect.MinY = sprite.Y
	}
	if len(g.enemies) == 0 {
		g.killEnemies = false
	}
}

func (g *Game) DrawEnemies(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()
}
