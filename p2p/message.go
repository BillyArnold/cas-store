package p2p

import "net"

// RPC stores any data sent over each transport between nodes
type RPC struct {
	From    net.Addr
	Payload []byte
}
