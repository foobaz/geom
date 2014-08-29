package geom

import (
	"reflect"
	"testing"
)

func TestBounds(t *testing.T) {

	var testCases = []struct {
		g      T
		bounds Bounds
	}{
		{
			Point{1, 2},
			Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			Point{1, 2, 3},
			Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			Point{1, 2, 3, 4},
			Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			LineString{{1, 2}, {3, 4}},
			Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			LineString{{1, 2, 3}, {4, 5, 6}},
			Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			LineString{{1, 2, 3, 4}, {5, 6, 7, 8}},
			Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			Polygon{{{1, 2}, {3, 4}, {5, 6}}},
			Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			MultiPoint{{1, 2}, {3, 4}},
			Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			MultiPoint{{1, 2, 3}, {4, 5, 6}},
			Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			MultiPoint{{1, 2, 3, 4}, {5, 6, 7, 8}},
			Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			Bounds{Point{1, 2}, Point{7, 8}},
		},
		{
			MultiLineString{{{1, 2, 3}, {4, 5, 6}}, {{7, 8, 9}, {10, 11, 12}}},
			Bounds{Point{1, 2}, Point{10, 11}},
		},
		{
			MultiLineString{{{1, 2, 3, 4}, {5, 6, 7, 8}}, {{9, 10, 11, 12}, {13, 14, 15, 16}}},
			Bounds{Point{1, 2}, Point{13, 14}},
		},
	}

	for _, tc := range testCases {
		if got := tc.g.Bounds(NewBounds()); !reflect.DeepEqual(got, tc.bounds) {
			t.Errorf("%#v.Bounds() == %#v, want %#v", tc.g, got, tc.bounds)
		}
	}

}

func TestBoundsEmpty(t *testing.T) {
	if got := NewBounds().Empty(); got != true {
		t.Errorf("NewBounds.Empty() == %#v, want true", got)
	}
}
