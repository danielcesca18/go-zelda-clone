package entities

type Potion struct {
	*Sprite
	AmtHeal uint
	Status  string
	Count   int
}
