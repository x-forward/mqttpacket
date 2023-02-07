package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
)

type Connect struct {
	Buffer      *bytes.Buffer
	FixedHeader *FixedHeader
	// variable header
	ProtocolName  []byte
	ProtocolLevel byte
	KeepAlive     uint16
	Flag          *Flag
	Properties    *Properties
	// payload
	ClientID       []byte
	WillProperties *Properties
	WillTopic      []byte
	WillMessage    []byte
	Username       []byte
	Password       []byte
}

type Flag struct {
	UserName     bool
	Password     bool
	WillRetain   bool
	WillQos      uint8
	Will         bool
	CleanSession bool
	Reserved     uint8
}

func NewConnect(fh *FixedHeader, buffer *bytes.Buffer) *Connect {
	return &Connect{
		Buffer:      buffer,
		FixedHeader: fh,
	}
}

func (c *Connect) Encode() (result []byte, err error) {
	fhBuf, err := EncodingFixedHeaderPacket(c.FixedHeader)
	if err != nil {
		return nil, err
	}

	_, err = c.Buffer.Write(fhBuf)
	if err != nil {
		return nil, err
	}

	vhBuf, err := c.encodeVariant()
	if err != nil {
		return nil, err
	}

	_, err = c.Buffer.Write(vhBuf)

	if err != nil {
		return nil, err
	}

	payload, err := c.encodePayload()
	if err != nil {
		return nil, err
	}

	_, err = c.Buffer.Write(payload)

	if err != nil {
		return nil, EncodePacketErr
	}
	return c.Buffer.Bytes(), nil

}

func (c *Connect) Decode() (*Connect, error) {
	if err := c.decodeVariant(); err != nil {
		return nil, err
	}

	if err := c.decodePayload(); err != nil {
		return nil, err
	}
	c.Buffer = nil
	return c, nil
}

func (c *Connect) decodeVariant() error {
	protocolName, err := ReadUTF8String(true, c.Buffer)
	if err != nil {
		return fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	protocolLevel, err := c.Buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	flag, err := c.Buffer.ReadByte()
	if err != nil || 1&flag != 0 {
		return fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	f := c.decodeFlag(flag)
	if f.WillRetain && !f.Will {
		return fmt.Errorf("%w: %s", ParsePacketErr, errors.New("will retain flag conflict with will flag"))
	}
	if !f.Will && f.WillQos > 0 {
		return fmt.Errorf("%w: %s", ParsePacketErr, errors.New("will qos flag conflict with will flag"))
	}
	ka := make([]byte, 2)
	if _, err := c.Buffer.Read(ka); err != nil {
		return fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	keepAlive := binary.BigEndian.Uint16(ka)

	c.ProtocolName = protocolName
	c.ProtocolLevel = protocolLevel
	c.Flag = f
	c.KeepAlive = keepAlive
	if protocolLevel >= Version5 {
		connectProperties, err := PropertiesDecodeHandler(c.Buffer)
		if err != nil {
			return err
		}
		c.Properties = connectProperties
	}
	return nil
}

func (c *Connect) encodeVariant() (result []byte, err error) {
	result = append(result, EncodingMSBAndLSB(uint16(len(c.ProtocolName)))...)
	result = append(result, c.ProtocolName...)
	// protocol level
	result = append(result, c.ProtocolLevel)
	// flag
	result = append(result, c.encodeFlag())
	// keepalive
	result = append(result, EncodingMSBAndLSB(c.KeepAlive)...)

	if c.ProtocolLevel >= Version5 && c.Properties != nil {
		bs, err := EncodingRemainingLength(c.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, c.Properties.Encode(CONNECTPropType)...)
	}
	return result, nil
}

func (c *Connect) decodePayload() error {
	cid, err := ReadUTF8String(true, c.Buffer)
	if err != nil {
		return err
	}
	c.ClientID = cid

	if c.Flag.Will {
		if c.ProtocolLevel >= Version5 {
			willProperties, err := PropertiesDecodeHandler(c.Buffer)
			if err != nil {
				return err
			}
			c.WillProperties = willProperties
		}
		willTopic, err := ReadUTF8String(true, c.Buffer)
		if err != nil {
			return err
		}
		willMsg, err := ReadUTF8String(true, c.Buffer)
		if err != nil {
			return err
		}
		c.WillTopic = willTopic
		c.WillMessage = willMsg
	}

	if c.Flag.UserName && c.Flag.Password {
		u, err := ReadUTF8String(true, c.Buffer)
		if err != nil {
			return err

		}
		p, err := ReadUTF8String(true, c.Buffer)
		if err != nil {
			return err
		}
		c.Password = p
		c.Username = u
	}
	return nil
}

func (c *Connect) encodePayload() (result []byte, err error) {
	// client
	result = append(result, EncodingMSBAndLSB(uint16(len(c.ClientID)))...)
	result = append(result, c.ClientID...)
	if c.Flag.Will {
		if c.ProtocolLevel >= Version5 {
			bs, err := EncodingRemainingLength(c.WillProperties.Length)
			if err != nil {
				return nil, err
			}
			result = append(result, bs...)
			result = append(result, c.WillProperties.Encode(WILLPropType)...)
		}
		if c.WillTopic != nil {
			result = append(result, EncodingMSBAndLSB(uint16(len(c.WillTopic)))...)
			result = append(result, c.WillTopic...)
		}
		if c.WillMessage != nil {
			result = append(result, EncodingMSBAndLSB(uint16(len(c.WillMessage)))...)
			result = append(result, c.WillMessage...)
		}
	}

	// username
	if c.Flag.UserName {
		result = append(result, EncodingMSBAndLSB(uint16(len(c.Username)))...)
		result = append(result, c.Username...)
	}
	// password
	if c.Flag.Password {
		result = append(result, EncodingMSBAndLSB(uint16(len(c.Password)))...)
		result = append(result, c.Password...)
	}

	return result, nil
}

func (c *Connect) decodeFlag(flag byte) *Flag {
	return &Flag{
		UserName:     (1 & (flag >> 7)) > 0,
		Password:     (1 & (flag >> 6)) > 0,
		WillRetain:   (1 & (flag >> 5)) > 0,
		WillQos:      3 & (flag >> 3),
		Will:         (1 & (flag >> 2)) > 0, //
		CleanSession: (1 & (flag >> 1)) > 0,
		Reserved:     0,
	}
}

func (c *Connect) encodeFlag() byte {
	var (
		username     = 0
		password     = 0
		willRetain   = 0
		will         = 0
		cleanSession = 0
		reserved     = 0
	)
	if c.Flag.UserName {
		username = 1 << 7
	}
	if c.Flag.Password {
		password = 1 << 6
	}
	if c.Flag.WillRetain {
		willRetain = 1 << 5
	}

	if c.Flag.Will {
		will = 4
	}
	if c.Flag.CleanSession {
		cleanSession = 1 << 1
	}
	qosFlag := 0
	switch c.Flag.WillQos {
	case 1:
		qosFlag = 8
	case 2:
		qosFlag = 16
	}
	flag := username | password | willRetain | qosFlag | will | cleanSession | reserved
	return uint8(flag)
}
