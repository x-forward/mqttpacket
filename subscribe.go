package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Subscribe struct {
	Buffer      *bytes.Buffer
	FixedHeader *FixedHeader
	Version     byte
	PacketID    uint16
	Properties  *Properties
	// payload
	Topic []Topic
}

type Topic struct {
	Opt  *TopicOpt
	Name []byte
}

type TopicOpt struct {
	Qos               byte
	RetainHandling    byte
	NoLocal           bool
	RetainAsPublished bool
}

func NewSubscribe(fh *FixedHeader, buffer *bytes.Buffer, version byte) *Subscribe {
	return &Subscribe{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
		Topic:       make([]Topic, 0),
	}
}

func (s *Subscribe) Encode() (result []byte, err error) {
	fhBuf, err := EncodingFixedHeaderPacket(s.FixedHeader)
	if err != nil {
		return nil, err
	}

	_, err = s.Buffer.Write(fhBuf)
	if err != nil {
		return nil, err
	}

	vhBuf, err := s.encodeVariant()
	if err != nil {
		return nil, err
	}

	_, err = s.Buffer.Write(vhBuf)

	if err != nil {
		return nil, err
	}

	payload, err := s.encodePayload()
	if err != nil {
		return nil, err
	}

	_, err = s.Buffer.Write(payload)

	if err != nil {
		return nil, EncodePacketErr
	}
	return s.Buffer.Bytes(), nil
}

func (s *Subscribe) Decode() (*Subscribe, error) {
	if err := s.decodeVariant(); err != nil {
		return nil, err
	}
	if err := s.decodePayload(); err != nil {
		return nil, err
	}
	s.Buffer = nil
	return s, nil
}

func (s *Subscribe) decodeVariant() error {
	pidBuf, err := ReadByteWithWidth(2, s.Buffer)
	if err != nil {
		return err
	}
	s.PacketID = binary.BigEndian.Uint16(pidBuf)

	if s.Version == Version5 {
		properties, err := PropertiesDecodeHandler(s.Buffer)
		if err != nil {
			return err
		}
		s.Properties = properties
	}

	return nil
}

func (s *Subscribe) decodePayload() error {
	for {
		// 2byte(MSB LSB) invalid packet identifier
		if s.Buffer.Len() < 2 {
			return errors.New("payload invalid length")
		}
		// read topic name
		topicName, err := ReadUTF8String(true, s.Buffer)
		if err != nil {
			return errors.New("payload invalid length")
		}
		// read topic options
		topic := Topic{Name: topicName}
		opts, err := s.Buffer.ReadByte()
		if err != nil {
			return errors.New("invalid topic opt packet identifier for subscribe")
		}

		if Version5 == s.Version {
			topic.Opt = &TopicOpt{
				Qos:               opts & 3,
				NoLocal:           (1 & (opts >> 2)) > 0,
				RetainAsPublished: (1 & (opts >> 3)) > 0,
				RetainHandling:    3 & (opts >> 4),
			}
		} else {
			topic.Opt = &TopicOpt{
				Qos: opts,
			}
		}

		s.Topic = append(s.Topic, topic)
		if s.Buffer.Len() == 0 {
			break
		}
	}
	return nil
}

func (s *Subscribe) encodeVariant() (result []byte, err error) {
	if s.PacketID > 0 {
		result = append(result, EncodingMSBAndLSB(s.PacketID)...)
	}
	if s.Version == Version5 {
		bs, err := EncodingRemainingLength(s.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, s.Properties.Encode(SUBSCRIBEPropType)...)
	}
	return result, nil
}

func (s *Subscribe) encodePayload() (result []byte, err error) {
	for _, topic := range s.Topic {
		// 2byte(msb+ lsb)+ variable length(topic name) + 1 byte(opts)
		result = append(result, EncodingMSBAndLSB(uint16(len(topic.Name)))...)
		result = append(result, topic.Name...)
		if Version5 == s.Version {
			var rap, nl byte
			if topic.Opt.NoLocal {
				nl = 4
			} else {
				nl = 0
			}
			if topic.Opt.RetainAsPublished {
				rap = 8
			} else {
				rap = 0
			}
			result = append(result, topic.Opt.Qos|nl|rap|topic.Opt.RetainHandling<<4)

		} else {
			result = append(result, topic.Opt.Qos)
		}
	}
	return result, nil
}
