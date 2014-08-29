package geom

type LineString []Point

func (lineString LineString) Bounds(b Bounds) Bounds {
	return b.ExtendPoints(lineString)
}
