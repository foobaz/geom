package wkt

import (
	"github.com/foobaz/geom"
)

func appendLineStringWKT(dst []byte, lineString geom.LineString, name []byte, dimension int) []byte {
	dst = append(dst, []byte("LINESTRING")...)
	dst = append(dst, name...)
	dst = append(dst, '(')
	dst = appendPointsCoords(dst, lineString, dimension)
	dst = append(dst, ')')
	return dst
}
