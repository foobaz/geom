package wkb

import (
	"encoding/binary"
	"github.com/foobaz/geom"
	"io"
)

func geometryCollectionReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geoms := make(geom.GeometryCollection, numGeometries)
	for i := range geoms {
		if g, err := Read(r); err == nil {
			var ok bool
			geoms[i], ok = g.(geom.Geom)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geoms, nil
}

func writeGeometryCollection(w io.Writer, byteOrder binary.ByteOrder, axes uint32, geometryCollection geom.GeometryCollection) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollection))); err != nil {
		return err
	}
	for _, geom := range geometryCollection {
		if err := Write(w, byteOrder, axes, geom); err != nil {
			return err
		}
	}
	return nil
}
