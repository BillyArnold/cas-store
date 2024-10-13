package p2p

import "net"

// Peer represents a remote node
type Peer interface {
	Send([]byte) error
	net.Conn
}

// Transport is anything handling communication between network nodes
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
