package wkb

import (
	"encoding/binary"
	"github.com/foobaz/geom"
	"io"
)

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	points, err := readPoints(r, byteOrder, dimension)
	if err != nil {
		return nil, err
	}
	return geom.LineString(points), nil
}

func writeLineString(w io.Writer, byteOrder binary.ByteOrder, dimension int, lineString geom.LineString) error {
	return writePoints(w, byteOrder, dimension, lineString)
}
