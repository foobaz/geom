package wkb

import (
	"encoding/binary"
	"github.com/foobaz/geom"
	"io"
)

func multiLineStringReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	var numLineStrings uint32
	if err := binary.Read(r, byteOrder, &numLineStrings); err != nil {
		return nil, err
	}
	lineStrings := make([]geom.LineString, numLineStrings)
	for i := range lineStrings {
		g, err := Read(r)
		if err != nil {
			return nil, err
		}

		var ok bool
		lineStrings[i], ok = g.(geom.LineString)
		if !ok {
			return nil, &UnexpectedGeometryError{g}
		}
	}
	return geom.MultiLineString(lineStrings), nil
}

func writeMultiLineString(w io.Writer, byteOrder binary.ByteOrder, axes uint32, multiLineString geom.MultiLineString) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiLineString))); err != nil {
		return err
	}
	for _, lineString := range multiLineString {
		if err := Write(w, byteOrder, axes, lineString); err != nil {
			return err
		}
	}
	return nil
}
