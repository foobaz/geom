package hex

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/foobaz/geom"
	"github.com/foobaz/geom/encoding/wkb"
)

// axes must be geom.TwoD, geom.Z, geom.M, or geom.ZM
func Encode(g geom.T, byteOrder binary.ByteOrder, axes uint32) (string, error) {
	wkb, err := wkb.Encode(g, byteOrder, axes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(wkb), nil
}

func Decode(s string) (geom.T, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return wkb.Decode(data)
}
