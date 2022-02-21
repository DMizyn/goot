package network

import (
	"github.com/rwxsu/goot/messages"
	"net"
)

const debug = true

// RecvMessage reads the incoming message length (first two bytes), followed by
// how many bytes the incoming message length is.
func RecvMessage(c net.Conn) *messages.Message {
	msg := messages.NewMessage()
	c.Read(msg.Buffer[0:2]) // incoming message length
	if msg.Length() == 0 {
		return nil
	}
	bytes := make([]uint8, msg.Length())
	c.Read(bytes)
	msg.Buffer = append(msg.Buffer, bytes...)
	if debug {
		msg.HexDump("recv")
	}
	return msg
}

// SendMessage sends a message to the given connection.
func SendMessage(dest net.Conn, msg *messages.Message) {
	dest.Write(msg.Buffer)
	if debug {
		msg.HexDump("send")
	}
}
