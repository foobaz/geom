package geomop

import (
	"github.com/foobaz/geom"

	"math"
)

// TODO: this should be a parameter and not a magic number
const minArea = 1e-8

// Organize takes a polygon with a number of rings of any winding order,
// and returns
// 1) a slice of polygons each of which have only one positive
// area ring and any holes that lie inside that outer ring.
// 2) a polygon with all the rings that couldn't be processed
func Organize(poly geom.Polygon) ([]geom.Polygon, geom.Polygon) {
	// output
	var result []geom.Polygon
	var leftovers geom.Polygon

	// degenerate case
	if len(poly) == 0 {
		return result, leftovers
	}

	// first, sort into negative and positive rings
	type posRing struct {
		p     geom.Ring
		area  float64
		holes []int
	}

	var neg []geom.Ring
	var pos []posRing

	for _, r := range poly {
		a := area(r)
		// some input polys may be optimized away to nothing.
		// drop any with ~zero area, put the others in the
		// positive and negative buckets
		switch {
		case math.Abs(a) < minArea:
			leftovers = append(leftovers, r)
		case a > 0:
			pos = append(pos, posRing{p: r,
				area: a,
			})
		default:
			neg = append(neg, r)
		}
	}

	// another degenerate case
	if len(pos) == 0 {
		return result, geom.Polygon(neg)
	}

	// ok, must have positive rings at least. match up the holes if any
	for n, h := range neg {
		// going to assume that the smallest positive ring containing
		// a given negative ring is the one it should be paired with.
		found := false
		bestArea := 0.0
		bestIndex := 0
		for nn, pr := range pos {
			// inefficient, fix sometime
			var outer, inner geom.Polygon
			outer = append(outer, pr.p)
			inner = append(inner, h)
			if Within(inner, outer) {
				if !found || pr.area < bestArea {
					bestIndex = nn
					bestArea = pr.area
					found = true
				}
			}
		}
		if found {
			pos[bestIndex].holes = append(pos[bestIndex].holes, n)
		} else {
			// fail. put it in the leftovers
			leftovers = append(leftovers, h)
		}
	}

	// return new polys that contain the negative rings
	for _, pring := range pos {
		var newp geom.Polygon
		newp = append(newp, pring.p)
		for _, index := range pring.holes {
			newp = append(newp, neg[index])
		}
		result = append(result, Clone(newp))
	}
	return result, leftovers
}
