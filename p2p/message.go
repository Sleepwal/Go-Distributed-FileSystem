package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC 可以是任意类型，两个节点之间的通信的数据结构
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
