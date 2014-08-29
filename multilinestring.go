package geom

type MultiLineString []LineString

func (multiLineString MultiLineString) Bounds(b Bounds) Bounds {
	for _, lineString := range multiLineString {
		b = lineString.Bounds(b)
	}

	return b
}
