package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingPubRecPacket(t *testing.T) {
	cases := [][]byte{
		// v3.x
		{PUBREC << 4, 2, 0, 10},
		// v5
		{PUBREC << 4, 3, 0, 10, 0 /*reason code*/ /* properties len */},
	}

	want := []*PubRec{
		{
			FixedHeader: &FixedHeader{
				Type:            PUBREC,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version: Version,

			PacketID: 0xa,
		},
		{
			FixedHeader: &FixedHeader{
				Type:            PUBREC,
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
		result, err := NewPubRec(fh, rd, want[i].Version).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestEncodingPubRecPacket(t *testing.T) {
	want := [][]byte{
		// v3.x
		{PUBREC << 4, 2, 0, 10},
		// v5
		{PUBREC << 4, 3, 0, 10, 0 /*reason code*/},
	}

	cases := []*PubRec{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBREC,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:  Version,
			PacketID: 0xa,
		},
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            PUBREC,
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
