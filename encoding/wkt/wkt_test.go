package wkt

import (
	"reflect"
	"testing"

	"github.com/foobaz/geom"
)

func TestWKT(t *testing.T) {
	var testCases = []struct {
		g   geom.T
		wkt []byte
		axes int
	}{
		{
			geom.Point{1, 2},
			[]byte(`POINT(1 2)`),
			geom.TwoD,
		},
		{
			geom.Point{1, 2, 3},
			[]byte(`POINTZ(1 2 3)`),
			geom.Z,
		},
		{
			geom.Point{1, 2, 3},
			[]byte(`POINTM(1 2 3)`),
			geom.M,
		},
		{
			geom.Point{1, 2, 3, 4},
			[]byte(`POINTZM(1 2 3 4)`),
			geom.ZM,
		},
		{
			geom.LineString{{1, 2}, {3, 4}},
			[]byte(`LINESTRING(1 2,3 4)`),
			geom.TwoD,
		},
		{
			geom.LineString{{1, 2, 3}, {4, 5, 6}},
			[]byte(`LINESTRINGZ(1 2 3,4 5 6)`),
			geom.Z,
		},
		{
			geom.LineString{{1, 2, 3}, {4, 5, 6}},
			[]byte(`LINESTRINGM(1 2 3,4 5 6)`),
			geom.M,
		},
		{
			geom.LineString{{1, 2, 3, 4}, {5, 6, 7, 8}},
			[]byte(`LINESTRINGZM(1 2 3 4,5 6 7 8)`),
			geom.ZM,
		},
		{
			geom.Polygon{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}},
			[]byte(`POLYGON((1 2,3 4,5 6,1 2))`),
			geom.TwoD,
		},
		{
			geom.Polygon{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}},
			[]byte(`POLYGONZ((1 2 3,4 5 6,7 8 9,1 2 3))`),
			geom.Z,
		},
		{
			geom.Polygon{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}},
			[]byte(`POLYGONM((1 2 3,4 5 6,7 8 9,1 2 3))`),
			geom.M,
		},
		{
			geom.Polygon{{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {1, 2, 3, 4}}},
			[]byte(`POLYGONZM((1 2 3 4,5 6 7 8,9 10 11 12,1 2 3 4))`),
			geom.ZM,
		},
	}
	for _, tc := range testCases {
		if got, err := Encode(tc.g, tc.axes); err != nil || !reflect.DeepEqual(got, tc.wkt) {
			t.Errorf("Encode(%#v, %d) == %#v, %#v, want %#v, nil", tc.g, tc.axes, string(got), err, string(tc.wkt))
		}
	}
}
