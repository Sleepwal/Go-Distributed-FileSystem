package p2p

import "net"

// TcpPeer 建立tcp连接后的一个远程节点
type TcpPeer struct {
	conn net.Conn
	// outbound == true：通过拨号方式连接到节点，并获取数据。
	// outbound == false：等待连接请求，建立连接，并从其中检索数据或信息。
	outbound bool
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		conn:     conn,
		outbound: outbound,
	}
}
