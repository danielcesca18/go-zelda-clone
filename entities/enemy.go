package entities

type Enemy struct {
	*Sprite
	Health        *uint
	FollowsPlayer bool
	Collider      Collider
}
