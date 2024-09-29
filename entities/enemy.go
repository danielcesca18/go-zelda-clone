package entities

type Enemy struct {
	*Sprite
	Status          *string
	Health          *uint
	HitCounter      int
	FollowsPlayer   bool
	Collider        Collider
	Knockback       Knockback
	PotionSpawnRate int
}

type Knockback struct {
	DirectionX, DirectionY, Velocity float64
}
