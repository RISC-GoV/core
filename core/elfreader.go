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
	Entry      uint64  // Entry point address
	PhOff      uint64  // Program header offset
	ShOff      uint64  // Section header offset
	Flags      uint32  // Processor-specific flags
	Ehsize     uint16  // ELF header size
	Phentsize  uint16  // Size of program header
	Phnum      uint16  // Number of program headers
	Shentsize  uint16  // Size of section header
	Shnum      uint16  // Number of section headers
	Shstrndx   uint16  // Section header string table index

	// Program headers and section headers
	ProgramHeaders []ProgramHeader
	SectionHeaders []SectionHeader
	SectionNames   []string

	// Machine code (raw data) from sections
	SectionsContent map[string][]byte
}

// ProgramHeader represents a program header in the ELF file.
type ProgramHeader struct {
	Type     uint32
	Flags    uint32
	Offset   uint64
	VAddr    uint64
	PAddr    uint64
	FileSize uint64
	MemSize  uint64
	Align    uint64
}

// SectionHeader represents a section header in the ELF file.
type SectionHeader struct {
	Name      uint32
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Linked    uint32
	Info      uint32
	AddrAlign uint64
	EntSize   uint64
}

func ReadELFFile(filePath string) (*ELFFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var elf ELFFile
	err = binary.Read(file, binary.LittleEndian, &elf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ELF header: %w", err)
	}

	if string(elf.Magic[:]) != "\x7fELF" { // Wrong file, exit here
		return nil, errors.New("invalid ELF magic number")
	}

	var byteOrder binary.ByteOrder
	if elf.Data == 1 { // little-endian
		byteOrder = binary.LittleEndian
	} else if elf.Data == 2 { // big-endian
		byteOrder = binary.BigEndian
	} else {
		return nil, fmt.Errorf("unsupported ELF data encoding: %d", elf.Data)
	}

	_, err = file.Seek(int64(elf.PhOff), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to program headers: %w", err)
	}

	elf.ProgramHeaders = make([]ProgramHeader, elf.Phnum)
	for i := 0; i < int(elf.Phnum); i++ {
		var ph ProgramHeader
		err = binary.Read(file, byteOrder, &ph)
		if err != nil {
			return nil, fmt.Errorf("failed to read program header: %w", err)
		}
		elf.ProgramHeaders[i] = ph
	}

	// Read section headers
	_, err = file.Seek(int64(elf.ShOff), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to section headers: %w", err)
	}

	elf.SectionHeaders = make([]SectionHeader, elf.Shnum)
	for i := 0; i < int(elf.Shnum); i++ {
		var sh SectionHeader
		err = binary.Read(file, byteOrder, &sh)
		if err != nil {
			return nil, fmt.Errorf("failed to read section header: %w", err)
		}
		elf.SectionHeaders[i] = sh
	}

	sectionNamesSection := elf.SectionHeaders[elf.Shstrndx]
	_, err = file.Seek(int64(sectionNamesSection.Offset), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to section names: %w", err)
	}

	elf.SectionNames = make([]string, elf.Shnum)
	for i := 0; i < int(elf.Shnum); i++ {
		var nameOffset uint32
		err = binary.Read(file, byteOrder, &nameOffset)
		if err != nil {
			return nil, fmt.Errorf("failed to read section name offset: %w", err)
		}

		_, err = file.Seek(int64(sectionNamesSection.Offset+uint64(nameOffset)), 0)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to section name: %w", err)
		}

		var sectionName string
		for {
			var b byte
			err = binary.Read(file, byteOrder, &b)
			if err != nil {
				return nil, fmt.Errorf("failed to read section name byte: %w", err)
			}
			if b == 0 {
				break
			}
			sectionName += string(b)
		}

		elf.SectionNames[i] = sectionName
	}

	elf.SectionsContent = make(map[string][]byte)
	for i := 0; i < int(elf.Shnum); i++ {
		section := elf.SectionHeaders[i]
		if section.Size == 0 {
			continue // skip empty sections
		}

		_, err = file.Seek(int64(section.Offset), 0)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to section content: %w", err)
		}

		content := make([]byte, section.Size)
		err = binary.Read(file, byteOrder, &content)
		if err != nil {
			return nil, fmt.Errorf("failed to read section content: %w", err)
		}

		// Store content by section name
		sectionName := elf.SectionNames[i]
		elf.SectionsContent[sectionName] = content
	}

	return &elf, nil
}
