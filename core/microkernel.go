package core

import (
	"os"
)

const (
	GETCWD   = 17
	MKDIRAT  = 34
	UNLINKAT = 35
	CHDIR    = 49
	FCHDIR   = 50
	OPENAT   = 56
	CLOSE    = 57
	READ     = 63
	WRITE    = 64
	EXIT     = 93
)

const AT_FDCWD = -100

type MicroKernel struct {
	CWD             string
	FileDescriptors []string
}

var Kernel MicroKernel

func IsValidFileDescriptor(fd int32) bool {
	if fd < 0 || fd >= int32(len(Kernel.FileDescriptors)) || Kernel.FileDescriptors[fd] == "" {
		return false
	}
	return true
}

func GetPath(path string, dirfd int32) string {

	if path[0] != '.' {
		return path
	}
	if dirfd == AT_FDCWD {
		return "./" + Kernel.CWD + path[1:]
	}
	if !IsValidFileDescriptor(dirfd) {
		return ""
	}
	return Kernel.FileDescriptors[dirfd] + path[1:]
}

func (c *CPU) HandleECALL() int {
	switch c.ReadRegister(ARG_SEVEN) { //a7 contains the function we are trying to call
	case GETCWD:
		address := c.ReadRegister(ARG_ZERO) //Pointer to start of Buffer we write to
		cwd := []byte(Kernel.CWD)
		for i := range len(cwd) {
			c.Memory.WriteByte(address+uint32(i), cwd[i])
		}
	case MKDIRAT:
		dirfd := c.ReadRegister(ARG_ZERO)
		address := c.ReadRegister(ARG_ONE) //Pointer to start of string we are reading from
		mode := c.ReadRegister(ARG_TWO)
		path := GetPath(c.Memory.ReadString(address), int32(dirfd))
		err := os.Mkdir(path, os.FileMode(mode))
		if err != nil {
			return IO_ERROR
		}
		Kernel.FileDescriptors = append(Kernel.FileDescriptors, path)
		c.Memory.WriteWord(address, uint32(len(Kernel.FileDescriptors)-1)) //Return the file descriptor as the return value
	case UNLINKAT:
		address := c.ReadRegister(ARG_ONE) //Pointer to start of string we are reading from
		path := GetPath(c.Memory.ReadString(address), int32(c.ReadRegister(ARG_ZERO)))
		err := os.Remove(path)
		if err != nil {
			return IO_ERROR
		}
	case CHDIR:
		address := c.ReadRegister(ARG_ZERO) //Pointer to start of string we are reading from
		path := c.Memory.ReadString(address)
		Kernel.CWD = path
	case FCHDIR:
		fd := int32(c.ReadRegister(ARG_ZERO))
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR
		}
		Kernel.CWD = Kernel.FileDescriptors[fd]
	case OPENAT:
		address := c.ReadRegister(ARG_ONE) //Pointer to start of string we are reading from
		path := GetPath(c.Memory.ReadString(address), int32(c.ReadRegister(ARG_ZERO)))
		flags := c.ReadRegister(ARG_TWO)
		mode := c.ReadRegister(ARG_THREE)
		if flags&0x0100 != 0 { //O_CREAT
			file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(mode))
			if err != nil {
				return IO_ERROR
			}
			file.Close()
		}
		Kernel.FileDescriptors = append(Kernel.FileDescriptors, path)
	case CLOSE:
		fd := int32(c.ReadRegister(ARG_ZERO))
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR
		}
		Kernel.FileDescriptors[fd] = ""
	case READ:
		fd := int32(c.ReadRegister(ARG_ZERO))
		dest := c.ReadRegister(ARG_ONE)
		size := c.ReadRegister(ARG_TWO)
		offset := c.ReadRegister(ARG_THREE)
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR
		}
		file, err := os.OpenFile(Kernel.FileDescriptors[fd], os.O_RDONLY, 0)
		if err != nil {
			return IO_ERROR
		}
		defer file.Close()
		buf := make([]byte, offset+size)
		amt, err := file.Read(buf)
		if uint32(amt) != size+offset || err != nil {
			return IO_ERROR
		}
		for i := 0; uint32(i) < size; i++ {
			c.Memory.WriteByte(dest+uint32(i), buf[uint32(i)+offset])
		}
	case WRITE:
		fd := int32(c.ReadRegister(ARG_ZERO))
		source := c.ReadRegister(ARG_ONE)
		size := c.ReadRegister(ARG_TWO)
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR
		}
		file, err := os.OpenFile(Kernel.FileDescriptors[fd], os.O_WRONLY, 0)
		if err != nil {
			return IO_ERROR
		}
		defer file.Close()
		buf := make([]byte, size)
		for i := 0; uint32(i) < size; i++ {
			buf[i] = c.Memory.ReadByte(source + uint32(i))
		}
		_, err = file.Write(buf)
		if err != nil {
			return IO_ERROR
		}
	case EXIT:
		returnCode := c.ReadRegister(ARG_ZERO)
		if returnCode != 0 {
			return PROGRAM_EXIT_FAILURE
		}
		return PROGRAM_EXIT

	default:
		return 0
	}
	return 0
}
