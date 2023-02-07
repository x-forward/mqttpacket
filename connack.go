package packet

import (
	"bytes"
	"fmt"
)

const (
	ConnAckAccepted                           = 0x00
	ConnAckRefusedWithInvalidMqttProtocol     = 0x01
	ConnAckRefusedWithInvalidClientID         = 0x02
	ConnAckRefusedWithInvalidServer           = 0x03
	ConnAckRefusedWithInvalidUsernamePassword = 0x04
	ConnAckRefusedServerRejected              = 0x05
)

type ConnAck struct {
	Buffer         *bytes.Buffer
	FixedHeader    *FixedHeader
	SessionPresent byte
	ResponseCode   byte
	Version        byte
	Properties     *Properties
}

func NewConnAck(fh *FixedHeader, buffer *bytes.Buffer, version byte) *ConnAck {
	return &ConnAck{
		FixedHeader: fh,
		Buffer:      buffer,
		Version:     version,
	}
}

func (c *ConnAck) Encode() ([]byte, error) {
	fhBuf, err := EncodingFixedHeaderPacket(c.FixedHeader)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	_, err = c.Buffer.Write(fhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	vhBuf, err := c.encodeVariant()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	_, err = c.Buffer.Write(vhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	return c.Buffer.Bytes(), nil
}

func (c *ConnAck) Decode() (*ConnAck, error) {
	if err := c.decodeVariant(); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	c.Buffer = nil
	return c, nil
}

func (c *ConnAck) decodeVariant() error {
	sp, err := c.Buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("decoding connack session present error, %w", err)
	}
	rc, err := c.Buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("decoding connack response code error, %w", err)
	}
	c.SessionPresent = sp
	c.ResponseCode = rc
	if c.Version == Version5 {
		c.Properties, err = PropertiesDecodeHandler(c.Buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConnAck) encodeVariant() (result []byte, err error) {
	result = append(result, c.SessionPresent, c.ResponseCode)
	if c.Version == Version5 && c.Properties != nil {
		bs, err := EncodingRemainingLength(c.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, c.Properties.Encode(CONNACKPropType)...)
	}
	return result, nil
}
