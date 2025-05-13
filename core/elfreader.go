package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type ELFFile struct {
	Magic      [4]byte // ELF Magic number
	Class      byte    // Class (32-bit or 64-bit)
	Data       byte    // Data encoding (little-endian or big-endian)
	Version    byte    // Version
	OSABI      byte    // OS/ABI
	ABIVersion byte    // ABI version
	Padding    [7]byte // Padding
	Type       uint16  // ELF Type
	Machine    uint16  // Machine type
	Version2   uint32  // Version
	Entry      uint32  // Entry point address
	PhOff      uint32  // Program header offset
	ShOff      uint32  // Section header offset
	Flags      uint32  // Processor-specific flags
	Ehsize     uint16  // ELF header size
	Phentsize  uint16  // Size of program header
	Phnum      uint16  // Number of program headers
	Shentsize  uint16  // Size of section header
	Shnum      uint16  // Number of section headers
	Shstrndx   uint16  // Section header string table index

	// Program headers and section headers
	ProgramHeaders []ProgramHeader

	// Machine code (raw data) from sections
	MachineCode [][]byte
}

// ProgramHeader represents a program header in the ELF file.
type ProgramHeader struct {
	Type     uint32
	Flags    uint32
	Offset   uint32
	VAddr    uint32
	PAddr    uint32
	FileSize uint32
	MemSize  uint32
	Align    uint32
}

func ReadProgramHeader(bytes []byte) *ProgramHeader {
	ph := &ProgramHeader{}
	ph.Type = binary.LittleEndian.Uint32(bytes[0:4])
	ph.Flags = binary.LittleEndian.Uint32(bytes[4:8])
	ph.Offset = binary.LittleEndian.Uint32(bytes[8:12])
	ph.VAddr = binary.LittleEndian.Uint32(bytes[12:16])
	ph.PAddr = binary.LittleEndian.Uint32(bytes[16:20])
	ph.FileSize = binary.LittleEndian.Uint32(bytes[20:24])
	ph.MemSize = binary.LittleEndian.Uint32(bytes[24:28])
	ph.Align = binary.LittleEndian.Uint32(bytes[28:32])
	return ph
}

func ReadELFFile(filePath string) (*ELFFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	var elf ELFFile
	var offset int32 = 0
	var buffer []byte = make([]byte, 1048576) // 1MB buffer
	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if string(buffer[0:3]) != "\x7fELF" { // Wrong file, exit here
		return nil, errors.New("invalid ELF magic number")
	}
	copy(elf.Magic[:], buffer[offset:offset+3])
	offset += 3
	elf.Class = buffer[offset]
	offset++
	if elf.Class != 1 {
		return nil, fmt.Errorf("unsupported ELF class: %d", elf.Class)
	}

	elf.Data = buffer[offset]
	offset++

	var byteOrder binary.ByteOrder
	if elf.Data == 1 { // little-endian
		byteOrder = binary.LittleEndian
	} else if elf.Data == 2 { // big-endian
		byteOrder = binary.BigEndian
	} else {
		return nil, fmt.Errorf("unsupported ELF data encoding: %d", elf.Data)
	}

	elf.Version = buffer[offset]
	offset++
	if elf.Version != 1 {
		return nil, fmt.Errorf("unsupported ELF version: %d", elf.Version)
	}

	elf.OSABI = buffer[offset]
	offset++
	if elf.OSABI != 0 {
		return nil, fmt.Errorf("OSABI is different from 0x03/Linux: %d", elf.OSABI)
	}
	elf.ABIVersion = buffer[offset]
	offset += 8 // 1 for ABI Version and 7 for padding
	elf.Padding = [7]byte{}
	elf.Type = byteOrder.Uint16(buffer[offset : offset+2])
	offset += 2
	if elf.Type != 2 {
		return nil, fmt.Errorf("unsupported ELF type: %d", elf.Type)
	}

	elf.Machine = byteOrder.Uint16(buffer[offset : offset+2])
	offset += 2
	if elf.Machine != 0xf3 {
		return nil, fmt.Errorf("sachine is different from 0xF3/RISC-V: %d", elf.Machine)
	}
	elf.Version2 = byteOrder.Uint32(buffer[offset : offset+4])
	offset += 4
	elf.Entry = byteOrder.Uint32(buffer[offset : offset+4])
	offset += 4
	elf.PhOff = byteOrder.Uint32(buffer[offset : offset+4])
	offset += 4
	elf.ShOff = byteOrder.Uint32(buffer[offset : offset+4])
	offset += 4
	elf.Flags = byteOrder.Uint32(buffer[offset : offset+4])
	offset += 4
	elf.Ehsize = byteOrder.Uint16(buffer[offset : offset+2])
	offset += 2
	elf.Phentsize = byteOrder.Uint16(buffer[offset : offset+2])
	offset += 2
	elf.Phnum = byteOrder.Uint16(buffer[offset : offset+2])
	offset += 8 // We ignore Sections as they are not relevant for execution

	elf.ProgramHeaders = make([]ProgramHeader, elf.Phnum)
	for i := 0; i < int(elf.Phnum); i++ {
		phbuffer := make([]byte, elf.Phentsize)
		copy(phbuffer, buffer[offset:])
		offset += int32(elf.Phentsize)
		elf.ProgramHeaders[i] = *ReadProgramHeader(phbuffer)
	}

	elf.MachineCode = make([][]byte, elf.Phnum)
	for i := 0; i < int(elf.Phnum); i++ {
		code := make([]byte, elf.ProgramHeaders[i].FileSize)
		copy(code, buffer[elf.ProgramHeaders[i].Offset:])
		elf.MachineCode[i] = code
	}
	return &elf, nil
}

func (ELFFile *ELFFile) CopyToMemory(mem *Memory) error {
	for i, ph := range ELFFile.ProgramHeaders {
		if ph.Type != 1 { // PT_LOAD
			continue
		}
		copy(mem.mem[ph.VAddr:], ELFFile.MachineCode[i])
	}

	return nil
}
