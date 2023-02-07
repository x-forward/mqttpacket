package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingPublishPacket(t *testing.T) {
	cases := [][]byte{
		// qos 0
		{PUBLISH<<4 | (0 | 0 | 0<<1), 33, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 123, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		// qos 1 & retain true
		{PUBLISH<<4 | (1 | 1 | 1<<1), 36, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 231, 83, 123, 32, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
	}

	want := []*Publish{
		{

			Version:     Version,
			FixedHeader: &FixedHeader{Type: PUBLISH, RemainingLength: 33, Flag: 0},
			Dup:         false,
			Qos:         0,
			Retain:      false,
			TopicName:   []byte("testtopic/#"),
			PacketID:    0, // qos 0 not have packet identifier
			Payload:     []byte{123, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		},
		{

			Version:     Version,
			FixedHeader: &FixedHeader{Type: PUBLISH, RemainingLength: 36, Flag: 3},
			Dup:         false,
			Qos:         1,
			Retain:      true,
			TopicName:   []byte("testtopic/#"),
			PacketID:    59219, // qos 0 not have packet identifier
			Payload:     []byte{123, 32, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)
		result, err := NewPublish(fh, rd, Version).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestEncodingPublishPacket(t *testing.T) {
	want := [][]byte{
		// qos 0
		{PUBLISH<<4 | (0 | 0 | 0<<1), 33, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 123, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		// qos 1 & retain true
		{PUBLISH<<4 | (1 | 1 | 1<<1), 36, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 231, 83, 123, 32, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
	}

	cases := []*Publish{
		{
			Buffer:      &bytes.Buffer{},
			Version:     Version,
			FixedHeader: &FixedHeader{Type: PUBLISH, RemainingLength: 33, Flag: 0},
			Dup:         false,
			Qos:         0,
			Retain:      false,
			TopicName:   []byte("testtopic/#"),
			PacketID:    0, // qos 0 not have packet identifier
			Payload:     []byte{123, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		},
		{
			Buffer:      &bytes.Buffer{},
			Version:     Version,
			FixedHeader: &FixedHeader{Type: PUBLISH, RemainingLength: 36, Flag: 3},
			Dup:         false,
			Qos:         1,
			Retain:      true,
			TopicName:   []byte("testtopic/#"),
			PacketID:    59219, // qos 0 not have packet identifier
			Payload:     []byte{123, 32, 10, 32, 32, 34, 109, 115, 103, 34, 58, 32, 34, 104, 101, 108, 108, 111, 34, 10, 125},
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}
