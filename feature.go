package geom

import ()

// Create a Feature to serialize GeoJSON with additional arbitrary properties.
// Properties may be any JSON-serializable value. Encoding a Feature to
// another format, like WKT, will only encode the geometry, not the properties.
type Feature struct {
	T
	Properties interface{}
}

func NewFeature(t T, properties interface{}) Feature {
	return Feature{t, properties}
}
