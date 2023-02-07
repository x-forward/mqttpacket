package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingPubCompPacket(t *testing.T) {
	cases := [][]byte{
		// v3.x
		{PUBCOMP << 4, 2, 0, 10},
		// v5
		{PUBCOMP << 4, 3, 0, 10, 0 /*reason code*/ /* properties len */},
	}

	want := []*PubRel{
		{
			FixedHeader: &FixedHeader{
				Type:            PUBCOMP,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version: Version,

			PacketID: 0xa,
		},
		{
			FixedHeader: &FixedHeader{
				Type:            PUBCOMP,
				RemainingLength: 3,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:    Version5,
			ReasonCode: 0,
			PacketID:   0xa,
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		result, err := NewPubRel(fh, rd, want[i].Version).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestEncodingPubCompPacket(t *testing.T) {
	want := [][]byte{
		// v3.x
		{PUBCOMP << 4, 2, 0, 10},
		// v5
		{PUBCOMP << 4, 3, 0, 10, 0 /*reason code*/},
	}

	cases := []*PubComp{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBCOMP,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:  Version,
			PacketID: 0xa,
		},
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBCOMP,
				RemainingLength: 3,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:    Version5,
			ReasonCode: 0,
			PacketID:   0xa,
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
