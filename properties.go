package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type PropType byte

const (
	CONNECTPropType PropType = iota
	WILLPropType
	PUBLISHPropType
	CONNACKPropType
	PUBACKPropType
	PUBRECPropType
	PUBRELPropType
	PUBCOMPPropType
	SUBSCRIBEPropType
	SUBACKPropType
	UNSUBSCRIBEPropType
	UNSUBACKPropType
	DISCONNECTPropType
	AUTHPropType
)

// https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901046
const (
	PayloadFormatIndicator          = 0x01
	MessageExpiryInterval           = 0x02
	ContentType                     = 0x03
	ResponseTopic                   = 0x08
	CorrelationData                 = 0x09
	SubscriptionIdentifier          = 0x0B
	SessionExpiryInterval           = 0x11
	AssignedClientIdentifier        = 0x12
	ServerKeepAlive                 = 0x13
	AuthenticationMethod            = 0x15
	AuthenticationData              = 0x16
	RequestProblemInformation       = 0x17
	WillDelayInterval               = 0x18
	RequestResponseInformation      = 0x19
	ResponseInformation             = 0x1A
	ServerReference                 = 0x1C
	ReasonString                    = 0x1F
	ReceiveMaximum                  = 0x21
	TopicAliasMaximum               = 0x22
	TopicAlias                      = 0x23
	MaximumQoS                      = 0x24
	RetainAvailable                 = 0x25
	UserProperty                    = 0x26
	MaximumPacketSize               = 0x27
	WildcardSubscriptionAvailable   = 0x28
	SubscriptionIdentifierAvailable = 0x29
	SharedSubscriptionAvailable     = 0x2A
)

type User struct {
	Key   []byte
	Value []byte
}

type Properties struct {
	PayloadFormatIndicator *byte
	// 	Four Byte Integer
	MessageExpiryInterval []byte
	// UTF-8 Encoded String
	ContentType []byte
	// UTF-8 Encoded String
	ResponseTopic []byte
	// Binary Data
	CorrelationData []byte
	// Variable Byte Integer
	SubscriptionIdentifier []byte
	// Four Byte Integer
	SessionExpiryInterval []byte
	// UTF-8 Encoded String
	AssignedClientIdentifier []byte
	// Two Byte Integer
	ServerKeepAlive []byte
	// UTF-8 Encoded String
	AuthenticationMethod []byte
	// Binary Data
	AuthenticationData []byte
	// Byte
	RequestProblemInformation *byte
	// Four Byte Integer
	WillDelayInterval []byte
	// Byte
	RequestResponseInformation *byte
	// UTF-8 Encoded String
	ResponseInformation []byte
	// UTF-8 Encoded String
	ServerReference []byte
	// UTF-8 Encoded String
	ReasonString []byte
	// Two Byte Integer
	ReceiveMaximum []byte
	// Two Byte Integer
	TopicAliasMaximum []byte
	// Two Byte Integer
	TopicAlias []byte
	// Byte
	MaximumQoS *byte
	// Byte
	RetainAvailable *byte
	// UTF-8 String Pair
	UserProperty []User
	// Four Byte Integer
	MaximumPacketSize []byte
	// Byte
	WildcardSubscriptionAvailable *byte
	// Byte
	SubscriptionIdentifierAvailable *byte
	// Byte
	SharedSubscriptionAvailable *byte
	Length                      int
}

