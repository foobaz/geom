package wkt

import (
	"github.com/foobaz/geom"
	"reflect"
)

// axes must be geom.TwoD, geom.Z, geom.M, or geom.ZM
func Encode(t geom.T, axes int) ([]byte, error) {
	name := []byte{}
	dimension := 0
	switch axes {
	case geom.TwoD:
		name = []byte("")
		dimension = 2
	case geom.Z:
		name = []byte("Z")
		dimension = 3
	case geom.M:
		name = []byte("M")
		dimension = 3
	case geom.ZM:
		name = []byte("ZM")
		dimension = 4
	default:
		return nil, UnsupportedAxesError{axes}
	}

	switch g := t.(type) {
	case geom.Point:
		return appendPointWKT(nil, g, name, dimension), nil
	case geom.LineString:
		return appendLineStringWKT(nil, g, name, dimension), nil
	case geom.MultiLineString:
		return appendMultiLineStringWKT(nil, g, name, dimension), nil
	case geom.Polygon:
		return appendPolygonWKT(nil, g, name, dimension), nil
	case geom.MultiPolygon:
		return appendMultiPolygonWKT(nil, g, name, dimension), nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}
