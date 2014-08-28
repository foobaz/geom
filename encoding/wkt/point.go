package wkt

import (
	"math"
	"strconv"

	"github.com/foobaz/geom"
)

func appendPointCoords(dst []byte, point geom.Point, dimension int) []byte {
	nan := math.NaN()
	elementCount := len(point)

	for i := 0; i < dimension; i++ {
		if i != 0 {
			dst = append(dst, ' ')
		}

		if i < elementCount {
			dst = strconv.AppendFloat(dst, point[i], 'g', -1, 64)
		} else {
			dst = strconv.AppendFloat(dst, nan, 'g', -1, 64)
		}
	}

	return dst
}

func appendPointsCoords(dst []byte, points []geom.Point, dimension int) []byte {
	for i, point := range points {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPointCoords(dst, point, dimension)
	}
	return dst
}

func appendPointssCoords(dst []byte, pointss geom.Polygon, dimension int) []byte {
	for i, points := range pointss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPointsCoords(dst, points, dimension)
		dst = append(dst, ')')
	}
	return dst
}

func appendPointWKT(dst []byte, point geom.Point, name []byte, dimension int) []byte {
	dst = append(dst, []byte("POINT")...)
	dst = append(dst, name...)
	dst = append(dst, '(')
	dst = appendPointCoords(dst, point, dimension)
	dst = append(dst, ')')

	return dst
}
