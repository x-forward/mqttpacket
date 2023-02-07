package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodingUnsubscribePacket(t *testing.T) {
	cases := []*Unsubscribe{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            UNSUBSCRIBE,
				Flag:            FixedHeaderUnsubscribeFlag,
				RemainingLength: 28,
			},
			Version:  0x00,
			PacketID: 1,
			Topic:    []string{"testtopic_1", "testtopic_2"},
		},
	}

	want := [][]byte{
		{UNSUBSCRIBE<<4 | FixedHeaderUnsubscribeFlag, 0x1c, 0x0, 0x1, 0x0, 0xb, 0x74, 0x65, 0x73, 0x74, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x31, 0x0, 0xb, 0x74, 0x65, 0x73, 0x74, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x32},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestDecodingUnsubscribePacket(t *testing.T) {
	want := []*Unsubscribe{
		{
			FixedHeader: &FixedHeader{
				Type:            UNSUBSCRIBE,
				Flag:            FixedHeaderUnsubscribeFlag,
				RemainingLength: 28,
			},
			Version:  Version,
			PacketID: 1,
			Topic:    []string{"testtopic_1", "testtopic_2"},
		},
	}

	cases := [][]byte{
		{UNSUBSCRIBE<<4 | FixedHeaderUnsubscribeFlag, 0x1c, 0x0, 0x1, 0x0, 0xb, 0x74, 0x65, 0x73, 0x74, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x31, 0x0, 0xb, 0x74, 0x65, 0x73, 0x74, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x32},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)

		result, err := NewUnsubscribe(fh, rd, Version).Decode()
		assert.Nil(t, err)

		assert.Equal(t, want[i], result)
	}
}
