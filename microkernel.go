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
func (c *CPU) HandleECALL() (int, error) {
	a7, err := c.ReadRegister(ARG_SEVEN) // a7 contains the function we are trying to call
	if err != nil {
		return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
	}
	switch a7 {
	case GETCWD:
		address, err := c.ReadRegister(ARG_ZERO) // Pointer to start of Buffer we write to
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		cwd := []byte(Kernel.CWD)
		for i := range len(cwd) {
			c.Memory.WriteByte(address+uint32(i), cwd[i])
		}
	case MKDIRAT:
		dirfd, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		address, err := c.ReadRegister(ARG_ONE) // Pointer to start of string we are reading from
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		mode, err := c.ReadRegister(ARG_TWO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		val, err := c.Memory.ReadString(address)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		path := GetPath(val, int32(dirfd))
		err = os.Mkdir(path, os.FileMode(mode))
		if err != nil {
			return IO_ERROR, nil
		}
		Kernel.FileDescriptors = append(Kernel.FileDescriptors, path)
		c.Memory.WriteWord(address, uint32(len(Kernel.FileDescriptors)-1)) // Return the file descriptor as the return value
	case UNLINKAT:
		_, err := c.ReadRegister(ARG_ZERO) // Not used directly, but for completeness
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		address, err := c.ReadRegister(ARG_ONE) // Pointer to start of string we are reading from
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		val, err := c.Memory.ReadString(address)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		zero, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		path := GetPath(val, int32(zero))
		err = os.Remove(path)
		if err != nil {
			return IO_ERROR, nil
		}
	case CHDIR:
		address, err := c.ReadRegister(ARG_ZERO) // Pointer to start of string we are reading from
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		path, err := c.Memory.ReadString(address)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		Kernel.CWD = path
	case FCHDIR:
		fdVal, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		fd := int32(fdVal)
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR, nil
		}
		Kernel.CWD = Kernel.FileDescriptors[fd]
	case OPENAT:
		zero, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		address, err := c.ReadRegister(ARG_ONE) // Pointer to start of string we are reading from
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		val, err := c.Memory.ReadString(address)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		path := GetPath(val, int32(zero))
		flags, err := c.ReadRegister(ARG_TWO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		mode, err := c.ReadRegister(ARG_THREE)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		if flags&0x0100 != 0 { // O_CREAT
			file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(mode))
			if err != nil {
				return IO_ERROR, nil
			}
			file.Close()
		}
		Kernel.FileDescriptors = append(Kernel.FileDescriptors, path)
		c.WriteRegister(ARG_ZERO, uint32(len(Kernel.FileDescriptors)-1)) // Return the file descriptor as the return value
	case CLOSE:
		fdVal, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		fd := int32(fdVal)
		if !IsValidFileDescriptor(fd) {
			return IO_ERROR, nil
		}
		Kernel.FileDescriptors[fd] = ""
	case READ:
		fdVal, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		dest, err := c.ReadRegister(ARG_ONE)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		size, err := c.ReadRegister(ARG_TWO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		offset, err := c.ReadRegister(ARG_THREE)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		fd := int32(fdVal)
		var file *os.File
		if IsValidFileDescriptor(fd) {
			file, err = os.OpenFile(Kernel.FileDescriptors[fd], os.O_RDONLY, 0)
			if err != nil {
				return IO_ERROR, nil
			}
			defer file.Close()
		} else {
			if fd != 0 {
				return IO_ERROR, nil
			}
			file = os.Stdin
		}
		buf := make([]byte, offset+size)
		amt, err := file.Read(buf)
		if err != nil {
			return IO_ERROR, nil
		}
		for i := 0; uint32(i) < size; i++ {
			c.Memory.WriteByte(dest+uint32(i), buf[uint32(i)+offset])
		}
		c.WriteRegister(ARG_ZERO, uint32(amt)) // Return the number of bytes read
	case WRITE:
		fdVal, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		source, err := c.ReadRegister(ARG_ONE)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		size, err := c.ReadRegister(ARG_TWO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		fd := int32(fdVal)
		var file *os.File
		if IsValidFileDescriptor(fd) {
			file, err = os.OpenFile(Kernel.FileDescriptors[fd], os.O_RDONLY, 0)
			if err != nil {
				return IO_ERROR, nil
			}
			defer file.Close()
		} else {
			switch fd {
			case 1:
				file = os.Stdout
			case 2:
				file = os.Stderr
			default:
				return IO_ERROR, nil
			}
		}
		buf := make([]byte, size)
		for i := 0; uint32(i) < size; i++ {
			buf[i], err = c.Memory.ReadByte(source + uint32(i))
			if err != nil {
				return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
			}
		}
		written, err := file.Write(buf)
		if err != nil {
			return IO_ERROR, nil
		}
		c.WriteRegister(ARG_ZERO, uint32(written)) // Return the number of bytes written

	case EXIT:
		returnCode, err := c.ReadRegister(ARG_ZERO)
		if err != nil {
			return 0, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		if returnCode != 0 {
			return PROGRAM_EXIT_FAILURE, nil
		}
		return PROGRAM_EXIT, nil

	default:
		return 0, nil
	}
	return 0, nil
}
