package p2p

import (
	"log"
	"net"
	"sync"
)

type TcpTransport struct {
	listenAddress string       // 监听地址
	listener      net.Listener // 监听器

	mu    sync.Mutex        // 锁
	peers map[net.Addr]Peer // 节点集合
}

func NewTCPTransport(listenAddr string) Transport {
	return &TcpTransport{
		listenAddress: listenAddr,
	}
}

func (t *TcpTransport) ListenAddr() string {
	return t.listenAddress
}

func (t *TcpTransport) ListenAndAccept() (err error) {
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return
	}

	go t.AcceptLoop()
	return
}

func (t *TcpTransport) AcceptLoop() {
	conn, err := t.listener.Accept() // 接收连接
	if err != nil {
		log.Println("TCP accept error:", err)
	}
	go t.handleConn(conn) // 处理连接
}

func (t *TcpTransport) handleConn(conn net.Conn) {
	peer := NewTcpPeer(conn, true)
	log.Printf("a new connection comes: %+v\n", peer)
}
