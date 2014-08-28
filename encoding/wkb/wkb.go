package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/foobaz/geom"
)

const (
	wkbXDR = iota
	wkbNDR
)

var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

const (
	wkbPoint              = 1
	wkbLineString         = 2
	wkbPolygon            = 3
	wkbMultiPoint         = 4
	wkbMultiLineString    = 5
	wkbMultiPolygon       = 6
	wkbGeometryCollection = 7
	wkbPolyhedralSurface  = 15
	wkbTIN                = 16
	wkbTriangle           = 17
)

type UnexpectedGeometryError struct {
	Geom geom.T
}

func (e UnexpectedGeometryError) Error() string {
	return fmt.Sprintf("wkb: unexpected geometry %v", e.Geom)
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkb: unsupported type: " + e.Type.String()
}

type UnsupportedAxesError struct {
	Axes uint32
}

func (e UnsupportedAxesError) Error() string {
	return fmt.Sprintf("wkb: unsupported axes %d", e.Axes)
}

type wkbReader func(io.Reader, binary.ByteOrder, int) (geom.T, error)

var wkbReaders map[uint32]wkbReader

func init() {
	wkbReaders = make(map[uint32]wkbReader)
	wkbReaders[wkbPoint] = pointReader
	wkbReaders[wkbLineString] = lineStringReader
	wkbReaders[wkbPolygon] = polygonReader
	wkbReaders[wkbMultiPoint] = multiPointReader
	wkbReaders[wkbMultiLineString] = multiLineStringReader
	wkbReaders[wkbMultiPolygon] = multiPolygonReader
	wkbReaders[wkbGeometryCollection] = geometryCollectionReader
}

func Read(r io.Reader) (geom.T, error) {

	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return nil, err
	}
	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case wkbXDR:
		byteOrder = binary.BigEndian
	case wkbNDR:
		byteOrder = binary.LittleEndian
	default:
		return nil, fmt.Errorf("invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return nil, err
	}

	axes := wkbGeometryType / 1000
	dimension := dimensionsInAxes(axes)
	if dimension == 0 {
		return nil, UnsupportedAxesError{axes}
	}

	baseType := wkbGeometryType - (axes * 1000)
	reader, ok := wkbReaders[baseType]
	if !ok {
		return nil, fmt.Errorf("unsupported geometry type %d", wkbGeometryType)
	}

	return reader(r, byteOrder, dimension)
}

func Decode(buf []byte) (geom.T, error) {
	return Read(bytes.NewBuffer(buf))
}

func writeMany(w io.Writer, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		if err := binary.Write(w, byteOrder, datum); err != nil {
			return err
		}
	}
	return nil
}

func Write(w io.Writer, byteOrder binary.ByteOrder, axes uint32, g geom.T) error {
	var wkbByteOrder uint8
	switch byteOrder {
	case XDR:
		wkbByteOrder = wkbXDR
	case NDR:
		wkbByteOrder = wkbNDR
	default:
		return fmt.Errorf("unsupported byte order %v", byteOrder)
	}
	if err := binary.Write(w, byteOrder, wkbByteOrder); err != nil {
		return err
	}

	var wkbGeometryType uint32
	switch g.(type) {
	case geom.Point:
		wkbGeometryType = wkbPoint
	case geom.LineString:
		wkbGeometryType = wkbLineString
	case geom.Polygon:
		wkbGeometryType = wkbPolygon
	case geom.MultiPoint:
		wkbGeometryType = wkbMultiPoint
	case geom.MultiLineString:
		wkbGeometryType = wkbMultiLineString
	case geom.MultiPolygon:
		wkbGeometryType = wkbMultiPolygon
	case geom.GeometryCollection:
		wkbGeometryType = wkbGeometryCollection
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
	wkbGeometryType += (axes * 1000)
	if err := binary.Write(w, byteOrder, wkbGeometryType); err != nil {
		return err
	}

	dimension := dimensionsInAxes(axes)
	if dimension == 0 {
		return UnsupportedAxesError{axes}
	}

	switch g.(type) {
	case geom.Point:
		return writePoint(w, byteOrder, dimension, g.(geom.Point))
	case geom.LineString:
		return writeLineString(w, byteOrder, dimension, g.(geom.LineString))
	case geom.Polygon:
		return writePolygon(w, byteOrder, dimension, g.(geom.Polygon))
	case geom.MultiPoint:
		return writeMultiPoint(w, byteOrder, axes, g.(geom.MultiPoint))
	case geom.MultiLineString:
		return writeMultiLineString(w, byteOrder, axes, g.(geom.MultiLineString))
	case geom.MultiPolygon:
		return writeMultiPolygon(w, byteOrder, axes, g.(geom.MultiPolygon))
	case geom.GeometryCollection:
		return writeGeometryCollection(w, byteOrder, axes, g.(geom.GeometryCollection))
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}

func Encode(g geom.T, byteOrder binary.ByteOrder, axes uint32) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, axes, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func dimensionsInAxes(axes uint32) int {
	dimension := 0
	switch axes {
	case geom.TwoD:
		dimension = 2
	case geom.Z, geom.M:
		dimension = 3
	case geom.ZM:
		dimension = 4
	}

	return dimension
}
