package geomop

// Simplify performs Ramer-Douglas-Peucker simplification.

import (
	"github.com/foobaz/geom"
	"math"
)

type algorithm int

const (
	TypeRDP algorithm = iota
	TypeGrid
)

// there are two components that might be extremely slow
// the ring simplification algorithm searches for the
// maximally separated nodes, which is O(n^2), and the
// RDP simplification itself is typically O(n log n) but
// is O(n^2) in its worst case.
//
// If a poly with more than maxNodesToSimplify is used
// as input to the simplifier it will return the original.
const maxNodesToSimplify = 10000

// gridded simplification is followed by an attempt to clip off
// degenerate structures such as ears or zero-width peninsulas
// this is the number of nodes that will be searched ahead
// when looking for these degenerate structures
//
// TODO: get rid of this magic number.
// search through the whole poly for ears and in each case
// keep the side with the greater area.
const earSearchNodes = 9

func SimplifyLine(l geom.LineString, eps float64) geom.LineString {
	newL := rdpSimplify(l, eps)
	return newL
}

func Simplify(g geom.T, eps float64, alg algorithm) {
	switch (g).(type) {
	case geom.Polygon:
		p := g.(geom.Polygon)
		for n, r := range p {
			// note if the first and last points are the same
			inputDupe := d(r[0], r[len(r)-1]) < tolerance
			newRing := simplifyRing(r, eps, alg)
			// if our first and last were the same before, make it so now
			if inputDupe {
				outputDupe := d(newRing[0], newRing[len(newRing)-1]) < tolerance
				if !outputDupe {
					newRing = append(newRing, newRing[0])
				}
			}
			p[n] = newRing

		}
	case geom.MultiPolygon:
		for _, p := range g.(geom.MultiPolygon) {
			Simplify(p, eps, alg)
		}
	case geom.GeometryCollection:
		for _, g := range g.(geom.GeometryCollection) {
			Simplify(g, eps, alg)
		}
	}
}

// Simplify a single ring; return the new ring.
func simplifyRing(r geom.Ring, eps float64, alg algorithm) geom.Ring {
	// rule out degenerate case
	if len(r) < 4 {
		return r
	}
	// which type of simplification?
	switch alg {
	case TypeRDP:
		return simplifyRingRDP(r, eps)
	case TypeGrid:
		return simplifyRingGrid(r, eps)
	default:
		panic("unknown poly simplification method invoked")
	}
}
func simplifyRingRDP(r geom.Ring, eps float64) geom.Ring {
	// to prevent this from bogging down on a ridiculous
	// input
	numPoints := len(r)
	if numPoints > maxNodesToSimplify {
		return r
	}
	// find the two points with maximal separation,
	// and use those two points to create the polylines
	// which will be simplified
	//
	// this should result in stable simplification,
	// meaning that the resulting simplified polygon
	// will be the same when processed with a
	// different starting point.
	maxDistance := 0.0
	startIndex := 0
	endIndex := 0
	for i := 0; i < (numPoints - 1); i++ {
		for j := i + 1; j < numPoints; j++ {
			dist := d(r[i], r[j])
			if dist > maxDistance {
				maxDistance = dist
				startIndex = i
				endIndex = j
			}
		}
	}
	var s1, s2 []geom.Point
	// note: the nodes at startIndex and endIndex must be in both polylines
	s1 = append(s1, r[startIndex:endIndex+1]...)
	s2 = append(s2, r[endIndex:]...)
	s2 = append(s2, r[:startIndex+1]...)

	news1 := rdpSimplify(s1, eps)
	news2 := rdpSimplify(s2, eps)

	// seam these back together
	var newRing geom.Ring
	newRing = append(newRing, news1...)
	newRing = append(newRing, news2[1:len(news2)-1]...)

	return newRing
}

func simplifyRingGrid(r geom.Ring, eps float64) geom.Ring {
	var newCoords []geom.Point
	for i := 0; i < len(r); i++ {
		newPoint := roundedPoint(r[i], eps)
		count := len(newCoords)
		// add if it's the first point or different from the most recent point
		if count == 0 || d(newPoint, newCoords[count-1]) > tolerance {
			newCoords = append(newCoords, newPoint)
		}
	}

	// check for collinearity
	// this check is a little more robust due to the removal of repeated points above
	cursor := 0
	polySize := len(newCoords)
	for cursor < polySize {
		// get indices for point before and after
		before := (cursor + polySize - 1) % polySize
		after := (cursor + 1) % polySize
		// maybe optimize this sometime
		triangle := []geom.Point{newCoords[before], newCoords[cursor], newCoords[after]}
		if math.Abs(area(triangle)) < tolerance {
			//if distPointToSegment(newCoords[cursor], newCoords[before], newCoords[after]) < tolerance {
			// remove the point, reduce the size, leave the cursor alone
			newCoords = append(newCoords[:cursor], newCoords[cursor+1:]...)
			polySize--
		} else {
			// keep this point, bump the cursor
			cursor++
		}
	}

	// check for ears/degenerate structures
	cursor = 0
	polySize = len(newCoords)
	for cursor < polySize {
		clipped := false
		for i := 2; !clipped && i < earSearchNodes; i++ {
			after := (cursor + i) % polySize
			if d(newCoords[cursor], newCoords[after]) < tolerance {
				// make sure this won't result in a degenerate poly
				// if we remove it...
				if (polySize - i) > 3 {
					// safe to remove
					clipped = true
					polySize -= i
					// removing from the middle?
					if after > cursor {
						newCoords = append(newCoords[:cursor], newCoords[cursor+i:]...)
					} else {
						// potentially removing from both ends
						newCoords = newCoords[after:cursor]
					}
				}
			}
		}
		// if we didn't clip, the cursor needs to be advanced.
		if !clipped {
			cursor++
		}
	}

	return geom.Ring(newCoords)
}

func rdpSimplify(points []geom.Point, epsilon float64) []geom.Point {
	return points[:rdpCompress(points, epsilon)]
}

func rdpCompress(points []geom.Point, epsilon float64) int {
	end := len(points)

	if end < 3 {
		// return points
		return end
	}

	// Find the point with the maximum distance
	var (
		first = points[0]
		last  = points[end-1]

		dmax  float64
		index int
	)

	for i := 1; i < end-1; i++ {
		d := distPointToSegment(points[i], first, last)
		if d > dmax {
			dmax, index = d, i
		}
	}

	// If max distance is lte to epsilon, return segment containing
	// the first and last points.
	if dmax <= epsilon {
		// return []point{first, last}
		points[1] = last
		return 2
	}

	// Recursive call
	r1 := rdpCompress(points[:index+1], epsilon)
	r2 := rdpCompress(points[index:], epsilon)

	// Build the result list
	// return append(r1[:len(r1)-1], r2...)
	x := r1 - 1
	n := copy(points[x:], points[index:index+r2])

	return x + n
}

func roundedPoint(p geom.Point, eps float64) geom.Point {
	// protect against stupidity
	if eps == 0.0 {
		return p
	}
	var newPoint geom.Point
	for i := 0; i < len(p); i++ {
		coord := math.Floor(p[i]/eps+0.5) * eps
		newPoint = append(newPoint, coord)
	}
	return newPoint
}
