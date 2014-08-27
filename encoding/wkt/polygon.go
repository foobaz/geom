package wkt

import (
	"github.com/foobaz/geom"
)

func appendPolygonWKT(dst []byte, polygon geom.Polygon, name []byte, dimension int) []byte {
	dst = append(dst, []byte("POLYGON")...)
	dst = append(dst, name...)
	dst = append(dst, '(')
	dst = appendPointssCoords(dst, polygon, dimension)
	dst = append(dst, ')')

	return dst
}
