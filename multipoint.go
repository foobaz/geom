package geom

type MultiPoint []Point

func (multiPoint MultiPoint) Bounds(b Bounds) Bounds {
	for _, point := range multiPoint {
		b = point.Bounds(b)
	}

	return b
}
