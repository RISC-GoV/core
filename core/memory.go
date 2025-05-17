package core

import (
	"fmt"
)

type Memory struct {
	mem []byte
}

func NewMemory() *Memory {
	return &Memory{
		mem: make([]byte, 1048576), //1MB
	}
}
func (m *Memory) ReadSingleByte(addr uint32) (uint32, error) {
	if uint32(len(m.mem))-1 < addr {
		return 0, fmt.Errorf("read single byte out of range at addr=%d", addr)
	}
	return uint32(m.mem[addr]), nil
}

func (m *Memory) ReadByte(addr uint32) (byte, error) {
	if uint32(len(m.mem))-1 < addr {
		return 0, fmt.Errorf("read byte out of range at addr=%d", addr)
	}
	return m.mem[addr], nil
}

func (m *Memory) WriteSingleByte(addr uint32, val uint32) error {
	if uint32(len(m.mem))-1 < addr {
		return fmt.Errorf("write single byte out of range at addr=%d", addr)
	}
	m.mem[addr] = byte(val)
	return nil
}

func (m *Memory) WriteByte(addr uint32, val byte) error {
	if uint32(len(m.mem))-1 < addr {
		return fmt.Errorf("write byte out of range at addr=%d", addr)
	}
	m.mem[addr] = val
	return nil
}

func (m *Memory) ReadHalfWord(addr uint32) (uint32, error) {
	if uint32(len(m.mem)) < addr {
		return 0, fmt.Errorf("read half-word out of range at addr=%d", addr)
	}
	return uint32(m.mem[addr]) | uint32(m.mem[addr+1])<<8, nil
}

func (m *Memory) WriteHalfWord(addr uint32, val uint32) error {
	if uint32(len(m.mem)) < addr {
		return fmt.Errorf("write half-word out of range at addr=%d", addr)
	}
	m.mem[addr] = byte(val)
	m.mem[addr+1] = byte(val >> 8)
	return nil
}

func (m *Memory) ReadWord(addr uint32) (uint32, error) {
	if uint32(len(m.mem)) < addr+2 {
		return 0, fmt.Errorf("read word out of range at addr=%d", addr)
	}
	return uint32(m.mem[addr]) | uint32(m.mem[addr+1])<<8 | uint32(m.mem[addr+2])<<16 | uint32(m.mem[addr+3])<<24, nil
}

func (m *Memory) WriteWord(addr uint32, val uint32) error {
	if uint32(len(m.mem)) < addr+2 {
		return fmt.Errorf("write word out of range at addr=%d", addr)
	}
	m.mem[addr] = byte(val)
	m.mem[addr+1] = byte(val >> 8)
	m.mem[addr+2] = byte(val >> 16)
	m.mem[addr+3] = byte(val >> 24)
	return nil
}

func (m *Memory) ReadString(addr uint32) (string, error) {
	maxInRange := 256
	if uint32(len(m.mem)) < addr+255 {
		maxInRange = int(uint32(len(m.mem)) - addr)
		if maxInRange <= 0 {
			return "", fmt.Errorf("read string out of range at addr=%d", addr)
		}
	}

	strbytes := make([]byte, 0, maxInRange)
	for i := 0; i < maxInRange; i++ {
		b, err := m.ReadByte(addr + uint32(i))
		if err != nil {
			return string(strbytes), fmt.Errorf("read string error at index %d: %w", i, err)
		}
		if b == 0 { // Null terminator, string is fully read
			break
		}
		strbytes = append(strbytes, b)
	}

	return string(strbytes), nil
}
