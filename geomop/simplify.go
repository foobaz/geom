package geomop

import (
	"github.com/foobaz/geom"
)

// Simplify performs Ramer-Douglas-Peucker simplification.

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
	p0 := r[0]
	maxDistance := 0.0
	maxIndex := 0
	for i := 1; i < numPoints; i++ {
		dist := d(p0, r[i])
		if dist > maxDistance {
			maxDistance = dist
			maxIndex = i
		}
	}
	var s1, s2 []geom.Point
	s1 = append(s1, r[0:maxIndex+1]...)
	s2 = append(s2, r[maxIndex:]...)
	s2 = append(s2, r[0])

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
