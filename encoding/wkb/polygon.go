package wkb

import (
	"encoding/binary"
	"io"

	"github.com/foobaz/geom"
)

func polygonReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make(geom.Polygon, numRings)
	for i := uint32(0); i < numRings; i++ {
		if points, err := readPoints(r, byteOrder, dimension); err != nil {
			return nil, err
		} else {
			rings[i] = points
		}
	}
	return rings, nil
}

func writePolygon(w io.Writer, byteOrder binary.ByteOrder, dimension int, polygon geom.Polygon) error {
	return writePointss(w, byteOrder, dimension, polygon)
}
