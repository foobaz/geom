package geom

type GeometryCollection struct {
	Geoms []Geom
}

func (geometryCollection GeometryCollection) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geom := range geometryCollection.Geoms {
		b = geom.Bounds(b)
	}
	return b
}
