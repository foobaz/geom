package geojson

import (
	"encoding/json"
	"reflect"

	"github.com/foobaz/geom"
)

func ToGeoJSON(t geom.T) (interface{}, error) {
	switch g := t.(type) {
	case geom.Point:
		return Geometry{
			Type:        "Point",
			Coordinates: g,
		}, nil
	case geom.LineString:
		err := validateLineString(g)
		if err != nil {
			return nil, err
		}

		return Geometry{
			Type:        "LineString",
			Coordinates: g,
		}, nil
	case geom.MultiLineString:
		err := validateMultiLineString(g)
		if err != nil {
			return nil, err
		}

		return Geometry{
			Type:        "MultiLineString",
			Coordinates: g,
		}, nil
	case geom.Polygon:
		err := validatePolygon(g)
		if err != nil {
			return nil, err
		}

		return Geometry{
			Type:        "Polygon",
			Coordinates: g,
		}, nil
	case geom.MultiPolygon:
		err := validateMultiPolygon(g)
		if err != nil {
			return nil, err
		}

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

func validateLineString(l geom.LineString) error {
	pointCount := len(l)
	if pointCount < 2 {
		return InsufficientPointsError{pointCount}
	}

	return nil
}

func validateMultiLineString(m geom.MultiLineString) error {
	for _, l := range m {
		err := validateLineString(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePolygon(p geom.Polygon) error {
	if len(p) == 0 {
		return InsufficientPointsError{0}
	}

	for _, r := range p {
		pointCount := len(r)
		if pointCount < 4 {
			return InsufficientPointsError{pointCount}
		}
	}

	return nil
}

func validateMultiPolygon(m geom.MultiPolygon) error {
	for _, p := range m {
		err := validatePolygon(p)
		if err != nil {
			return err
		}
	}

	return nil
}
