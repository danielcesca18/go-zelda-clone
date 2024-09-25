package entities

type Player struct {
	*Sprite
	Collider Collider
	Health   uint
}
