package main

import (
	"fmt"
	"go-game/entities"
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	PowerUps        = []entities.PowerUp{}
	PowerUpsIndexes = []int{}
)

func (g *Game) DrawLevelUp(screen *ebiten.Image) {
	// Draw clickable rectangles
	screenWidth, screenHeight := g.Layout(0, 0)
	rectWidth, rectHeight := 80, 80
	spacing := 10
	totalWidth := 3*rectWidth + 2*spacing
	startX := (screenWidth - totalWidth) / 2
	y := (screenHeight - rectHeight) / 2

	rects := []image.Rectangle{
		image.Rect(startX+2*(rectWidth+spacing), y, startX+3*rectWidth+2*spacing, y+rectHeight),
		image.Rect(startX+rectWidth+spacing, y, startX+2*rectWidth+spacing, y+rectHeight),
		image.Rect(startX, y, startX+rectWidth, y+rectHeight),
	}

	// pu1, pu2, pu3 := rand.Intn(3), rand.Intn(3), rand.Intn(3)

	if len(PowerUpsIndexes) == 0 {
		for _, rect := range rects {

			var puIndex int
			for {
				puIndex = rand.Intn(len(PowerUps))
				duplicate := false
				for _, index := range PowerUpsIndexes {
					if index == puIndex {
						duplicate = true
						break
					}
				}
				if !duplicate {
					break
				}
			}
			PowerUpsIndexes = append(PowerUpsIndexes, puIndex)

			vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), color.RGBA{255, 0, 0, 255}, false)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(rect.Dx())/float64(PowerUps[puIndex].Img.Bounds().Dx()), float64(rect.Dy())/float64(PowerUps[puIndex].Img.Bounds().Dy()))
			op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
			screen.DrawImage(PowerUps[puIndex].Img, op)

			// Draw text below each image
			text := PowerUps[puIndex].Description
			textX := rect.Min.X + (rect.Dx()-len(text)*7)/2 // Center the text
			textY := rect.Max.Y + 10                        // Position below the rectangle
			ebitenutil.DebugPrintAt(screen, text, textX, textY)
		}
	}

	for i, rect := range rects {

		vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), color.RGBA{255, 0, 0, 255}, false)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(rect.Dx())/float64(PowerUps[PowerUpsIndexes[i]].Img.Bounds().Dx()), float64(rect.Dy())/float64(PowerUps[PowerUpsIndexes[i]].Img.Bounds().Dy()))
		op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		screen.DrawImage(PowerUps[PowerUpsIndexes[i]].Img, op)

		// Draw text below each image
		text := PowerUps[PowerUpsIndexes[i]].Description
		textX := rect.Min.X + (rect.Dx()-len(text)*7)/2 // Center the text
		textY := rect.Max.Y + 10                        // Position below the rectangle
		ebitenutil.DebugPrintAt(screen, text, textX, textY)
	}

	// Handle clicks on the rectangles
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i, rect := range rects {
			if rect.Min.X <= x && x <= rect.Max.X && rect.Min.Y <= y && y <= rect.Max.Y {

				if g.player.PowerUps == nil {
					g.player.PowerUps = make(map[string]int)
				}
				if g.player.PowerUps[PowerUps[PowerUpsIndexes[i]].Name] != 0 {
					g.player.PowerUps[PowerUps[PowerUpsIndexes[i]].Name]++
				} else {
					g.player.PowerUps[PowerUps[PowerUpsIndexes[i]].Name] = 1
				}

				for j := 0; j < len(PowerUps); j++ {
					if PowerUps[PowerUpsIndexes[i]].Name == "health" {
						g.player.MaxHealth += 2
						*g.player.Health += 2
						fmt.Println("health")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "attack" {
						g.player.Attack.Damage += 1
						fmt.Println("attack")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "speed" {
						g.player.AttackSpeed -= 3
						g.player.Speed += 0.1
						fmt.Println("speed")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "death" {
						g.killEnemies = true
						fmt.Println("death")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "defense" {
						fmt.Println("defense")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "hitbox" {
						g.player.Hitbox.Width += 2
						g.player.Hitbox.Distance += 1
						fmt.Println("hitbox")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "revive" {
						g.player.Revives++
						fmt.Println("revive")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "thornmail" {
						fmt.Println("thornmail")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "vampirism" {
						g.HitCounter = 0
						g.player.Vampirism += 1
						fmt.Println("vampirism")

					} else if PowerUps[PowerUpsIndexes[i]].Name == "punch" {
						g.player.Punch += 0.2
						fmt.Println("punch")

					}

					PowerUpsIndexes = []int{}

					g.GameState = "RUNNING"
					LevelUpSoundPlayer.Rewind()
					LevelUpSoundPlayer.Play()

					break
				}
			}
		}
	}
}