func PropertiesDecodeHandler(buffer *bytes.Buffer) (*Properties, error) {
	length, err := DecodingRemainingLength(buffer)
	if err != nil {
		return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
	}
	if length == 0 {
		return nil, nil
	}
	p := &Properties{
		Length: length,
	}
	rd := bytes.NewBuffer(buffer.Next(length))
	for {
		propertyType, err := rd.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		switch propertyType {
		case PayloadFormatIndicator:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.PayloadFormatIndicator = &b
		case MessageExpiryInterval:
			p.MessageExpiryInterval, err = ReadByteWithWidth(4, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ContentType:
			p.ContentType, err = ReadUTF8String(true, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ResponseTopic:
			p.ResponseTopic, err = ReadUTF8String(true, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case CorrelationData:
			p.CorrelationData, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case SubscriptionIdentifier:
			si, err := DecodingRemainingLength(rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			bs := make([]byte, 4)
			binary.BigEndian.PutUint32(bs, uint32(si))
			p.SubscriptionIdentifier = append(p.SubscriptionIdentifier, bs...)
		case SessionExpiryInterval:
			p.SessionExpiryInterval, err = ReadByteWithWidth(4, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case AssignedClientIdentifier:
			p.AssignedClientIdentifier, err = ReadUTF8String(true, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ServerKeepAlive:
			p.ServerKeepAlive, err = ReadByteWithWidth(2, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case AuthenticationMethod:
			if p.AuthenticationMethod != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.AuthenticationMethod, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case AuthenticationData:
			p.AuthenticationData, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case RequestProblemInformation:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.RequestProblemInformation = &b
		case WillDelayInterval:
			p.WillDelayInterval, err = ReadByteWithWidth(4, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case RequestResponseInformation:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}

			p.RequestResponseInformation = &b

		case ResponseInformation:

			p.ResponseInformation, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ServerReference:

			p.ServerReference, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ReasonString:
			p.ReasonString, err = ReadUTF8String(false, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case ReceiveMaximum:
			p.ReceiveMaximum, err = ReadByteWithWidth(2, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case TopicAliasMaximum:
			p.TopicAliasMaximum, err = ReadByteWithWidth(2, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case TopicAlias:
			p.TopicAlias, err = ReadByteWithWidth(2, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case MaximumQoS:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.MaximumQoS = &b
		case RetainAvailable:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.RetainAvailable = &b
		case UserProperty:
			k, err := ReadUTF8String(true, rd)
			if err != nil {
				return nil, err
			}
			v, err := ReadUTF8String(true, rd)
			if err != nil {
				return nil, err
			}
			p.UserProperty = append(p.UserProperty, User{Key: k, Value: v})
		case MaximumPacketSize:
			p.MaximumPacketSize, err = ReadByteWithWidth(4, rd)
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
		case WildcardSubscriptionAvailable:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.WildcardSubscriptionAvailable = &b
		case SubscriptionIdentifierAvailable:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.SubscriptionIdentifierAvailable = &b
		case SharedSubscriptionAvailable:
			b, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("%w:%s", DecodePropertiesErr, err)
			}
			p.SharedSubscriptionAvailable = &b
		}
	}
	return p, nil
}

func (p *Properties) Encode(t PropType) []byte {
	var result []byte
	switch t {
	case CONNECTPropType:
		// 	SessionExpiryInterval, AuthenticationMethod, AuthenticationData, RequestProblemInformation, RequestResponseInformation, ReceiveMaximum, TopicAliasMaximum, UserProperty, MaximumPacketSize
		if p.SessionExpiryInterval != nil {
			result = append(result, SessionExpiryInterval)
			result = append(result, p.SessionExpiryInterval...)
		}
		if p.AuthenticationMethod != nil {
			result = append(result, AuthenticationMethod)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationMethod)))...)
			result = append(result, p.AuthenticationMethod...)
		}

		if p.AuthenticationData != nil {
			result = append(result, AuthenticationData)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationData)))...)
			result = append(result, p.AuthenticationData...)
		}

		if p.RequestProblemInformation != nil {
			result = append(result, RequestProblemInformation)
			result = append(result, *p.RequestProblemInformation)
		}
		if p.RequestResponseInformation != nil {
			result = append(result, RequestResponseInformation)
			result = append(result, *p.RequestResponseInformation)
		}

		if p.ReceiveMaximum != nil {
			result = append(result, ReceiveMaximum)
			result = append(result, p.ReceiveMaximum...)
		}

		if p.TopicAliasMaximum != nil {
			result = append(result, TopicAliasMaximum)
			result = append(result, p.TopicAliasMaximum...)
		}

		if p.MaximumPacketSize != nil {
			result = append(result, MaximumPacketSize)
			result = append(result, p.MaximumPacketSize...)
		}
	case WILLPropType, PUBLISHPropType:
		// WILL: PayloadFormatIndicator, MessageExpiryInterval, ContentType, ResponseTopic, CorrelationData, WillDelayInterval, UserProperty
		// PUBLISH more than:  SubscriptionIdentifier, TopicAlias
		if p.WillDelayInterval != nil {
			result = append(result, WillDelayInterval)
			result = append(result, p.WillDelayInterval...)
		}

		if p.MessageExpiryInterval != nil {
			result = append(result, MessageExpiryInterval)
			result = append(result, p.MessageExpiryInterval...)
		}

		if p.ContentType != nil {
			result = append(result, ContentType)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ContentType)))...)
			result = append(result, p.ContentType...)
		}

		if p.PayloadFormatIndicator != nil {
			result = append(result, PayloadFormatIndicator)
			result = append(result, *p.PayloadFormatIndicator)
		}

		if p.ResponseTopic != nil {
			result = append(result, ResponseTopic)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ResponseTopic)))...)
			result = append(result, p.ResponseTopic...)
		}

		if p.CorrelationData != nil {
			result = append(result, CorrelationData)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.CorrelationData)))...)
			result = append(result, p.CorrelationData...)
		}

		if t == PUBLISHPropType {
			if p.TopicAlias != nil {
				result = append(result, TopicAlias)
				result = append(result, p.TopicAlias...)
			}
			if p.SubscriptionIdentifier != nil {
				result = append(result, SubscriptionIdentifier)
				result = append(result, p.SubscriptionIdentifier...)
			}
		}

	case SUBSCRIBEPropType:
		if p.SubscriptionIdentifier != nil {
			result = append(result, SubscriptionIdentifier)
			result = append(result, p.SubscriptionIdentifier...)
		}
	case DISCONNECTPropType:
		if p.SessionExpiryInterval != nil {
			result = append(result, SessionExpiryInterval)
			result = append(result, p.SessionExpiryInterval...)
		}
		if p.ServerReference != nil {
			result = append(result, ServerReference)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ServerReference)))...)
			result = append(result, p.ServerReference...)
		}
		if p.ReasonString != nil {
			result = append(result, ReasonString)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ReasonString)))...)
			result = append(result, p.ReasonString...)
		}
	case AUTHPropType:
		if p.AuthenticationMethod != nil {
			result = append(result, AuthenticationMethod)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationMethod)))...)
			result = append(result, p.AuthenticationMethod...)
		}
		if p.AuthenticationData != nil {
			result = append(result, AuthenticationData)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationData)))...)
			result = append(result, p.AuthenticationData...)
		}
		if p.ReasonString != nil {
			result = append(result, ReasonString)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ReasonString)))...)
			result = append(result, p.ReasonString...)
		}
	case PUBACKPropType, PUBRECPropType, PUBRELPropType, PUBCOMPPropType, SUBACKPropType, UNSUBACKPropType:
		if p.ReasonString != nil {
			result = append(result, ReasonString)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ReasonString)))...)
			result = append(result, p.ReasonString...)
		}
	case CONNACKPropType:
		// SessionExpiryInterval, AssignedClientIdentifier,ServerKeepAlive, AuthenticationMethod,
		// AuthenticationData, ResponseInformation, ServerReference, ReasonString, ReceiveMaximum,
		// TopicAliasMaximum, MaximumQoS, RetainAvailable, UserProperty, MaximumPacketSize,
		// WildcardSubscriptionAvailable, SubscriptionIdentifierAvailable, SharedSubscriptionAvailable
		if p.SessionExpiryInterval != nil {
			result = append(result, SessionExpiryInterval)
			result = append(result, p.SessionExpiryInterval...)
		}
		if p.AssignedClientIdentifier != nil {
			result = append(result, AssignedClientIdentifier)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AssignedClientIdentifier)))...)
			result = append(result, p.AssignedClientIdentifier...)
		}
		if p.ServerKeepAlive != nil {
			result = append(result, ServerKeepAlive)
			result = append(result, p.ServerKeepAlive...)
		}

		if p.AuthenticationMethod != nil {
			result = append(result, AuthenticationMethod)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationMethod)))...)
			result = append(result, p.AuthenticationMethod...)
		}
		if p.AuthenticationData != nil {
			result = append(result, AuthenticationData)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.AuthenticationData)))...)
			result = append(result, p.AuthenticationData...)
		}
		if p.RequestResponseInformation != nil {
			result = append(result, RequestResponseInformation)
			result = append(result, *p.RequestResponseInformation)
		}
		if p.ServerReference != nil {
			result = append(result, ServerReference)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ServerReference)))...)
			result = append(result, p.ServerReference...)
		}
		if p.ReasonString != nil {
			result = append(result, ReasonString)
			result = append(result, EncodingMSBAndLSB(uint16(len(p.ReasonString)))...)
			result = append(result, p.ReasonString...)
		}
		if p.ReceiveMaximum != nil {
			result = append(result, ReceiveMaximum)
			result = append(result, p.ReceiveMaximum...)
		}
		if p.TopicAliasMaximum != nil {
			result = append(result, TopicAliasMaximum)
			result = append(result, p.TopicAliasMaximum...)
		}
		if p.MaximumQoS != nil {
			result = append(result, MaximumQoS)
			result = append(result, *p.MaximumQoS)
		}
		if p.RetainAvailable != nil {
			result = append(result, RetainAvailable)
			result = append(result, *p.RetainAvailable)
		}
		if p.MaximumPacketSize != nil {
			result = append(result, MaximumPacketSize)
			result = append(result, p.MaximumPacketSize...)
		}
		if p.WildcardSubscriptionAvailable != nil {
			result = append(result, WildcardSubscriptionAvailable)
			result = append(result, *p.WildcardSubscriptionAvailable)
		}
		if p.SubscriptionIdentifierAvailable != nil {
			result = append(result, SubscriptionIdentifierAvailable)
			result = append(result, *p.SubscriptionIdentifierAvailable)
		}
		if p.SharedSubscriptionAvailable != nil {
			result = append(result, SharedSubscriptionAvailable)
			result = append(result, *p.SharedSubscriptionAvailable)
		}
	}

	if p.UserProperty != nil {
		for _, up := range p.UserProperty {
			result = append(result, UserProperty)
			result = append(result, EncodingMSBAndLSB(uint16(len(up.Key)))...)
			result = append(result, up.Key...)
			result = append(result, EncodingMSBAndLSB(uint16(len(up.Value)))...)
			result = append(result, up.Value...)
		}
	}

	return result
}
