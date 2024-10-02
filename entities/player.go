package entities

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	*Sprite
	Collider           Collider
	MaxHealth          uint
	Health             *uint
	Attack             Attack
	Speed              float64
	Hitbox             *Hitbox
	Status             string
	HitCounter         int
	Invencible         bool
	Experience         uint
	Level              uint
	PowerUps           []PowerUp
	AttackSpeedCounter int
	AttackSpeed        int
}

type Attack struct {
	Damage uint
	Img    *ebiten.Image
}

type PowerUp struct {
	Name string
	Img  *ebiten.Image
}
