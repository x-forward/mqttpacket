package packet

import (
	"bytes"
	"fmt"
)

type Disconnect struct {
	Buffer      *bytes.Buffer
	Version     byte
	FixedHeader *FixedHeader
	Properties  *Properties
}

func NewDisconnect(fh *FixedHeader, buffer *bytes.Buffer, version byte) *Disconnect {
	return &Disconnect{
		Buffer:      buffer,
		Version:     version,
		FixedHeader: fh,
	}
}

func (d *Disconnect) Encode() ([]byte, error) {
	fhBuf, err := EncodingFixedHeaderPacket(d.FixedHeader)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	_, err = d.Buffer.Write(fhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	vhBuf, err := d.encodeVariant()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}

	_, err = d.Buffer.Write(vhBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", EncodePacketErr, err)
	}
	return d.Buffer.Bytes(), nil
}

func (d *Disconnect) Decode() (*Disconnect, error) {
	if err := d.decodeVariant(); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsePacketErr, err)
	}
	d.Buffer = nil
	return d, nil
}

func (d *Disconnect) decodeVariant() (err error) {
	if d.Version == Version5 {
		d.Properties, err = PropertiesDecodeHandler(d.Buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Disconnect) encodeVariant() (result []byte, err error) {
	if d.Version == Version5 && d.Properties != nil {
		bs, err := EncodingRemainingLength(d.Properties.Length)
		if err != nil {
			return nil, err
		}
		result = append(result, bs...)
		result = append(result, d.Properties.Encode(DISCONNECTPropType)...)
	}
	return result, nil
}
