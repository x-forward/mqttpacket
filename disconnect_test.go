package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodingDisconnectPacket(t *testing.T) {
	cases := []*Disconnect{
		{Buffer: &bytes.Buffer{}, FixedHeader: &FixedHeader{Type: DISCONNECT, Flag: FixedHeaderReservedFlag}, Version: Version5},
	}

	want := [][]byte{
		{DISCONNECT << 4, 0},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.NoError(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestDecodingDisconnectPacket(t *testing.T) {
	want := []*Disconnect{
		{FixedHeader: &FixedHeader{Type: DISCONNECT, Flag: FixedHeaderReservedFlag}, Version: Version5},
	}

	cases := [][]byte{
		{DISCONNECT << 4, 0},
	}
	for i, c := range cases {
		rd := bytes.NewBuffer(c)

		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)

		result, err := NewDisconnect(fh, rd, Version5).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
