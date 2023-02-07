package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Publish struct {
	Buffer      *bytes.Buffer
	Version     byte
	FixedHeader *FixedHeader
	Dup         bool
	Qos         uint8
	Retain      bool
	TopicName   []byte
	PacketID    uint16
	Payload     []byte
	Properties  *Properties
}

func NewPublish(fh *FixedHeader, buffer *bytes.Buffer, version byte) *Publish {
	return &Publish{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
	}
}

func (p *Publish) Encode() (result []byte, err error) {
	fhBuf, err := EncodingFixedHeaderPacket(p.FixedHeader)
	if err != nil {
		return nil, err
	}

	_, err = p.Buffer.Write(fhBuf)
	if err != nil {
		return nil, err
	}

	vhBuf, err := p.encodeVariant()
	if err != nil {
		return nil, err
	}

	_, err = p.Buffer.Write(vhBuf)

	if err != nil {
		return nil, err
	}

	_, err = p.Buffer.Write(p.Payload)

	if err != nil {
		return nil, EncodePacketErr
	}
	return p.Buffer.Bytes(), nil
}

func (p *Publish) Decode() (*Publish, error) {
	p.decodeFlag()
	if err := p.decodeVariant(); err != nil {
		return nil, err
	}

	if err := p.decodePayload(); err != nil {
		return nil, err
	}
	p.Buffer = nil
	return p, nil
}

// FixedHeader: [1 byte(packetType + Flag), 1~4byte(RemainingLength)]
// Flag: [1 byte(Dup + QOS-H,  QOS-L, Retain)]
func (p *Publish) decodeFlag() {
	p.Dup = (1 & p.FixedHeader.Flag >> 3) > 0
	p.Qos = (p.FixedHeader.Flag >> 1) & 3
	if p.FixedHeader.Flag&1 == 1 {
		p.Retain = true
	}
}

func (p *Publish) decodeVariant() error {
	topicName, err := ReadUTF8String(true, p.Buffer)
	if err != nil {
		return err
	}
	p.TopicName = topicName

	if p.Qos > 0 {
		pidBuf, err := ReadByteWithWidth(2, p.Buffer)
		if err != nil {
			return err
		}
		p.PacketID = binary.BigEndian.Uint16(pidBuf)
	}

	if p.Version == Version5 {
		properties, err := PropertiesDecodeHandler(p.Buffer)
		if err != nil {
			return err
		}
		p.Properties = properties
	}

	return nil
}

func (p *Publish) decodePayload() error {
	if p.Buffer.Len() <= 0 {
		return errors.New("malformed publish payload packet")
	}
	p.Payload = p.Buffer.Next(p.Buffer.Len())
	return nil
}

func (p *Publish) encodeVariant() (result []byte, err error) {
	result = append(result, EncodingMSBAndLSB(uint16(len(p.TopicName)))...)
	result = append(result, p.TopicName...)

	if p.Qos > 0 {
		result = append(result, EncodingMSBAndLSB(p.PacketID)...)
	}
	if p.Version == Version5 {
		bs, err := EncodingRemainingLength(p.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, p.Properties.Encode(PUBLISHPropType)...)
	}
	return result, nil
}
