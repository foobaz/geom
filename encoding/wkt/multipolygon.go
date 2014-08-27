package wkt

import (
	"github.com/foobaz/geom"
)

func appendMultiPolygonWKT(dst []byte, multiPolygon geom.MultiPolygon, name []byte, dimension int) []byte {
	dst = append(dst, []byte("MULTIPOLYGON")...)
	dst = append(dst, name...)
	dst = append(dst, '(')
	dst = append(dst, '(')
	for i, pg := range multiPolygon{
		dst = appendPointssCoords(dst, pg, dimension)
		if i != len(multiPolygon)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
