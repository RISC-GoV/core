package core

import "bytes"

type Memory struct {
	mem []byte
}

func NewMemory() *Memory {
	return &Memory{
		mem: make([]byte, 1048576), //1MB
	}
}

func (m *Memory) ReadSingleByte(addr uint32) uint32 {
	return uint32(m.mem[addr])
}

func (m *Memory) ReadByte(addr uint32) byte {
	return m.mem[addr]
}

func (m *Memory) WriteSingleByte(addr uint32, val uint32) {
	m.mem[addr] = byte(val)
}

func (m *Memory) WriteByte(addr uint32, val byte) {
	m.mem[addr] = byte(val)
}

func (m *Memory) ReadHalfWord(addr uint32) uint32 {
	return uint32(m.mem[addr]) | uint32(m.mem[addr+1])<<8
}

func (m *Memory) WriteHalfWord(addr uint32, val uint32) {
	m.mem[addr] = byte(val)
	m.mem[addr+1] = byte(val >> 8)
}

func (m *Memory) ReadWord(addr uint32) uint32 {
	return uint32(m.mem[addr]) | uint32(m.mem[addr+1])<<8 | uint32(m.mem[addr+2])<<16 | uint32(m.mem[addr+3])<<24
}

func (m *Memory) WriteWord(addr uint32, val uint32) {
	m.mem[addr] = byte(val)
	m.mem[addr+1] = byte(val >> 8)
	m.mem[addr+2] = byte(val >> 16)
	m.mem[addr+3] = byte(val >> 24)
}

func (m *Memory) ReadString(addr uint32) string {
	strbytes := make([]byte, 256)
	for i := 0; i < 256; i++ { //Stop at 255 to avoid reading infinitely
		strbytes[i] = m.ReadByte(addr + uint32(i))
		if strbytes[i] == 0 { //Null terminator, string is fully read
			break
		}
	}
	strbytes = bytes.Trim(strbytes, "\x00")
	return string(strbytes)
}
