package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) DrawHUD(screen *ebiten.Image) {
	// Draw HUD
	// Draw health bar
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lv.: %d", g.player.Level), 160, 7)
	vector.DrawFilledRect(screen, 209, 9, 102, 12, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, 210, 10, 100, 10, color.RGBA{220, 30, 30, 255}, false)
	vector.DrawFilledRect(screen, 210, 10, float32(100*(*g.player.Health)/g.player.MaxHealth), 10, color.RGBA{0, 204, 0, 255}, false)

	// Draw experience bar
	vector.DrawFilledRect(screen, 209, 22, 102, 4, color.RGBA{255, 255, 255, 255}, false)
	expBarWidth := float32(100 * g.player.Experience / (g.player.Level * 10))
	vector.DrawFilledRect(screen, 210, 23, expBarWidth, 2, color.RGBA{0, 0, 204, 255}, false)
}
