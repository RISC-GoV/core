package core

import (
	"fmt"
)

const (
	RETURN_ADDRESS  = 1
	STACK_POINTER   = 2
	GLOBAL_POINTER  = 3
	THREAD_POINTER  = 4
	TEMPORARY_ZERO  = 5
	TEMPORARY_ONE   = 6
	TEMPORARY_TWO   = 7
	FRAME_POINTER   = 8
	SAVED_ZERO      = 8
	SAVED_ONE       = 9
	RETURN_ZERO     = 10
	RETURN_ONE      = 11
	ARG_ZERO        = 10
	ARG_ONE         = 11
	ARG_TWO         = 12
	ARG_THREE       = 13
	ARG_FOUR        = 14
	ARG_FIVE        = 15
	ARG_SIX         = 16
	ARG_SEVEN       = 17
	SAVED_TWO       = 18
	SAVED_THREE     = 19
	SAVED_FOUR      = 20
	SAVED_FIVE      = 21
	SAVED_SIX       = 22
	SAVED_SEVEN     = 23
	SAVED_EIGHT     = 24
	SAVED_NINE      = 25
	SAVED_TEN       = 26
	SAVED_ELEVEN    = 27
	TEMPORARY_THREE = 28
	TEMPORARY_FOUR  = 29
	TEMPORARY_FIVE  = 30
	TEMPORARY_SIX   = 31
)

const (
	OK                   = 0
	E_BREAK              = 1
	PROGRAM_EXIT         = 2
	PROGRAM_EXIT_FAILURE = -1
	UNKNOWN_INSTRUCTION  = -2
	IO_ERROR             = -3
)

type CPU struct {
	Memory    *Memory
	Registers [32]uint32
	PC        uint32
}

func NewCPU(mem *Memory) *CPU {
	return &CPU{
		Memory: mem,
	}
}

func (c *CPU) LoadFile(path string) error {
	elf, err := ReadELFFile(path)
	if err != nil {
		return err
	}
	err = elf.CopyToMemory(c.Memory)
	if err != nil {
		return err
	}
	c.PC = elf.Entry
	return nil
}

func (c *CPU) DebugFile(path string) error {
	err := c.LoadFile(path)
	if err != nil {
		return err
	}
	stepping := true
	onlyprint := false
	state := OK
	for {
		if onlyprint {
			onlyprint = false
			continue
		}
		if state != OK || stepping {
			var command string
			switch fmt.Scan(&command); command {
			case "p":
				fallthrough
			case "printRegisters":
				c.PrintRegisters()
				onlyprint = true
				continue
			case "e":
				fallthrough
			case "exit":
				return nil
			case "c":
				fallthrough
			case "continue":
				stepping = false
			case "s":
				fallthrough
			case "step":
				fallthrough
			default:
				stepping = true
			}
		}
		if state == E_BREAK {
			fmt.Printf("Breakpoint hit at 0x%08x\n", c.PC)
		}
		if state == PROGRAM_EXIT {
			fmt.Println("Program exited normally")
			return nil
		}
		if state == PROGRAM_EXIT_FAILURE {
			fmt.Println("Program exited with failure")
			return nil
		}
		state, err = c.ExecuteSingle()
		if err != nil {
			return err
		}
	}
}

