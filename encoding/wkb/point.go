package wkb

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/foobaz/geom"
)

func pointReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	point := make(geom.Point, dimension)
	if err := binary.Read(r, byteOrder, point); err != nil {
		return nil, err
	}
	return point, nil
}

func readPoints(r io.Reader, byteOrder binary.ByteOrder, dimension int) ([]geom.Point, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}

	components := make([]float64, int(numPoints) * dimension)
	if err := binary.Read(r, byteOrder, components); err != nil {
		return nil, err
	}

	points := make([]geom.Point, numPoints)
	for i := range points {
		j := i + 1
		points[i] = geom.Point(components[i*dimension:j*dimension])
	}

	return points, nil
}

func writePoint(w io.Writer, byteOrder binary.ByteOrder, dimension int, point geom.Point) error {
	clipped := make(geom.Point, dimension)
	nan := math.NaN()
	for i := range clipped {
		clipped[i] = nan
	}
	copy(clipped, point)
	return binary.Write(w, byteOrder, clipped)
}

func writePoints(w io.Writer, byteOrder binary.ByteOrder, dimension int, points []geom.Point) error {
	pointCount := len(points)
	countErr := binary.Write(w, byteOrder, uint32(pointCount))
	if countErr != nil {
		return countErr
	}

	nans := make(geom.Point, dimension)
	nan := math.NaN()
	for i := range nans {
		nans[i] = nan
	}
	clipped := make(geom.Point, dimension)
	for _, p := range points {
		copy(clipped, nans)
		copy(clipped, p)
		pointErr := binary.Write(w, byteOrder, clipped)
		if pointErr != nil {
			return pointErr
		}
	}

	return nil
}

func writePointss(w io.Writer, byteOrder binary.ByteOrder, dimension int, pointss [][]geom.Point) error {
	countErr := binary.Write(w, byteOrder, uint32(len(pointss)))
	if countErr != nil {
		return countErr
	}

	for _, points := range pointss {
		pointsErr := writePoints(w, byteOrder, dimension, points)
		if pointsErr != nil {
			return pointsErr
		}
	}
	return nil

}
