package entities

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	*Sprite
	Collider   Collider
	MaxHealth  uint
	Health     *uint
	Attack     Attack
	Hitbox     *Hitbox
	Status     string
	HitCounter int
	Invencible bool
}

type Attack struct {
	Damage uint
	Img    *ebiten.Image
}
