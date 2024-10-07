package p2p

// Peer represents a remote node
type Peer interface {
	Close() error
}

// Transport is anything handling communication between network nodes
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
