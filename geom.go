package geom

const (
	TwoD = iota
	Z
	M
	ZM
)

type T interface {
	Bounds(*Bounds) *Bounds
}

type Geom interface {
	T
}
