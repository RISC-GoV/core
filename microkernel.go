package core

import (
	"fmt"
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

func (k *MicroKernel) Init() {
	k.CWD = "/"
	k.FileDescriptors = []string{"stdin", "stdout", "stderr"}
}

var Kernel MicroKernel

func IsValidFileDescriptor(fd int32) bool {
	return !(fd < 3 || fd >= int32(len(Kernel.FileDescriptors)) || Kernel.FileDescriptors[fd] == "")
}

func IsSpecialFileDescriptor(fd int32) bool {
	return fd == 0 || fd == 1 || fd == 2
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
		val, err := c.Memory.ReadString(address)
		if err != nil {
			panic(fmt.Sprintf("crash at PC=%d with error:\n%s", c.PC-4, err.Error()))
		}
		path := GetPath(val, int32(dirfd))
		err = os.Mkdir(path, os.FileMode(mode))
		if err != nil {
			return IO_ERROR
		}
		Kernel.FileDescriptors = append(Kernel.FileDescriptors, path)
		c.Memory.WriteWord(address, uint32(len(Kernel.FileDescriptors)-1)) //Return the file descriptor as the return value
	case UNLINKAT:
		address := c.ReadRegister(ARG_ONE) //Pointer to start of string we are reading from
		val, err := c.Memory.ReadString(address)
		if err != nil {
			panic(fmt.Sprintf("crash at PC=%d with error:\n%s", c.PC-4, err.Error()))
		}
		path := GetPath(val, int32(c.ReadRegister(ARG_ZERO)))
		err = os.Remove(path)
		if err != nil {
			return IO_ERROR
		}
	case CHDIR:
		address := c.ReadRegister(ARG_ZERO) //Pointer to start of string we are reading from
		path, err := c.Memory.ReadString(address)
		if err != nil {
			panic(fmt.Sprintf("crash at PC=%d with error:\n%s", c.PC-4, err.Error()))
		}
		Kernel.CWD = path
	case FCHDIR:
		fd := int32(c.ReadRegister(ARG_ZERO))
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR
		}
		Kernel.CWD = Kernel.FileDescriptors[fd]
	case OPENAT:
		address := c.ReadRegister(ARG_ONE) //Pointer to start of string we are reading from
		val, err := c.Memory.ReadString(address)
		if err != nil {
			panic(fmt.Sprintf("crash at PC=%d with error:\n%s", c.PC-4, err.Error()))
		}
		path := GetPath(val, int32(c.ReadRegister(ARG_ZERO)))
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
		c.WriteRegister(ARG_ZERO, uint32(len(Kernel.FileDescriptors)-1)) //Return the file descriptor as the return value
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
		var file *os.File
		var err error
		if IsValidFileDescriptor(fd) {
			file, err = os.OpenFile(Kernel.FileDescriptors[fd], os.O_RDONLY, 0)
			if err != nil {
				return IO_ERROR
			}
			defer file.Close()
		} else {
			if fd != 0 {
				return IO_ERROR
			}
			file = os.Stdin
		}
		buf := make([]byte, offset+size)
		amt, err := file.Read(buf)
		if uint32(amt) != size+offset || err != nil {
			return IO_ERROR
		}
		for i := 0; uint32(i) < size; i++ {
			c.Memory.WriteByte(dest+uint32(i), buf[uint32(i)+offset])
		}
		c.WriteRegister(ARG_ZERO, uint32(amt)) //Return the number of bytes read
	case WRITE:
		fd := int32(c.ReadRegister(ARG_ZERO))
		source := c.ReadRegister(ARG_ONE)
		size := c.ReadRegister(ARG_TWO)
		var file *os.File
		var err error
		if IsValidFileDescriptor(fd) {
			file, err = os.OpenFile(Kernel.FileDescriptors[fd], os.O_RDONLY, 0)
			if err != nil {
				return IO_ERROR
			}
			defer file.Close()
		} else {
			switch fd {
			case 1:
				file = os.Stdout
			case 2:
				file = os.Stderr
			default:
				return IO_ERROR
			}
		}
		buf := make([]byte, size)
		for i := 0; uint32(i) < size; i++ {
			buf[i], err = c.Memory.ReadByte(source + uint32(i))
			if err != nil {
				panic(fmt.Sprintf("crash at PC=%d with error:\n%s", c.PC-4, err.Error()))
			}
		}
		written, err := file.Write(buf)
		if err != nil {
			return IO_ERROR
		}
		c.WriteRegister(ARG_ZERO, uint32(written)) //Return the number of bytes written

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
