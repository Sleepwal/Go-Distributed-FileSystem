package p2p

import (
	"errors"
	"log"
	"net"
)

type TcpTransportConfig struct {
	ListenAddress string
	HandshakeFunc
	Decoder
	OnPeer func(Peer) error
}

type TcpTransport struct {
	TcpTransportConfig
	listener net.Listener // 监听器
	rpcChan  chan RPC
}

func NewTCPTransport(conf TcpTransportConfig) Transport {
	return &TcpTransport{
		TcpTransportConfig: conf,
		rpcChan:            make(chan RPC, 1024),
	}
}

// Addr
// @Description: 返回监听地址
// @receiver t
// @return string
func (t *TcpTransport) Addr() string {
	return t.ListenAddress
}

// Consume
// @Description: 读取另一个节点的信息，该通道只读
// @receiver t
// @return <-chan
func (t *TcpTransport) Consume() <-chan RPC {
	return t.rpcChan
}

func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func (t *TcpTransport) Close() error { return t.listener.Close() }

func (t *TcpTransport) ListenAndAccept() (err error) {
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return
	}

	go t.AcceptLoop()
	log.Printf("tcp listening on %s\n", t.ListenAddress)
	return
}

func (t *TcpTransport) AcceptLoop() {
	for {
		conn, err := t.listener.Accept() // 接收连接
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			log.Println("TCP accept error:", err)
		}
		go t.handleConn(conn, false) // 处理连接
	}
}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		log.Printf("peer[%s] connection closed: %s\n", conn.RemoteAddr(), err)
		_ = conn.Close()
	}()

	peer := NewTcpPeer(conn, outbound)
	log.Printf("a new connection comes: %+v\n", peer)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// 循环读取消息
	for {
		rpc := RPC{}
		err = t.Decode(conn, &rpc) // 解码
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			log.Printf("[%s] incoming stream rpc: %+v\n", conn.RemoteAddr(), rpc)
			peer.wg.Wait()
			log.Printf("[%s] stream rpc done: %+v\n", conn.RemoteAddr(), rpc)
			continue
		}

		t.rpcChan <- rpc
	}
}
