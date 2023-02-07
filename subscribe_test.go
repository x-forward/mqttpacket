package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodingSubscribePacket(t *testing.T) {
	cases := [][]byte{
		{SUBSCRIBE<<4 | FixedHeaderSubscribeFlag, 16, 13, 111, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 0},
		{SUBSCRIBE<<4 | FixedHeaderSubscribeFlag, 18, 210, 134, 0, 13, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 115, 47, 35, 0},
	}

	want := []*Subscribe{
		{
			FixedHeader: &FixedHeader{
				Type:            SUBSCRIBE,
				RemainingLength: 16,
				Flag:            FixedHeaderSubscribeFlag,
			},
			Version:  Version,
			PacketID: 3439,
			Topic: []Topic{
				{
					Opt: &TopicOpt{
						Qos:               0,
						RetainHandling:    0,
						NoLocal:           false,
						RetainAsPublished: false,
					},
					Name: []byte("testtopic/#"),
				},
			},
		},

		{
			FixedHeader: &FixedHeader{
				Type:            SUBSCRIBE,
				RemainingLength: 18,
				Flag:            FixedHeaderSubscribeFlag,
			},
			Version:  Version,
			PacketID: 0xd286,
			Topic: []Topic{
				{
					Opt: &TopicOpt{
						Qos:               0,
						RetainHandling:    0,
						NoLocal:           false,
						RetainAsPublished: false,
					},
					Name: []byte("testtopic/s/#"),
				},
			},
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)

		sub, err := NewSubscribe(fh, rd, Version).Decode()
		assert.Nil(t, err)

		assert.Equal(t, want[i], sub)
	}
}

func TestEncodingSubscribePacket(t *testing.T) {
	want := [][]byte{
		{SUBSCRIBE<<4 | FixedHeaderSubscribeFlag, 16, 13, 111, 0, 11, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 35, 0},
		{SUBSCRIBE<<4 | FixedHeaderSubscribeFlag, 18, 210, 134, 0, 13, 116, 101, 115, 116, 116, 111, 112, 105, 99, 47, 115, 47, 35, 0},
	}

	cases := []*Subscribe{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            SUBSCRIBE,
				RemainingLength: 16,
				Flag:            FixedHeaderSubscribeFlag,
			},
			Version:  Version,
			PacketID: 3439,
			Topic: []Topic{
				{
					Opt: &TopicOpt{
						Qos:               0,
						RetainHandling:    0,
						NoLocal:           false,
						RetainAsPublished: false,
					},
					Name: []byte("testtopic/#"),
				},
			},
		},

		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            SUBSCRIBE,
				RemainingLength: 18,
				Flag:            FixedHeaderSubscribeFlag,
			},
			Version:  Version,
			PacketID: 0xd286,
			Topic: []Topic{
				{
					Opt: &TopicOpt{
						Qos:               0,
						RetainHandling:    0,
						NoLocal:           false,
						RetainAsPublished: false,
					},
					Name: []byte("testtopic/s/#"),
				},
			},
		},
	}

	for i, c := range cases {
		actual, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], actual)
	}
}
