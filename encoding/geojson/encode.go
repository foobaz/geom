package geojson

import (
	"encoding/json"
	"github.com/foobaz/geom"
	"reflect"
)

func ToGeoJSON(t geom.T) (*Geometry, error) {
	switch g := t.(type) {
	case geom.Point:
		return &Geometry{
			Type:        "Point",
			Coordinates: g,
		}, nil
	case geom.LineString:
		return &Geometry{
			Type:        "LineString",
			Coordinates: g,
		}, nil
	case geom.Polygon:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: g,
		}, nil
	case geom.MultiPolygon:
		return &Geometry{
			Type:        "MultiPolygon",
			Coordinates: g,
		}, nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g).String()}
	}
}

func Encode(g geom.T) ([]byte, error) {
	if object, err := ToGeoJSON(g); err == nil {
		return json.Marshal(object)
	} else {
		return nil, err
	}
}
