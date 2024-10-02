package entities

import (
	"math"
)

// Hitbox representa uma forma geométrica arbitrária com vértices definidos
type Hitbox struct {
	X, Y, DirectionX, DirectionY float64
	Width, Height, Distance      float64
	Vertices                     [4][2]float64 // Cada vértice tem suas próprias coordenadas X e Y
}

// NewHitbox cria uma nova hitbox com vértices definidos
func NewHitbox(x1, y1, x2, y2, x3, y3, x4, y4 float64) *Hitbox {
	return &Hitbox{
		Vertices: [4][2]float64{
			{x1, y1},
			{x2, y2},
			{x3, y3},
			{x4, y4},
		},
	}
}

// Overlaps verifica se a hitbox atual está colidindo/sobrepondo com um sprite
func (h *Hitbox) Overlaps(sprite *Sprite) bool {
	// Função auxiliar para projetar um ponto em um eixo
	project := func(vertex [2]float64, axis [2]float64) float64 {
		return (vertex[0]*axis[0] + vertex[1]*axis[1]) / (axis[0]*axis[0] + axis[1]*axis[1])
	}

	// Função auxiliar para verificar se há uma separação em um eixo
	isSeparated := func(axis [2]float64) bool {
		minA, maxA := math.Inf(1), math.Inf(-1)
		minB, maxB := math.Inf(1), math.Inf(-1)

		for _, vertex := range h.Vertices {
			projection := project(vertex, axis)
			if projection < minA {
				minA = projection
			}
			if projection > maxA {
				maxA = projection
			}
		}

		spriteVertices := [4][2]float64{
			{sprite.X, sprite.Y},
			{sprite.X + 16, sprite.Y},
			{sprite.X + 16, sprite.Y + 16},
			{sprite.X, sprite.Y + 16},
		}

		for _, vertex := range spriteVertices {
			projection := project(vertex, axis)
			if projection < minB {
				minB = projection
			}
			if projection > maxB {
				maxB = projection
			}
		}

		return maxA < minB || maxB < minA
	}

	// Verificar todos os eixos normais dos dois polígonos
	for i := 0; i < 4; i++ {
		axis1 := [2]float64{
			h.Vertices[i][1] - h.Vertices[(i+1)%4][1],
			h.Vertices[(i+1)%4][0] - h.Vertices[i][0],
		}
		if isSeparated(axis1) {
			return false
		}

		if i < 2 { // Apenas dois eixos são necessários para o sprite (retângulo alinhado aos eixos)
			spriteVertices := [4][2]float64{
				{sprite.X, sprite.Y},
				{sprite.X + 16, sprite.Y},
				{sprite.X + 16, sprite.Y + 16},
				{sprite.X, sprite.Y + 16},
			}
			axis2 := [2]float64{
				spriteVertices[i][1] - spriteVertices[(i+1)%4][1],
				spriteVertices[(i+1)%4][0] - spriteVertices[i][0],
			}
			if isSeparated(axis2) {
				return false
			}
		}
	}

	return true
}
