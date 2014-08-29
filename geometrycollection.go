package geom

type GeometryCollection []T

func (geometryCollection GeometryCollection) Bounds(b Bounds) Bounds {
	for _, t := range geometryCollection {
		b = t.Bounds(b)
	}

	return b
}
