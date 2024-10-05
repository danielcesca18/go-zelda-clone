package main

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) DrawHUD(screen *ebiten.Image) {

	// draw fps counter
	msg := fmt.Sprintf(
		"TPS: %0.2f\nEnemies: %d",
		ebiten.ActualTPS(),
		len(g.enemies),
	)
	ebitenutil.DebugPrintAt(screen, msg, 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Horde: %d", g.Horde), 0, 225)

	ebitenutil.DebugPrintAt(screen, "[ESC] Pause/PowerUps", 195, 225)

	// Draw HUD
	// Draw health bar
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lv.: %d", g.player.Level), 160, 7)
	if g.player.Invencible {
		vector.DrawFilledRect(screen, 207, 7, 106, 16, color.RGBA{220, 30, 30, 255}, false)
	} else {
		vector.DrawFilledRect(screen, 209, 9, 102, 12, color.RGBA{255, 255, 255, 255}, false)
	}
	vector.DrawFilledRect(screen, 210, 10, 100, 10, color.RGBA{220, 30, 30, 255}, false)
	vector.DrawFilledRect(screen, 210, 10, float32(100*(*g.player.Health)/g.player.MaxHealth), 10, color.RGBA{0, 204, 0, 255}, false)

	// Draw experience bar
	vector.DrawFilledRect(screen, 209, 22, 102, 4, color.RGBA{255, 255, 255, 255}, false)
	expBarWidth := float32(100 * g.player.Experience / (g.player.Level * 10))
	vector.DrawFilledRect(screen, 210, 23, expBarWidth, 2, color.RGBA{0, 0, 204, 255}, false)

	msg = fmt.Sprintf("Score: %d", g.Points)
	ebitenutil.DebugPrintAt(screen, msg, 210, 26)

	if g.GameState == "PAUSED" {
		ebitenutil.DebugPrintAt(screen, "PAUSED", 140, 110)

		// Collect power-up names and sort them
		names := make([]string, 0, len(g.player.PowerUps))
		for name := range g.player.PowerUps {
			names = append(names, name)
		}
		sort.Strings(names)

		// Display sorted power-ups
		for i, name := range names {
			count := g.player.PowerUps[name]
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s x%d", name, count), 210, 40+(i*14))
		}
	}
}
