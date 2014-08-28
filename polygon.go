package geom

type Ring []Point
type Polygon []Ring

func (polygon Polygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointss(polygon)
}
