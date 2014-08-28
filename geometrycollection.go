package geom

type GeometryCollection []Geom

func (geometryCollection GeometryCollection) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geom := range geometryCollection {
		b = geom.Bounds(b)
	}
	return b
}