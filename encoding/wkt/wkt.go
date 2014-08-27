package wkt

import (
	"fmt"
	"reflect"
)

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkt: unsupported geometry type: " + e.Type.String()
}

type UnsupportedAxesError struct {
	Axes int
}

func (e UnsupportedAxesError) Error() string {
	return fmt.Sprintf("wkt: unsupported axes %d", e.Axes)
}
