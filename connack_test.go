package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodingConnAckPacket(t *testing.T) {
	cases := []*ConnAck{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            CONNACK,
				RemainingLength: 2,
				Flag:            0,
			},
			Version:        Version5,
			SessionPresent: uint8(1),
			ResponseCode:   ConnAckAccepted,
		},

		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            CONNACK,
				RemainingLength: 2,
				Flag:            0,
			},
			SessionPresent: byte(0),
			Version:        Version5,
			ResponseCode:   ConnAckRefusedServerRejected,
		},
	}

	want := [][]byte{
		{0x20, 0x02, 0x01, 0x00},
		{0x20, 0x02, 0x00, 0x05},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestDecodingConnAckPacket(t *testing.T) {
	cases := [][]byte{
		{0x20, 0x02, 0x01, 0x00},
		{0x20, 0x02, 0x00, 0x05},
	}
	want := []*ConnAck{
		{

			FixedHeader: &FixedHeader{
				Type:            CONNACK,
				RemainingLength: 2,
				Flag:            0,
			},
			SessionPresent: byte(1),
			Version:        Version5,
			ResponseCode:   ConnAckAccepted,
		},

		{

			FixedHeader: &FixedHeader{
				Type:            CONNACK,
				RemainingLength: 2,
				Flag:            0,
			},
			SessionPresent: byte(0),
			Version:        Version5,
			ResponseCode:   ConnAckRefusedServerRejected,
		},
	}
	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)
		result, err := NewConnAck(fh, rd, Version5).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
