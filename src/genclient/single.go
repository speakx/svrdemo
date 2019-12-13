package genclient

import (
	"encoding/binary"
	"net"

	"github.com/smallnest/goframe"
)

type GenSingleClient struct {
	FC goframe.FrameConn
}

func (gsc *GenSingleClient) Connect() {
	conn, err := net.Dial("tcp", "10.211.55.27:11000")
	if err != nil {
		panic(err)
	}

	encoderConfig := goframe.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}

	decoderConfig := goframe.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}

	gsc.FC = goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)
}