func (c *CPU) PrintRegisters() {
	fmt.Println("Registers:")
	for i := 0; i < 32; i++ {
		fmt.Printf("x%-2d: 0x%08x  ", i, c.Registers[i])
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println("\nCurrent Instruction and Context:")
	if c.PC > 4 {
		fmt.Printf("Previous: ")
		c.PrintInstruction(c.PC - 4)
	}
	fmt.Printf("Current:  ")
	c.PrintInstruction(c.PC)
	fmt.Printf("Next 1:   ")
	c.PrintInstruction(c.PC + 4)
}

func (c *CPU) PrintInstruction(addr uint32) {
	instruction, err := c.FetchInstruction(addr)
	if err != nil {
		fmt.Println("Error fetching instruction:", err)
		return
	}
	val, err := c.Memory.ReadWord(addr)
	if err != nil {
		fmt.Println("Error reading instruction :", err)
	}
	fmt.Printf("0x%08x: 0x%08x    #%s %s, %s, %s\n", addr, val, RISCVInstructionToString(instruction.value), RegisterToString(instruction.operand0), RegisterToString(instruction.operand1), RegisterToString(instruction.operand2))
}

func (c *CPU) ExecuteFile(path string) error {
	err := c.LoadFile(path)
	if err != nil {
		return err
	}
	state := OK
	for state == OK || state == E_BREAK { //ignore breakpoints in normal execution mode
		state, err = c.ExecuteSingle()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CPU) FetchNextInstruction() (Instruction, error) {
	val, err := c.Memory.ReadWord(c.PC)
	if err != nil {
		return Instruction{}, err
	}
	return DecodeInstruction(val)
}

func (c *CPU) FetchInstruction(addr uint32) (Instruction, error) {
	val, err := c.Memory.ReadWord(addr)
	if err != nil {
		return Instruction{}, err
	}
	return DecodeInstruction(val)
}

// ExecuteSingle decodes and executes a single instruction from memory at the program counter (PC).
// It updates CPU registers based on the decoded instruction and increments the PC where applicable.
// Returns True if Execution should be stopped
func (c *CPU) ExecuteSingle() (int, error) {
	instruction, err := c.FetchNextInstruction()
	if err != nil {
		return -1, err
	}
	c.PC += 4
	switch instruction.value {
	case LUI:
		c.WriteRegister(instruction.operand0, instruction.operand1<<12)
	case AUIPC:
		c.WriteRegister(instruction.operand0, uint32(int32(c.PC)+int32(instruction.operand1)))
	case JAL:
		c.WriteRegister(instruction.operand0, c.PC)
		c.PC = uint32((int32(c.PC) - 4 + int32(instruction.operand1)))
	case BEQ:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 == val1 {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BNE:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 != val1 {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BLT:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 < val1 {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BGE:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 >= val1 {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BLTU:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 < val1 {
			c.PC += uint32(int32(instruction.operand2) - 4)
		}
	case BGEU:
		val0, err0 := c.ReadRegister(instruction.operand0)
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		if val0 >= val1 {
			c.PC += uint32(int32(instruction.operand2) - 4)
		}
	case JALR:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, c.PC+4)
		c.PC = uint32(int32(val1) + int32(instruction.operand2))
	case LB:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		retVal, err := c.Memory.ReadSingleByte(uint32(int32(val1) + int32(instruction.operand2)*2))
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		c.WriteRegister(instruction.operand0, retVal)
	case LH:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		retVal, err := c.Memory.ReadHalfWord(uint32(int32(val1) + int32(instruction.operand2)*2))
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		c.WriteRegister(instruction.operand0, retVal)
	case LW:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		retVal, err := c.Memory.ReadWord(uint32(int32(val1) + int32(instruction.operand2)*2))
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		c.WriteRegister(instruction.operand0, retVal)
	case LBU:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		retVal, err := c.Memory.ReadSingleByte(val1 + instruction.operand2*2)
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		c.WriteRegister(instruction.operand0, retVal)
	case LHU:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		retVal, err := c.Memory.ReadHalfWord(val1 + instruction.operand2*2)
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
		c.WriteRegister(instruction.operand0, retVal)
	case SB:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val0, err0 := c.ReadRegister(instruction.operand0)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		err := c.Memory.WriteSingleByte(uint32(int32(val1)+int32(instruction.operand2)*2), val0)
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
	case SH:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val0, err0 := c.ReadRegister(instruction.operand0)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		err := c.Memory.WriteHalfWord(uint32(int32(val1)+int32(instruction.operand2)*2), val0)
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
	case SW:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val0, err0 := c.ReadRegister(instruction.operand0)
		if err0 != nil || err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err0, err1)
		}
		err := c.Memory.WriteWord(uint32(int32(val1)+int32(instruction.operand2)*2), val0)
		if err != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err.Error())
		}
	case ADDI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, uint32(int32(val1)+int32(instruction.operand2)))
	case SLTI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		var b uint32 = 0
		if int32(val1) < int32(instruction.operand2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case SLTIU:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		var b uint32 = 0
		if val1 < instruction.operand2 {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case XORI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, val1^instruction.operand2)
	case ORI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, val1|instruction.operand2)
	case ANDI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, val1&instruction.operand2)
	case SLLI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, val1<<instruction.operand2)
	case SRLI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, val1>>instruction.operand2)
	case SRAI:
		val1, err1 := c.ReadRegister(instruction.operand1)
		if err1 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s", c.PC-4, err1)
		}
		c.WriteRegister(instruction.operand0, uint32(int32(val1)>>int32(instruction.operand2)))
	case EBREAK:
		return 1, nil
	case ECALL:
		return c.HandleECALL()
	case ADD:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1+val2)
	case SUB:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1-val2)
	case SLL:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1<<val2)
	case SLT:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		var b uint32 = 0
		if int32(val1) < int32(val2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case SLTU:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		var b uint32 = 0
		if val1 < val2 {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case XOR:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1^val2)
	case SRL:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1>>val2)
	case SRA:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, uint32(int32(val1)>>int32(val2)))
	case OR:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1|val2)
	case AND:
		val1, err1 := c.ReadRegister(instruction.operand1)
		val2, err2 := c.ReadRegister(instruction.operand2)
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("crash at PC=%d with error:\n%s%s", c.PC-4, err1, err2)
		}
		c.WriteRegister(instruction.operand0, val1&val2)
	case NOP:
		fallthrough
	default:
	}
	return OK, nil
}

func (c *CPU) ReadRegister(reg uint32) (uint32, error) {
	if reg == 0 {
		return 0, nil
	}
	if uint32(len(c.Registers)) < reg {
		return 0, fmt.Errorf("read register out of range at PC=%d", c.PC)
	}
	return c.Registers[reg], nil
}

func (c *CPU) WriteRegister(reg, val uint32) {
	if reg == 0 {
		return
	}
	if uint32(len(c.Registers)) < reg {
		return
	}
	c.Registers[reg] = val
}
