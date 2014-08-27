package wkb

import (
	"encoding/binary"
	"github.com/foobaz/geom"
	"io"
)

func multiPointReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make([]geom.Point, numPoints)
	for i := range points {
		if g, err := Read(r); err == nil {
			var ok bool
			points[i], ok = g.(geom.Point)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPoint(points), nil
}

func writeMultiPoint(w io.Writer, byteOrder binary.ByteOrder, axes uint32, multiPoint geom.MultiPoint) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPoint))); err != nil {
		return err
	}
	for _, point := range multiPoint {
		if err := Write(w, byteOrder, axes, point); err != nil {
			return err
		}
	}
	return nil
}
