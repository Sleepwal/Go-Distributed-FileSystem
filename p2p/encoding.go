package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct{}

// Decode
// @Description: 解码
// @receiver d
// @param r
// @param msg
// @return error
func (d *DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	header := make([]byte, 1)
	if _, err := r.Read(header); err != nil { // 读取第1个字节
		return nil
	}

	stream := header[0] == IncomingStream
	if stream { // 流信息
		msg.Stream = true
		return nil
	}

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}

type GobDecoder struct{}

func (g *GobDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}
