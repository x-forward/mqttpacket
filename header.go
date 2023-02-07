package packet

import (
	"bytes"
	"errors"
	"io"
)

const (
	FixedHeaderReservedFlag = iota
	FixedHeaderSubscribeFlag
	FixedHeaderUnsubscribeFlag
)

type FixedHeader struct {
	Type            byte
	RemainingLength int
	Flag            byte
}

func EncodingFixedHeaderPacket(fh *FixedHeader) (result []byte, err error) {
	firstByte := fh.Type<<4 | fh.Flag
	result = append(result, firstByte)
	rl, err := EncodingRemainingLength(fh.RemainingLength)
	if err != nil {
		return nil, err
	}
	result = append(result, rl...)
	return result, nil
}

func DecodingFixedHeaderPacket(rd *bytes.Buffer) (*FixedHeader, error) {
	fp, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}
	rl, err := DecodingRemainingLength(rd)
	if err != nil {
		return nil, err
	}
	return &FixedHeader{
		Type:            fp >> 4,
		RemainingLength: rl,
		Flag:            fp & 15,
	}, nil
}

func DecodingRemainingLength(rd *bytes.Buffer) (int, error) {
	var vbi uint32
	var multiplier uint32
	for {
		digit, err := rd.ReadByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
		vbi |= uint32(digit&127) << multiplier
		if vbi > 268435455 {
			return 0, err
		}
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(vbi), nil
}

func EncodingRemainingLength(length int) ([]byte, error) {
	var result []byte
	if length < 128 {
		result = make([]byte, 1)
	} else if length < 16384 {
		result = make([]byte, 2)
	} else if length < 2097152 {
		result = make([]byte, 3)
	} else if length < 268435456 {
		result = make([]byte, 4)
	} else {
		return nil, errors.New("invalid remaining length")
	}
	var i int
	for {
		encodedByte := length % 128
		length = length / 128
		if length > 0 {
			encodedByte = encodedByte | 128
		}
		result[i] = byte(encodedByte)
		i++
		if length <= 0 {
			break
		}
	}
	return result, nil
}
