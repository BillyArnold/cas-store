package p2p

import "net"

// Peer represents a remote node
type Peer interface {
	Send([]byte) error
	net.Conn
	CloseStream()
}

// Transport is anything handling communication between network nodes
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
