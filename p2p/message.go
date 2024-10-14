package p2p

// RPC stores any data sent over each transport between nodes
type RPC struct {
	From    string
	Payload []byte
}
