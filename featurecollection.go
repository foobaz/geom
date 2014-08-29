package geom

import ()

type FeatureCollection struct {
	Features   []T
	Properties interface{}
}

func (f FeatureCollection) Bounds(b Bounds) Bounds {
	for _, feature := range f.Features {
		b = feature.Bounds(b)
	}

	return b
}

func (f FeatureCollection) AppendFeature(newFeature Feature) FeatureCollection {
	f.Features = append(f.Features, newFeature)
	return f
}

func (f FeatureCollection) AppendGeometry(t T, properties interface{}) FeatureCollection {
	newFeature := NewFeature(t, properties)
	return f.AppendFeature(newFeature)
}
