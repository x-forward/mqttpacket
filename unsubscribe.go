package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Unsubscribe struct {
	Buffer      *bytes.Buffer
	Version     byte
	FixedHeader *FixedHeader
	PacketID    uint16
	Topic       []string
	Properties  *Properties
}

func NewUnsubscribe(fh *FixedHeader, buffer *bytes.Buffer, version byte) *Unsubscribe {
	return &Unsubscribe{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
		Topic:       make([]string, 0),
	}
}

func (u *Unsubscribe) Encode() (result []byte, err error) {
	fhBuf, err := EncodingFixedHeaderPacket(u.FixedHeader)
	if err != nil {
		return nil, err
	}

	_, err = u.Buffer.Write(fhBuf)
	if err != nil {
		return nil, err
	}

	vhBuf, err := u.encodeVariant()
	if err != nil {
		return nil, err
	}

	_, err = u.Buffer.Write(vhBuf)

	if err != nil {
		return nil, err
	}

	payload, err := u.encodePayload()
	if err != nil {
		return nil, err
	}

	_, err = u.Buffer.Write(payload)

	if err != nil {
		return nil, EncodePacketErr
	}
	return u.Buffer.Bytes(), nil
}

func (u *Unsubscribe) Decode() (*Unsubscribe, error) {
	if err := u.decodeVariant(); err != nil {
		return nil, err
	}
	if err := u.decodePayload(); err != nil {
		return nil, err
	}
	u.Buffer = nil
	return u, nil
}

func (u *Unsubscribe) decodeVariant() error {
	pidBuf, err := ReadByteWithWidth(2, u.Buffer)
	if err != nil {
		return err
	}
	u.PacketID = binary.BigEndian.Uint16(pidBuf)
	if u.Version == Version5 {
		properties, err := PropertiesDecodeHandler(u.Buffer)
		if err != nil {
			return err
		}
		u.Properties = properties
	}

	return nil
}

func (u *Unsubscribe) decodePayload() error {
	for {
		// 2byte(MSB LSB) invalid packet identifier
		if u.Buffer.Len() < 2 {
			return errors.New("payload invalid length")
		}
		// read topic name
		topicName, err := ReadUTF8String(true, u.Buffer)
		if err != nil {
			return err
		}

		u.Topic = append(u.Topic, string(topicName))
		if u.Buffer.Len() == 0 {
			break
		}
	}
	return nil
}

func (u *Unsubscribe) encodeVariant() (result []byte, err error) {
	if u.PacketID > 0 {
		result = append(result, EncodingMSBAndLSB(u.PacketID)...)
	}
	if u.Version == Version5 {
		bs, err := EncodingRemainingLength(u.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, u.Properties.Encode(UNSUBSCRIBEPropType)...)
	}
	return result, nil
}

func (u *Unsubscribe) encodePayload() (result []byte, err error) {
	for _, topic := range u.Topic {
		// 2byte(msb+ lsb)+ variable length(topic name) + 1 byte(opts)
		result = append(result, EncodingMSBAndLSB(uint16(len(topic)))...)
		result = append(result, []byte(topic)...)
	}
	return result, nil
}
