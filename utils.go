package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"

	"unicode/utf8"
)

func ReadUTF8String(mustUTF8 bool, rd *bytes.Buffer) (b []byte, err error) {
	msbAndLsb, err := ReadByteWithWidth(2, rd)
	if err != nil {
		return nil, errors.New("read utf8 string failed")
	}
	length := int(binary.BigEndian.Uint16(msbAndLsb))
	if length > rd.Len() {
		return nil, errors.New("invalid buffer length")
	}
	payload := rd.Next(length)
	if mustUTF8 {
		if !utf8.Valid(payload) {
			return nil, errors.New("invalid UTF-8 string")
		}
	}
	return payload, nil
}

// EncodingMSBAndLSB packed variable size packet into []byte
func EncodingMSBAndLSB(pack uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, pack)
	return b
}

func ReadByteWithWidth(width int, rd *bytes.Buffer) ([]byte, error) {
	buf := make([]byte, width)
	_, err := rd.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
