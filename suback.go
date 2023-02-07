package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SubAck struct {
	Buffer      *bytes.Buffer
	Version     byte
	FixedHeader *FixedHeader
	PacketID    uint16
	Properties  *Properties
	Payload     []byte
}

func NewSubAck(fh *FixedHeader, buffer *bytes.Buffer, version byte) *SubAck {
	return &SubAck{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
	}
}

func (s *SubAck) Encode() ([]byte, error) {
	fhBuf, err := EncodingFixedHeaderPacket(s.FixedHeader)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	_, err = s.Buffer.Write(fhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	vhBuf, err := s.encodeVariant()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	_, err = s.Buffer.Write(vhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	_, err = s.Buffer.Write(s.Payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	return s.Buffer.Bytes(), nil
}

func (s *SubAck) Decode() (*SubAck, error) {
	if err := s.decodeVariant(); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	payload, err := ReadByteWithWidth(s.Buffer.Len(), s.Buffer)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	s.Payload = payload
	s.Buffer = nil
	return s, nil
}

func (s *SubAck) decodeVariant() (err error) {
	pidBuf, err := ReadByteWithWidth(2, s.Buffer)
	if err != nil {
		return err
	}
	s.PacketID = binary.BigEndian.Uint16(pidBuf)

	if s.Version == Version5 {
		s.Properties, err = PropertiesDecodeHandler(s.Buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SubAck) encodeVariant() (result []byte, err error) {
	if s.PacketID > 0 {
		result = append(result, EncodingMSBAndLSB(s.PacketID)...)
	}
	if s.Version == Version5 {
		if s.Properties != nil {
			bs, err := EncodingRemainingLength(s.Properties.Length)
			if err != nil {
				return nil, err
			}
			result = append(result, bs...)
			result = append(result, s.Properties.Encode(SUBACKPropType)...)
		}
	}
	return result, nil
}
