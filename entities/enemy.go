package entities

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	Collider      Collider
}
