package geom

type Ring []Point
type Polygon []Ring

func (polygon Polygon) Bounds(b Bounds) Bounds {
	return b.ExtendPointss(polygon)
}
