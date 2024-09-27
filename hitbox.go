package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) UpdateHitbox() {
	// Obter a posição do mouse
	mouseX, mouseY := ebiten.CursorPosition()

	// Ajustar a posição do mouse de acordo com a câmera
	mouseX -= int(g.cam.X)
	mouseY -= int(g.cam.Y)

	// Calcular a direção do vetor entre o jogador e o mouse
	playerCenterX := g.player.X + 8 // Centro do jogador (16x16)
	playerCenterY := g.player.Y + 8 // Centro do jogador (16x16)
	directionX := float64(mouseX) - playerCenterX
	directionY := float64(mouseY) - playerCenterY

	// Normalizar o vetor de direção
	magnitude := math.Sqrt(directionX*directionX + directionY*directionY)
	if magnitude != 0 {
		directionX /= magnitude
		directionY /= magnitude
	}

	// Definir uma distância fixa do jogador
	distance := 14.0

	// Calcular a nova posição da hitbox
	hitboxX := playerCenterX + directionX*distance
	hitboxY := playerCenterY + directionY*distance

	// Calcular os vértices da hitbox com base na nova posição e direção
	halfWidth := 14.0  // Metade da largura da hitbox (16x16)
	halfHeight := 14.0 // Metade da altura da hitbox (16x16)

	vertices := [4][2]float64{
		{-halfWidth, -halfHeight},
		{halfWidth, -halfHeight},
		{halfWidth, halfHeight},
		{-halfWidth, halfHeight},
	}

	// Aplicar rotação aos vértices
	sin, cos := math.Sincos(math.Atan2(directionY, directionX))
	for i := range vertices {
		x := vertices[i][0]
		y := vertices[i][1]
		vertices[i][0] = x*cos - y*sin + hitboxX
		vertices[i][1] = x*sin + y*cos + hitboxY
	}

	// Atualizar os vértices da hitbox
	for i := range g.player.Hitbox.Vertices {
		g.player.Hitbox.Vertices[i][0] = vertices[i][0]
		g.player.Hitbox.Vertices[i][1] = vertices[i][1]
	}
}

func (g *Game) DrawHitbox(screen *ebiten.Image) {
	if g.showColliders {
		var path vector.Path

		x1 := float32(g.cam.X + g.player.Hitbox.Vertices[0][0])
		y1 := float32(g.cam.Y + g.player.Hitbox.Vertices[0][1])
		x2 := float32(g.cam.X + g.player.Hitbox.Vertices[1][0])
		y2 := float32(g.cam.Y + g.player.Hitbox.Vertices[1][1])
		x3 := float32(g.cam.X + g.player.Hitbox.Vertices[2][0])
		y3 := float32(g.cam.Y + g.player.Hitbox.Vertices[2][1])
		x4 := float32(g.cam.X + g.player.Hitbox.Vertices[3][0])
		y4 := float32(g.cam.Y + g.player.Hitbox.Vertices[3][1])

		path.MoveTo(x1, y1)
		path.LineTo(x2, y2)
		path.LineTo(x3, y3)
		path.LineTo(x4, y4)
		path.Close()

		var vs []ebiten.Vertex
		var is []uint16
		op := &vector.StrokeOptions{}
		op.Width = 1
		op.LineJoin = vector.LineJoinRound
		vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, op)

		op2 := &ebiten.DrawTrianglesOptions{}
		op2.AntiAlias = false

		op2.FillRule = ebiten.NonZero

		whiteImage := ebiten.NewImage(3, 3)
		whiteImage.Fill(color.White)

		// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// 	whiteImage.Fill(color.RGBA{255, 0, 0, 100})
		// 	// whiteSubImage.Fill(color.RGBA{255, 0, 0, 100})
		// }

		screen.DrawTriangles(vs, is, whiteImage, op2)
	}
}
