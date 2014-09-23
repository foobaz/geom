package geom

const (
	X = iota
	Y
)

// The first two components of Point are X and Y. You may use as many additional components as you like for Z, M, etc. They are ignored, but preserved for GeoJSON encoding.
type Point []float64

func New2Point(x, y float64) Point {
	return Point{x, y}
}

func NewNPoint(s ...float64) Point {
	return Point(s)
}

func (point Point) Bounds(b Bounds) Bounds {
	return b.ExtendPoint(point)
}

func (point Point) Equal(other Point) bool {
	n := len(point)
	if n != len(other) {
		return false
	}

	for i := 0; i < n; i++ {
		if point[i] != other[i] {
			return false
		}
	}

	return true
}
