package geojson

import (
	"fmt"
)

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type InvalidGeometryError struct{}

func (e InvalidGeometryError) Error() string {
	return "geojson: invalid geometry"
}

type UnsupportedGeometryError struct {
	Type string
}

func (e UnsupportedGeometryError) Error() string {
	return "geojson: unsupported geometry type " + e.Type
}

type InsufficientElementsError struct {
	ElementCount int
}

func (e InsufficientElementsError) Error() string {
	return fmt.Sprintf("geojson: need at least two elements in point, got %d", e.ElementCount)
}

type InsufficientPointsError struct {
	PointCount int
}

func (e InsufficientPointsError) Error() string {
	return fmt.Sprintf("geojson: need more than %d points", e.PointCount)
}
