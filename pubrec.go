package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PubRec struct {
	Buffer      *bytes.Buffer
	Version     byte
	FixedHeader *FixedHeader
	PacketID    uint16
	Properties  *Properties
	ReasonCode  int
}

func NewPubRec(fh *FixedHeader, buffer *bytes.Buffer, version byte) *PubRec {
	return &PubRec{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
	}
}

func (p *PubRec) Encode() ([]byte, error) {
	fhBuf, err := EncodingFixedHeaderPacket(p.FixedHeader)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	_, err = p.Buffer.Write(fhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	vhBuf, err := p.encodeVariant()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	_, err = p.Buffer.Write(vhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	return p.Buffer.Bytes(), nil
}

func (p *PubRec) Decode() (*PubRec, error) {
	if err := p.decodeVariant(); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	p.Buffer = nil
	return p, nil
}

func (p *PubRec) decodeVariant() (err error) {
	pidBuf, err := ReadByteWithWidth(2, p.Buffer)
	if err != nil {
		return err
	}
	p.PacketID = binary.BigEndian.Uint16(pidBuf)

	if p.Version == Version5 {
		code, err := p.Buffer.ReadByte()
		if err != nil {
			return err
		}
		p.ReasonCode = int(code)
		p.Properties, err = PropertiesDecodeHandler(p.Buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PubRec) encodeVariant() (result []byte, err error) {
	if p.PacketID > 0 {
		result = append(result, EncodingMSBAndLSB(p.PacketID)...)
	}
	if p.Version == Version5 {
		result = append(result, byte(p.ReasonCode))
		if p.Properties != nil {
			bs, err := EncodingRemainingLength(p.Properties.Length)
			if err != nil {
				return nil, err
			}
			result = append(result, bs...)
			result = append(result, p.Properties.Encode(PUBRECPropType)...)
		}
	}
	return result, nil
}
