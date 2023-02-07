package packet

type PingResp struct {
	Version     byte
	FixedHeader *FixedHeader
}

func DecodingPingRespPacket(fh *FixedHeader) (*PingResp, error) {
	return &PingResp{FixedHeader: fh}, nil
}

func EncodingPingRespPacket(pingResp *PingResp) (result []byte, err error) {
	return EncodingFixedHeaderPacket(pingResp.FixedHeader)
}
