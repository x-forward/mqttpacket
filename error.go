package packet

import "errors"

var (
	DecodePropertiesErr = errors.New("decode property packet err")
	ReadUTF8BufferErr   = errors.New("read utf-8 buffer err")
	ParsePacketErr      = errors.New("parse packet err")
	EncodePacketErr     = errors.New("encode packet err")
)
