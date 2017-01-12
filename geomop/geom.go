// Copyright (c) 2011 Mateusz Czapliński (Go port)
// Copyright (c) 2011 Mahir Iqbal (as3 version)
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// based on http://code.google.com/p/as3polyclip/ (MIT licensed)
// and code by Martínez et al: http://wwwdi.ujaen.es/~fmartin/bool_op.html (public domain)

// Package geomop provides implementation of algorithms for geometry operations.
// For further details, consult the description of Polygon.Construct method.
package geomop

import (
	"github.com/foobaz/geom"
	"math"
	"reflect"
)

// Equals returns true if both p1 and p2 describe the same point within
// a tolerance limit.
func PointEquals(p1, p2 geom.Point) bool {
	//	return (p1[0] == p2[0] && p1[1] == p2[1])
	return (p1[0] == p2[0] && p1[1] == p2[1]) ||
		(math.Abs(p1[0]-p2[0])/math.Abs(p1[0]+p2[0]) < tolerance &&
			math.Abs(p1[1]-p2[1])/math.Abs(p1[1]+p2[1]) < tolerance)
}

func pointSubtract(p1, p2 geom.Point) geom.Point {
	return geom.Point{p1[0] - p2[0], p1[1] - p2[1]}
}

// Length returns distance from p to point (0, 0).
func lengthToOrigin(p geom.Point) float64 {
	return math.Sqrt(p[0]*p[0] + p[1]*p[1])
}

// Used to represent an edge of a polygon.
type segment struct {
	start, end geom.Point
}

// Contour represents a sequence of vertices connected by line segments, forming a closed shape.
type Contour []geom.Point

func (c Contour) segment(index int) segment {
	if index == len(c)-1 {
		return segment{c[len(c)-1], c[0]}
	}
	return segment{c[index], c[index+1]}
	// if out-of-bounds, we expect panic detected by runtime
}

// Checks if a point is inside a contour using the "point in polygon" raycast method.
// This works for all polygons, whether they are clockwise or counter clockwise,
// convex or concave.
// See: http://en.wikipedia.org/wiki/Point_in_polygon#Ray_casting_algorithm
// Returns true if p is inside the polygon defined by contour.
func (c Contour) Contains(p geom.Point) bool {
	intersections := 0
	for i := range c {
		curr := c[i]
		ii := i + 1
		if ii == len(c) {
			ii = 0
		}
		next := c[ii]

		// see if a ray cast to the right crosses this segment
		if rayCrosses(p, curr, next) {
			intersections++
		}
	}
	return intersections%2 != 0
}

func rayCrosses(rayOrigin, start, end geom.Point) bool {
	p := geom.New2Point(rayOrigin[0], rayOrigin[1])
	// ensure the segment is flat or heads upward
	if start[1] > end[1] {
		start, end = end, start
	}

	// nudge the point if it matches a segment component
	// this should not affect correctness for any point
	// that is fully inside the polygon, and by definition
	// the algorithm is unpredictable due to rounding error
	// for points on the polygon boundary.
	for p[1] == start[1] || p[1] == end[1] {
		p[1] = math.Nextafter(p[1], math.Inf(1))
	}
	for p[0] == start[0] || p[0] == end[0] {
		p[0] = math.Nextafter(p[0], math.Inf(1))
	}

	// is an intersection even possible?
	if p[1] < start[1] || p[1] > end[1] {
		return false
	}

	// ok, the y-coords indicate a possible intersection
	// check to see if an intersection is certain or impossible
	if start[0] > end[0] {
		if p[0] > start[0] {
			// point to right of rightmost segment point, crossing impossible
			return false
		}
		if p[0] < end[0] {
			// point to left of leftmost segment point, crossing certain
			return true
		}
	} else {
		if p[0] > end[0] {
			// point to right of rightmost segment point, crossing impossible
			return false
		}
		if p[0] < start[0] {
			// point to left of leftmost segment point, crossing certain
			return true
		}
	}
	// compare slopes to see if the ray crosses
	return (p[1]-start[1])/(p[0]-start[0]) >= (end[1]-start[1])/(end[0]-start[0])
}

