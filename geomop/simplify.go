package geomop

// Simplify performs Ramer-Douglas-Peucker simplification.

import (
	"github.com/foobaz/geom"
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

func Simplify(g geom.T, eps float64) {
	switch (g).(type) {
	case geom.Polygon:
		p := g.(geom.Polygon)
		for n, r := range p {
			p[n] = SimplifyRing(r, eps)
		}
	case geom.MultiPolygon:
		for _, p := range g.(geom.MultiPolygon) {
			Simplify(p, eps)
		}
	case geom.GeometryCollection:
		for _, g := range g.(geom.GeometryCollection) {
			Simplify(g, eps)
		}
	}
}

// Simplify a single ring; return the new ring.
func SimplifyRing(r geom.Ring, eps float64) geom.Ring {
	// handle the ring case, break it down into two lines.
	// rule out degenerate case
	numPoints := len(r)
	if numPoints < 4 {
		return r
	}
	// to prevent this from bogging down on a ridiculous
	// input
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
