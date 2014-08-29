package geom

type MultiPolygon []Polygon

func (multiPolygon MultiPolygon) Bounds(b Bounds) Bounds {
	for _, polygon := range multiPolygon {
		b = polygon.Bounds(b)
	}

	return b
}
