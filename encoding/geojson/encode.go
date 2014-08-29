package geojson

import (
	"encoding/json"
	"github.com/foobaz/geom"
	"reflect"
)

func ToGeoJSON(t geom.T) (interface{}, error) {
	switch g := t.(type) {
	case geom.Point:
		return &Geometry{
			Type:        "Point",
			Coordinates: g,
		}, nil
	case geom.LineString:
		return Geometry{
			Type:        "LineString",
			Coordinates: g,
		}, nil
	case geom.Polygon:
		return Geometry{
			Type:        "Polygon",
			Coordinates: g,
		}, nil
	case geom.MultiPolygon:
		return Geometry{
			Type:        "MultiPolygon",
			Coordinates: g,
		}, nil
	case geom.Feature:
		serializable, err := ToGeoJSON(g.T)
		if err != nil {
			return nil, err
		}

		geometry, ok := serializable.(Geometry)
		if !ok {
			return nil, &UnsupportedGeometryError{reflect.TypeOf(geometry).String()}
		}

		return Feature{
			Type:       "Feature",
			Geometry:   geometry,
			Properties: g.Properties,
		}, nil
	case geom.FeatureCollection:
		features := make([]Feature, 0, len(g.Features))
		for _, geomFeature := range g.Features {
			serializable, err := ToGeoJSON(geomFeature)
			if err != nil {
				continue
			}

			jsonFeature, ok := serializable.(Feature)
			if !ok {
				continue
			}

			features = append(features, jsonFeature)
		}
		return FeatureCollection{
			Type:       "FeatureCollection",
			Features:   features,
			Properties: g.Properties,
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
