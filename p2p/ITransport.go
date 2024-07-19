package p2p

// Peer 远程节点
type Peer interface {
}

// Transport 处理节点之间的通信（使用tcp、UDP、WebSocket）
type Transport interface {
	Addr() string
	ListenAndAccept() error
	Dial(string) error
	Consume() <-chan RPC
}
