package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingUnSubAckPacket(t *testing.T) {
	cases := [][]byte{
		{UNSUBACK << 4, 2, 0, 10, 0x0},
	}

	want := []*UnSubAck{
		{
			FixedHeader: &FixedHeader{
				Type:            UNSUBACK,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:  Version,
			PacketID: 0xa,
			Payload:  []byte{0x0},
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)

		subAck, err := NewUnSubAck(fh, rd, Version).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], subAck)
	}
}

func TestEncodingUnSubAckPacket(t *testing.T) {
	want := [][]byte{
		{UNSUBACK << 4, 2, 0, 10, 0},
	}

	cases := []*UnSubAck{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            UNSUBACK,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:  Version,
			PacketID: 0xa,
			Payload:  []byte{0x0},
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
