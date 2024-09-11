package HEXX

import (
	"cmsApp/pkg/utils/stringx"
	"encoding/hex"
	"errors"
)

func EncodeHex(src string) (string, error) {
	byteSrc := stringx.String2bytes(src)
	if len(byteSrc) == 0 {
		return "", errors.New("encode hex error")
	}
	buf := make([]byte, hex.EncodedLen(len(byteSrc)))
	hex.Encode(buf, byteSrc)
	return stringx.Bytes2string(buf), nil
}

func DecodeHex(src string) (string, error) {
	byteSrc := stringx.String2bytes(src)
	if len(byteSrc) == 0 {
		return "", errors.New("decode hex error")
	}
	buf := make([]byte, hex.DecodedLen(len(byteSrc)))
	n, err := hex.Decode(buf, byteSrc)
	if err != nil {
		return "", errors.New("decode hex error")
	}
	if n > 0 {
		return stringx.Bytes2string(buf), nil
	}
	return "", nil
}
