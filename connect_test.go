package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func intWrapPtr(s uint8) *byte {
	return &s
}

func TestDecodingConnectPacket(t *testing.T) {
	cases := [][]byte{
		{CONNECT << 4, 41, 0, 4,
			77, 81, 84, 84, 4, 194, 0, 60, 0, 14, 109, 113, 116, 116, 120, 95, 100, 51, 98, 49, 99, 56, 98, 99, 0, 5, 97, 100, 109, 105, 110, 0, 6, 49, 50, 51, 52, 53, 54},
		{CONNECT << 4, 26, 0, 4,
			77, 81, 84, 84, 4, 0, 0, 120, 0, 14, 109, 113, 116, 116, 120, 95, 100, 51, 98, 49, 99, 56, 98, 99},
		// with  properties
		{
			CONNECT << 4, 124, 0, 4,
			77, 81, 84, 84, 5, 238, 0, 60, 8, 17, 0, 0, 0, 30, 33, 99, 211, 0, 14, 109, 113,
			116, 116, 120, 95, 54, 54, 49, 54, 52, 55, 51, 97, 42, 24, 0, 0, 0, 10, 2, 0, 0, 0, 10,
			3, 0, 27, 123, 10, 32, 32, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 32, 34, 100, 101,
			115, 99, 114, 105, 98, 101, 34, 10, 125, 1, 1, 0, 4, 119, 105, 108, 108, 0, 23, 123, 10,
			32, 32, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 32, 34, 119, 105, 108, 108, 34, 10,
			125, 0, 5, 97, 100, 109, 105, 110, 0, 6, 49, 50, 51, 52, 53, 54,
		},
	}

	want := []*Connect{
		{
			FixedHeader: &FixedHeader{
				Type:            CONNECT,
				RemainingLength: 41,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_d3b1c8bc"),
			Username:      []byte("admin"),
			Password:      []byte("123456"),
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x04,
			KeepAlive:     60,
			Flag: &Flag{
				UserName:     true,
				Password:     true,
				WillRetain:   false,
				WillQos:      0,
				Will:         false,
				CleanSession: true,
				Reserved:     0,
			},
		},
		{
			FixedHeader: &FixedHeader{
				Type:            0x1,
				RemainingLength: 26,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_d3b1c8bc"),
			Username:      nil,
			Password:      nil,
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x04,
			KeepAlive:     120,
			Flag: &Flag{
				UserName:     false,
				Password:     false,
				WillRetain:   false,
				WillQos:      0,
				Will:         false,
				CleanSession: false,
				Reserved:     0,
			},
		},
		{
			FixedHeader: &FixedHeader{
				Type:            0x1,
				RemainingLength: 124,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_6616473a"),
			Username:      []byte("admin"),
			Password:      []byte("123456"),
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x05,
			KeepAlive:     60,
			Flag: &Flag{
				UserName:     true,
				Password:     true,
				WillRetain:   true,
				WillQos:      1,
				Will:         true,
				CleanSession: true,
				Reserved:     0,
			},
			// WILL: PayloadFormatIndicator, MessageExpiryInterval, ContentType, ResponseTopic, CorrelationData, WillDelayInterval, UserProperty
			WillProperties: &Properties{
				PayloadFormatIndicator: intWrapPtr(1),
				MessageExpiryInterval:  []byte{0x00, 0x00, 0x00, 0xa},
				ContentType:            []byte{0x7b, 0x0a, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x22, 0x0a, 0x7d},
				ResponseTopic:          nil,
				WillDelayInterval:      []byte{0x00, 0x00, 0x00, 0xa},
				CorrelationData:        nil,
				UserProperty:           nil,
				Length:                 42,
			},
			// SessionExpiryInterval, AuthenticationMethod, AuthenticationData, RequestProblemInformation, RequestResponseInformation, ReceiveMaximum, TopicAliasMaximum, UserProperty, MaximumPacketSize
			Properties: &Properties{
				SessionExpiryInterval:      []byte{0x00, 0x00, 0x00, 0x1e},
				AuthenticationMethod:       nil,
				AuthenticationData:         nil,
				RequestProblemInformation:  nil,
				RequestResponseInformation: nil,
				ResponseInformation:        nil,
				ReceiveMaximum:             []byte{0x63, 0xd3},
				TopicAliasMaximum:          nil,
				UserProperty:               nil,
				MaximumPacketSize:          nil,
				Length:                     8,
			},
			WillTopic:   []byte("will"),
			WillMessage: []byte{0x7b, 0x0a, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x77, 0x69, 0x6c, 0x6c, 0x22, 0x0a, 0x7d},
		},
	}

	for i, c := range cases {
		rd := bytes.NewBuffer(c)
		fh, err := DecodingFixedHeaderPacket(rd)
		assert.Nil(t, err)
		result, err := NewConnect(fh, rd).Decode()
		assert.Nil(t, err)
		assert.Equal(t, want[i], result)
	}
}

func TestEncodingConnectPacket(t *testing.T) {
	want := [][]byte{
		{CONNECT << 4, 41, 0, 4,
			77, 81, 84, 84, 4, 194, 0, 60, 0, 14, 109, 113, 116, 116, 120, 95, 100, 51, 98, 49, 99, 56, 98, 99, 0, 5, 97, 100, 109, 105, 110, 0, 6, 49, 50, 51, 52, 53, 54},
		{CONNECT << 4, 26, 0, 4,
			77, 81, 84, 84, 4, 0, 0, 120, 0, 14, 109, 113, 116, 116, 120, 95, 100, 51, 98, 49, 99, 56, 98, 99},
		//with  properties
		{
			CONNECT << 4, 124, 0, 4,
			77, 81, 84, 84, 5, 238, 0, 60, 8, 17, 0, 0, 0, 30, 33, 99, 211, 0, 14, 109, 113,
			116, 116, 120, 95, 54, 54, 49, 54, 52, 55, 51, 97, 42, 24, 0, 0, 0, 10, 2, 0, 0, 0, 10,
			3, 0, 27, 123, 10, 32, 32, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 32, 34, 100, 101,
			115, 99, 114, 105, 98, 101, 34, 10, 125, 1, 1, 0, 4, 119, 105, 108, 108, 0, 23, 123, 10,
			32, 32, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 32, 34, 119, 105, 108, 108, 34, 10,
			125, 0, 5, 97, 100, 109, 105, 110, 0, 6, 49, 50, 51, 52, 53, 54,
		},
	}

	cases := []*Connect{
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            CONNECT,
				RemainingLength: 41,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_d3b1c8bc"),
			Username:      []byte("admin"),
			Password:      []byte("123456"),
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x04,
			KeepAlive:     60,
			Flag: &Flag{
				UserName:     true,
				Password:     true,
				WillRetain:   false,
				WillQos:      0,
				Will:         false,
				CleanSession: true,
				Reserved:     0,
			},
		},
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            0x1,
				RemainingLength: 26,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_d3b1c8bc"),
			Username:      nil,
			Password:      nil,
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x04,
			KeepAlive:     120,
			Flag: &Flag{
				UserName:     false,
				Password:     false,
				WillRetain:   false,
				WillQos:      0,
				Will:         false,
				CleanSession: false,
				Reserved:     0,
			},
		},
		{
			Buffer: &bytes.Buffer{},
			FixedHeader: &FixedHeader{
				Type:            0x1,
				RemainingLength: 124,
				Flag:            0x0,
			},
			ClientID:      []byte("mqttx_6616473a"),
			Username:      []byte("admin"),
			Password:      []byte("123456"),
			ProtocolName:  []byte("MQTT"),
			ProtocolLevel: 0x05,
			KeepAlive:     60,
			Flag: &Flag{
				UserName:     true,
				Password:     true,
				WillRetain:   true,
				WillQos:      1,
				Will:         true,
				CleanSession: true,
				Reserved:     0,
			},
			// WILL: PayloadFormatIndicator, MessageExpiryInterval, ContentType, ResponseTopic, CorrelationData, WillDelayInterval, UserProperty
			WillProperties: &Properties{
				PayloadFormatIndicator: intWrapPtr(1),
				MessageExpiryInterval:  []byte{0x00, 0x00, 0x00, 0xa},
				ContentType:            []byte{0x7b, 0x0a, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x22, 0x0a, 0x7d},
				ResponseTopic:          nil,
				WillDelayInterval:      []byte{0x00, 0x00, 0x00, 0xa},
				CorrelationData:        nil,
				UserProperty:           nil,
				Length:                 42,
			},
			// SessionExpiryInterval, AuthenticationMethod, AuthenticationData, RequestProblemInformation, RequestResponseInformation, ReceiveMaximum, TopicAliasMaximum, UserProperty, MaximumPacketSize
			Properties: &Properties{
				SessionExpiryInterval:      []byte{0x00, 0x00, 0x00, 0x1e},
				AuthenticationMethod:       nil,
				AuthenticationData:         nil,
				RequestProblemInformation:  nil,
				RequestResponseInformation: nil,
				ResponseInformation:        nil,
				ReceiveMaximum:             []byte{0x63, 0xd3},
				TopicAliasMaximum:          nil,
				UserProperty:               nil,
				MaximumPacketSize:          nil,
				Length:                     8,
			},
			WillTopic:   []byte("will"),
			WillMessage: []byte{0x7b, 0x0a, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x77, 0x69, 0x6c, 0x6c, 0x22, 0x0a, 0x7d},
		},
	}

	for i, c := range cases {
		result, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, len(want[i]), len(result))
		assert.Equal(t, want[i], result)
	}
}
