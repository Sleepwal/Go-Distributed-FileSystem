package main

import (
	"Go-Distributed-FileSystem/p2p"
	"log"
)

func main() {
	conf := p2p.TcpTransportConfig{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
	}
	transport := p2p.NewTCPTransport(conf)
	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
