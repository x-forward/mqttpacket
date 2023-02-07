package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingSubAckPacket(t *testing.T) {
	cases := [][]byte{
		{SUBACK << 4, 5, 0, 10, 0, 1, 2},
	}

	want := []*SubAck{
		{
			FixedHeader: &FixedHeader{
				Type:            SUBACK,
				RemainingLength: 5,
				Flag:            0,
			},
			Version:  Version,
			PacketID: 0xa,
			Payload:  []byte{0x0, 0x01, 0x02},
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)
		subAck, err := NewSubAck(fh, rd, Version).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], subAck)
	}
}

func TestEncodingSubAckPacket(t *testing.T) {
	want := [][]byte{
		{SUBACK << 4, 5, 0, 10, 0, 1, 2},
	}

	cases := []*SubAck{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            SUBACK,
				RemainingLength: 5,
				Flag:            0,
			},
			Version:  Version,
			PacketID: 0xa,
			Payload:  []byte{0x0, 0x01, 0x02},
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
