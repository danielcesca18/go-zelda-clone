package entities

type Enemy struct {
	*Sprite
	Status        *string
	Health        *uint
	HitCounter    int
	FollowsPlayer bool
	Collider      Collider
	Knockback     Knockback
}

type Knockback struct {
	DirectionX, DirectionY, Velocity float64
}
