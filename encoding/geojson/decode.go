package geojson

import (
	"encoding/json"

	"github.com/foobaz/geom"
)

func decodeCoordinates(jsonCoordinates interface{}) []float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}

	elementCount := len(array)
	if elementCount < 2 {
		panic(&InsufficientElementsError{elementCount})
	}

	coordinates := make([]float64, len(array))
	for i, element := range array {
		var ok bool
		if coordinates[i], ok = element.(float64); !ok {
			panic(&InvalidGeometryError{})
		}
	}

	return coordinates
}

func decodeCoordinates2(jsonCoordinates interface{}) [][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}

	coordinates := make([][]float64, len(array))
	for i, element := range array {
		coordinates[i] = decodeCoordinates(element)
	}

	return coordinates
}

func decodeCoordinates3(jsonCoordinates interface{}) [][][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}

	coordinates := make([][][]float64, len(array))
	for i, element := range array {
		coordinates[i] = decodeCoordinates2(element)
	}

	return coordinates
}

func makeLinearRing(coordinates [][]float64) []geom.Point {
	points := make([]geom.Point, len(coordinates))

	for i, c := range coordinates {
		points[i] = c
	}

	return points
}

func makeLinearRings(coordinates [][][]float64) []geom.Ring {
	pointss := make([]geom.Ring, len(coordinates))

	for i, element := range coordinates {
		pointss[i] = makeLinearRing(element)
	}

	return pointss
}

func doFromGeoJSON(g *Geometry) geom.T {
	switch g.Type {
	case "Point":
		coordinates := decodeCoordinates(g.Coordinates)
		return geom.Point(coordinates)
	case "LineString":
		coordinates := decodeCoordinates2(g.Coordinates)
		coordinateCount := len(coordinates)
		if coordinateCount < 2 {
			panic(&InsufficientPointsError{coordinateCount})
		}

		ring := makeLinearRing(coordinates)
		return geom.LineString(ring)
	case "Polygon":
		coordinates := decodeCoordinates3(g.Coordinates)
		for _, ring := range coordinates {
			coordinateCount := len(ring)
			if coordinateCount < 4 {
				panic(&InvalidGeometryError{})
			}
		}

		rings := makeLinearRings(coordinates)
		return geom.Polygon(rings)
	default:
		panic(&UnsupportedGeometryError{g.Type})
	}
}

func FromGeoJSON(geom *Geometry) (g geom.T, err error) {
	defer func() {
		if e := recover(); e != nil {
			g = nil
			err = e.(error)
		}
	}()
	return doFromGeoJSON(geom), nil
}

func Decode(data []byte) (geom.T, error) {
	var geom Geometry
	if err := json.Unmarshal(data, &geom); err == nil {
		return FromGeoJSON(&geom)
	} else {
		return nil, err
	}
}
