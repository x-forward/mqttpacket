package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingAuthPacket(t *testing.T) {
	cases := [][]byte{
		// v5
		{AUTH << 4, 2, 0, 0},
	}

	want := []*Auth{

		{
			FixedHeader: &FixedHeader{
				Type:            AUTH,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:                Version5,
			AuthenticateReasonCode: 0,
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		result, err := NewAuth(fh, rd, want[i].Version).Decode()

		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestEncodingAuthPacket(t *testing.T) {
	want := [][]byte{
		// v5
		{AUTH << 4, 2, 0},
	}

	cases := []*Auth{

		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            AUTH,
				RemainingLength: 2,
				Flag:            FixedHeaderReservedFlag,
			},
			Version:                Version5,
			AuthenticateReasonCode: 0,
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
