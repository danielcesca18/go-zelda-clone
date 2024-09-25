package entities

type FloatRect struct {
	MinX, MinY, MaxX, MaxY float64
}

// Função para verificar se dois retângulos de float se sobrepõem
func (r FloatRect) Overlaps(other FloatRect) bool {
	return r.MinX < other.MaxX && r.MaxX > other.MinX &&
		r.MinY < other.MaxY && r.MaxY > other.MinY
}
