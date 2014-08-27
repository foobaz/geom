package wkt

import (
	"github.com/foobaz/geom"
)

func appendMultiLineStringWKT(dst []byte, multiLineString geom.MultiLineString, name []byte, dimension int) []byte {
	dst = append(dst, []byte("MULTILINESTRING")...)
	dst = append(dst, name...)
	dst = append(dst, '(')
	dst = append(dst, '(')
	for i, ls := range multiLineString {
		dst = appendPointsCoords(dst, ls, dimension)
		if i != len(multiLineString)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
