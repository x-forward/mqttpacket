package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingPubRelPacket(t *testing.T) {
	cases := [][]byte{
		// v3.x
		{PUBREL << 4, 2, 0, 10},
		// v5
		{PUBREL << 4, 3, 0, 10, 0 /*reason code*/ /* properties len */},
	}

	want := []*PubRel{
		{
			FixedHeader: &FixedHeader{
				Type:            PUBREL,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version: Version,

			PacketID: 0xa,
		},
		{
			FixedHeader: &FixedHeader{
				Type:            PUBREL,
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

func TestEncodingPubRelPacket(t *testing.T) {
	want := [][]byte{
		// v3.x
		{PUBREL << 4, 2, 0, 10},
		// v5
		{PUBREL << 4, 3, 0, 10, 0 /*reason code*/},
	}

	cases := []*PubRel{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBREL,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:  Version,
			PacketID: 0xa,
		},
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBREL,
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
