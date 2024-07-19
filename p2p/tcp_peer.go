package p2p

import (
	"net"
	"sync"
)

// TcpPeer 建立tcp连接后的一个远程节点
type TcpPeer struct {
	net.Conn
	// outbound == true：通过拨号方式连接到节点，并获取数据。
	// outbound == false：等待连接请求，建立连接，并从其中检索数据或信息。
	outbound bool
	wg       *sync.WaitGroup
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

// CloseStream
// @Description: 关闭连接
// @receiver p
func (p *TcpPeer) CloseStream() {
	p.wg.Done()
}

// Send
// @Description: 发送数据
// @receiver p
// @param b
// @return error
func (p *TcpPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}