// Clone returns a copy of a contour.
func (c Contour) Clone() Contour {
	return append([]geom.Point{}, c...)
}

// NumVertices returns total number of all vertices of all contours of a polygon.
func NumVertices(p geom.Polygon) int {
	num := 0
	for _, c := range p {
		num += len(c)
	}
	return num
}

// Clone returns a duplicate of a polygon.
func Clone(p geom.Polygon) geom.Polygon {
	var r geom.Polygon
	r = make(geom.Polygon, len(p))
	for i, rr := range p {
		r[i] = make(geom.Ring, len(rr))
		for j, pp := range p[i] {
			r[i][j] = make(geom.Point, len(pp))
			copy(r[i][j], pp)
		}
	}
	return r
}

// Op describes an operation which can be performed on two polygons.
type Op int

const (
	UNION Op = iota
	INTERSECTION
	DIFFERENCE
	XOR
)

// Construct computes a 2D polygon, which is a result of performing
// specified Boolean operation on the provided pair of polygons (p <Op> clipping).
// It uses algorithm described by F. Martínez, A. J. Rueda, F. R. Feito
// in "A new algorithm for computing Boolean operations on polygons"
// - see: http://wwwdi.ujaen.es/~fmartin/bool_op.html
// The paper describes the algorithm as performing in time O((n+k) log n),
// where n is number of all edges of all polygons in operation, and
// k is number of intersections of all polygon edges.
// "subject" and "clipping" can both be of type geom.Polygon,
// geom.MultiPolygon, geom.LineString, or geom.MultiLineString.
func Construct(subject, clipping geom.T, operation Op) geom.T {
	// Prepare the input shapes
	var c clipper
	switch clipping.(type) {
	case geom.Polygon, geom.MultiPolygon:
		c.subject = convertToPolygon(subject)
		c.clipping = convertToPolygon(clipping)
		switch subject.(type) {
		case geom.Polygon, geom.MultiPolygon:
			c.outType = outputPolygons
		case geom.LineString, geom.MultiLineString:
			c.outType = outputLines
		}

	case geom.LineString, geom.MultiLineString:
		switch subject.(type) {
		case geom.Polygon, geom.MultiPolygon:
			// swap clipping and subject
			c.subject = convertToPolygon(clipping)
			c.clipping = convertToPolygon(subject)
			c.outType = outputLines
		case geom.LineString, geom.MultiLineString:
			c.subject = convertToPolygon(subject)
			c.clipping = convertToPolygon(clipping)
			c.outType = outputPoints
		}
	}
	// Run the clipper
	return c.compute(operation)
}

// convert input shapes to polygon to make internal processing easier
func convertToPolygon(t geom.T) geom.Polygon {
	var out geom.Polygon
	switch g := t.(type) {
	case geom.Polygon:
		out = g
	case geom.MultiPolygon:
		out = make(geom.Polygon, 0)
		for _, p := range g {
			for _, r := range p {
				out = append(out, r)
			}
		}
	case geom.LineString:
		out = make(geom.Polygon, 1)
		out[0] = geom.Ring(g)
	case geom.MultiLineString:
		out = make(geom.Polygon, len(g))
		for i, ls := range g {
			out[i] = geom.Ring(ls)
		}
	default:
		panic(NewError(g))
	}
	// The clipper doesn't work well if a shape is made up of only two points.
	// To get around this problem, if there are only 2 points, we add a third
	// one a small distance from the second point.
	// However, if there is only 1 point, we just delete the shape.
	for i, r := range out {
		if len(r) == 0 {
			continue
		} else if len(r) == 1 {
			out[i] = make([]geom.Point, 0)
		} else if len(r) == 2 {
			const delta = 0.00001
			newpt := geom.Point{r[1][0] + (r[1][0]-r[0][0])*delta,
				r[1][1] - (r[1][1]-r[0][1])*delta}
			out[i] = append(r, newpt)
		}
	}
	return out
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func NewError(g geom.T) UnsupportedGeometryError {
	return UnsupportedGeometryError{reflect.TypeOf(g)}
}

func (e UnsupportedGeometryError) Error() string {
	return "Unsupported geometry type: " + e.Type.String()
}
