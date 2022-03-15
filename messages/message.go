package messages

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

// NewMessage allocates the first two bytes for the message length and sets the
// cursor position right after. This is the recommended way of creating
// messages.
func NewMessage() *Message {
	var p Message
	p.Buffer = make([]uint8, 2)
	p.Cursor = 2
	return &p
}

// Message is a simple buffer with functions to read write in the little endian
// ordering. First two bytes are reserved for message length.
type Message struct {
	Buffer []uint8
	Cursor uint16
}

// Checks if cursor position is past last byte in buffer
func (p *Message) overflow() bool {
	return p.Cursor >= (uint16)(len(p.Buffer))
}

// ReadUint8 reads a single byte from buffer and advances cursor.
func (p *Message) ReadUint8() uint8 {
	if p.overflow() {
		return 0
	}
	v := p.Buffer[p.Cursor]
	p.Cursor++
	return v
}

// ReadUint16 reads 2 bytes from buffer and advances cursor.
func (p *Message) ReadUint16() uint16 {
	if p.overflow() {
		return 0
	}
	v := binary.LittleEndian.Uint16(p.Buffer[p.Cursor : p.Cursor+2])
	p.Cursor += 2
	return v
}

// ReadUint32 reads 4 bytes from buffer and advances cursor.
func (p *Message) ReadUint32() uint32 {
	if p.overflow() {
		return 0
	}
	v := binary.LittleEndian.Uint32(p.Buffer[p.Cursor : p.Cursor+4])
	p.Cursor += 4
	return v
}

// ReadString reads the string length followed by the string.
func (p *Message) ReadString() string {
	if p.overflow() {
		return ""
	}
	var str string
	strlen := p.ReadUint16()
	log.Printf("New String lenght is: %v", strlen)
	for i := (uint16)(0); i < strlen; i++ {
		str += (string)(p.ReadUint8())
	}
	return str

}

// WriteUint8 writes the given byte to the message buffer. Increments message
// length by one and advances cursor.
func (p *Message) WriteUint8(v uint8) {
	p.Buffer = append(p.Buffer, v)
	binary.LittleEndian.PutUint16(p.Buffer[0:2], (uint16)(len(p.Buffer)-2))
	p.Cursor++
}

// WriteUint16 writes 2 bytes to the buffer.
func (p *Message) WriteUint16(v uint16) {
	bytes := make([]uint8, 2)
	binary.LittleEndian.PutUint16(bytes, v)
	p.WriteUint8(bytes[0])
	p.WriteUint8(bytes[1])
}

// WriteUint32 writes 4 bytes to the buffer.
func (p *Message) WriteUint32(v uint32) {
	bytes := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bytes, v)
	p.WriteUint8(bytes[0])
	p.WriteUint8(bytes[1])
	p.WriteUint8(bytes[2])
	p.WriteUint8(bytes[3])
}

// WriteString writes the string length followed by the actual string to the
// buffer.
func (p *Message) WriteString(str string) {
	p.WriteUint16((uint16)(len(str)))
	for i := 0; i < len(str); i++ {
		p.WriteUint8((uint8)(str[i]))
	}
}

// Length returns the message length stored at the first two bytes in buffer
func (p *Message) Length() uint16 {
	return binary.LittleEndian.Uint16(p.Buffer[0:2])
}

// SkipBytes advances the buffer by the given length n. If an overflow happens,
// the cursor returns to the previous state.
func (p *Message) SkipBytes(n uint16) {
	p.Cursor += n
	if p.overflow() {
		p.Cursor -= n
	}
}

// HexDump is the same as hexdump -C in the terminal and is useful for debugging
func (p *Message) HexDump(prefix string) {
	fmt.Printf("\n[%s]\n%s", prefix, hex.Dump(p.Buffer))
}
