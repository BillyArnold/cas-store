package p2p

const IncomingStream = 0x2
const IncomingMessage = 0x1

// RPC stores any data sent over each transport between nodes
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
