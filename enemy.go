package main

import (
	"go-game/entities"
	"image"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) spawnEnemy() {
	if g.Tick%180 == 0 && g.spawnEnemies {
		// Base number of enemies to spawn
		baseEnemies := 1

		// Increase the number of enemies based on the player's score
		additionalEnemies := g.Tick / 600

		totalEnemies := baseEnemies + additionalEnemies

		for i := 0; i < totalEnemies; i++ {
			spawnDistance := float64(rand.Intn(10)+14) * 16 // Distance from the player to spawn enemies

			// Gerar um ponto aleatório no mapa
			randomX := float64(rand.Intn(1000)) * 16
			randomY := float64(rand.Intn(1000)) * 16

			// Randomly transform randomX and randomY to positive or negative
			if rand.Intn(2) == 0 {
				randomX = -randomX
			}
			if rand.Intn(2) == 0 {
				randomY = -randomY
			}

			// Calcular o vetor de direção entre o jogador e o ponto aleatório
			dirX := randomX - g.player.X
			dirY := randomY - g.player.Y

			// Normalizar o vetor de direção
			length := math.Sqrt(dirX*dirX + dirY*dirY)
			dirX /= length
			dirY /= length

			// Calcular a posição de spawn baseada na direção e na distância fixa
			offsetX := g.player.X + spawnDistance*dirX
			offsetY := g.player.Y + spawnDistance*dirY

			g.newEnemy(offsetX, offsetY)
		}
	}
}

func (g *Game) newEnemyType(x, y float64, enemyType int) {
	var enemyImg *ebiten.Image
	var weight float64
	var health int
	var potionSpawnRate int
	var err error

	if enemyType == 1 {
		enemyImg, _, err = ebitenutil.NewImageFromFile("assets/images/enemy1.png")
		if err != nil {
			log.Fatal(err)
		}
		weight = 20
		health = 15
		potionSpawnRate = 2
	} else if enemyType == 2 {
		enemyImg, _, err = ebitenutil.NewImageFromFile("assets/images/enemy2.png")
		if err != nil {
			log.Fatal(err)
		}
		weight = 20
		health = 30
		potionSpawnRate = 5
	} else if enemyType == 3 {
		enemyImg, _, err = ebitenutil.NewImageFromFile("assets/images/enemy3.png")
		if err != nil {
			log.Fatal(err)
		}
		weight = 40
		health = 45
		potionSpawnRate = 8
	} else if enemyType == 4 {
		enemyImg, _, err = ebitenutil.NewImageFromFile("assets/images/enemy4.png")
		if err != nil {
			log.Fatal(err)
		}
		weight = 60
		health = 60
		potionSpawnRate = 12
	} else if enemyType == 5 {
		enemyImg, _, err = ebitenutil.NewImageFromFile("assets/images/enemy5.png")
		if err != nil {
			log.Fatal(err)
		}
		weight = 80
		health = 80
		potionSpawnRate = 15
	}

	newCollider := entities.Collider{
		Rect: &entities.FloatRect{
			MinX: x,
			MinY: y,
			MaxX: x + 16,
			MaxY: y + 16,
		},
		Weight: weight,
	}

	g.softColliders = append(g.softColliders, newCollider)

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
		PotionSpawnRate: potionSpawnRate,
	})
}

func (g *Game) newEnemy(x, y float64) {
	// Define the probabilities for each enemy type based on the game tick
	// 2000 ticks is around 30 seconds
	var enemyType int
	switch {
	case g.Tick > 5000:
		if rand.Intn(100) < 20 {
			enemyType = 5
		} else if rand.Intn(100) < 40 {
			enemyType = 4
		} else if rand.Intn(100) < 60 {
			enemyType = 3
		} else {
			enemyType = 2
		}
	case g.Tick > 4000:
		if rand.Intn(100) < 20 {
			enemyType = 4
		} else if rand.Intn(100) < 40 {
			enemyType = 3
		} else if rand.Intn(100) < 60 {
			enemyType = 2
		} else {
			enemyType = 1
		}
	case g.Tick > 3000:
		if rand.Intn(100) < 20 {
			enemyType = 3
		} else if rand.Intn(100) < 40 {
			enemyType = 2
		} else {
			enemyType = 1
		}
	case g.Tick > 2000:
		if rand.Intn(100) < 20 {
			enemyType = 2
		} else {
			enemyType = 1
		}
	case g.Tick > 1000:
		enemyType = 1

	default:
		enemyType = 1
	}

	g.newEnemyType(x, y, enemyType)
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
			KillSoundPlayer.Play()
			g.Points++
			g.player.Experience += 3

			// Drop a potion based on enemy's potionSpawnRate
			if rand.Intn(100) < enemy.PotionSpawnRate {
				potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
				if err != nil {
					// handle error
					log.Fatal(err)
				}
				g.potions = append(g.potions, &entities.Potion{
					Sprite: &entities.Sprite{
						Img: potionImg,
						X:   enemy.X,
						Y:   enemy.Y,
					},
					AmtHeal: 2,
					Status:  "DROPPING",
					Count:   0,
				})
			}

			g.killEnemy(enemy)
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

		g.IsTouchingPlayer(enemy.Collider, *g.player)

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
