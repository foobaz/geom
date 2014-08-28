package geom

import (
	"math"
	"reflect"
)

func similar(a, b, e float64) bool {
	return math.Abs(a-b) < e
}

func pointSimilar(p1, p2 Point, e float64) bool {
	if len(p1) != len(p2) {
		return false
	}

	for i := range p1 {
		if !similar(p1[i], p2[i], e) {
			return false
		}
	}

	return true
}

func pointsSimilar(p1s, p2s []Point, e float64) bool {
	if len(p1s) != len(p2s) {
		return false
	}

	for i := range p1s {
		if !pointSimilar(p1s[i], p2s[i], e) {
			return false
		}
	}

	return true
}

func pointssSimilar(p1ss, p2ss Polygon, e float64) bool {
	if len(p1ss) != len(p2ss) {
		return false
	}

	for i := range p1ss {
		if !pointsSimilar(p1ss[i], p2ss[i], e) {
			return false
		}
	}

	return true
}

func Similar(t1, t2 T, e float64) bool {
	if reflect.TypeOf(t1) != reflect.TypeOf(t2) {
		return false
	}
	switch t1.(type) {
	case Point:
		return pointSimilar(t1.(Point), t2.(Point), e)
	case LineString:
		return pointsSimilar(t1.(LineString), t2.(LineString), e)
	case Polygon:
		return pointssSimilar(t1.(Polygon), t2.(Polygon), e)
	case MultiPoint:
		return pointsSimilar(t1.(MultiPoint), t2.(MultiPoint), e)
	default:
		return false
	}
}
