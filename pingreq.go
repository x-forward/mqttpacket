package packet

type PingReq struct {
	Version     byte
	FixedHeader *FixedHeader
}

func DecodingPingReqPacket(fh *FixedHeader) (*PingReq, error) {
	return &PingReq{FixedHeader: fh}, nil
}

func EncodingPingReqPacket(pingReq *PingReq) (result []byte, err error) {
	return EncodingFixedHeaderPacket(pingReq.FixedHeader)
}
