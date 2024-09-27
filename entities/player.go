package entities

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	*Sprite
	Collider Collider
	Health   uint
	Attack   Attack
	Hitbox   *Hitbox
	Status   string
}

type Attack struct {
	Damage uint
	Img    *ebiten.Image
}
