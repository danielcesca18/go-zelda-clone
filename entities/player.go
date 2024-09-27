package entities

type Player struct {
	*Sprite
	Collider Collider
	Health   uint
	Damage   uint
	Hitbox   *Hitbox
}
