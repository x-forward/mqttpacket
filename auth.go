package packet

import (
	"bytes"
	"fmt"
)

type Auth struct {
	Buffer                 *bytes.Buffer
	Version                byte
	FixedHeader            *FixedHeader
	Properties             *Properties
	AuthenticateReasonCode int
}

func NewAuth(fh *FixedHeader, buffer *bytes.Buffer, version byte) *Auth {
	return &Auth{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
	}
}

func (p *Auth) Encode() ([]byte, error) {
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

func (p *Auth) Decode() (*Auth, error) {
	if err := p.decodeVariant(); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	p.Buffer = nil
	return p, nil
}

func (p *Auth) decodeVariant() (err error) {
	code, err := p.Buffer.ReadByte()
	if err != nil {
		return err
	}
	p.AuthenticateReasonCode = int(code)
	if p.Version == Version5 {

		p.Properties, err = PropertiesDecodeHandler(p.Buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Auth) encodeVariant() (result []byte, err error) {
	result = append(result, byte(p.AuthenticateReasonCode))
	if p.Version == Version5 && p.Properties != nil {
		bs, err := EncodingRemainingLength(p.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, p.Properties.Encode(AUTHPropType)...)
	}
	return result, nil
}
