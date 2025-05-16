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
			case "print":
				c.PrintRegisters()
				onlyprint = true
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
		state = c.ExecuteSingle()
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
	fmt.Printf("0x%08x: 0x%08x    #%s %s, %s, %s\n", addr, c.Memory.ReadWord(addr), RISCVInstructionToString(instruction.value), RegisterToString(instruction.operand0), RegisterToString(instruction.operand1), RegisterToString(instruction.operand2))
}

func (c *CPU) ExecuteFile(path string) error {
	err := c.LoadFile(path)
	if err != nil {
		return err
	}
	state := OK
	for state == OK || state == E_BREAK { //ignore breakpoints in normal execution mode
		state = c.ExecuteSingle()
	}
	return nil
}

func (c *CPU) FetchNextInstruction() (Instruction, error) {
	return DecodeInstruction(c.Memory.ReadWord(c.PC))
}

func (c *CPU) FetchInstruction(addr uint32) (Instruction, error) {
	return DecodeInstruction(c.Memory.ReadWord(addr))
}

// ExecuteSingle decodes and executes a single instruction from memory at the program counter (PC).
// It updates CPU registers based on the decoded instruction and increments the PC where applicable.
// Returns True if Execution should be stopped
func (c *CPU) ExecuteSingle() int {
	instruction, err := c.FetchNextInstruction()
	if err != nil {
		return -1
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
		if c.ReadRegister(instruction.operand0) == c.ReadRegister(instruction.operand1) {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BNE:
		if c.ReadRegister(instruction.operand0) != c.ReadRegister(instruction.operand1) {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BLT:
		if c.ReadRegister(instruction.operand0) < c.ReadRegister(instruction.operand1) {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BGE:
		if c.ReadRegister(instruction.operand0) >= c.ReadRegister(instruction.operand1) {
			c.PC = uint32(int32(c.PC) - 4 + int32(instruction.operand2))
		}
	case BLTU:
		if c.ReadRegister(instruction.operand0) < c.ReadRegister(instruction.operand1) {
			c.PC += uint32(int32(instruction.operand2) - 4)
		}
	case BGEU:
		if c.ReadRegister(instruction.operand0) >= c.ReadRegister(instruction.operand1) {
			c.PC += uint32(int32(instruction.operand2) - 4)
		}
	case JALR:
		c.WriteRegister(instruction.operand0, c.PC+4)
		c.PC = uint32(int32(c.ReadRegister(instruction.operand1)) + int32(instruction.operand2))
	case LB:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadSingleByte(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
	case LH:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadHalfWord(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
	case LW:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadWord(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
	case LBU:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadSingleByte(
				c.ReadRegister(instruction.operand2)+instruction.operand1))
	case LHU:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadHalfWord(
				c.ReadRegister(instruction.operand2)+instruction.operand1))
	case SB:
		c.Memory.WriteSingleByte(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
	case SH:
		c.Memory.WriteHalfWord(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
	case SW:
		c.Memory.WriteWord(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
	case ADDI:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))+int32(instruction.operand2)))
	case SLTI:
		var b uint32 = 0
		if int32(c.ReadRegister(instruction.operand1)) < int32(instruction.operand2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case SLTIU:
		var b uint32 = 0
		if c.ReadRegister(instruction.operand1) < instruction.operand2 {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case XORI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)^instruction.operand2)
	case ORI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)|instruction.operand2)
		//and the blind forest
	case ANDI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)&instruction.operand2)
	case SLLI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)<<instruction.operand2)
	case SRLI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)>>instruction.operand2)
	case SRAI:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))>>int32(instruction.operand2)))
	case EBREAK:
		return 1
	case ECALL:
		return c.HandleECALL()
	case ADD:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)+c.ReadRegister(instruction.operand2))
	case SUB:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)-c.ReadRegister(instruction.operand2))
	case SLL:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)<<c.ReadRegister(instruction.operand2))
	case SLT:
		var b uint32 = 0
		if int32(c.ReadRegister(instruction.operand1)) < int32(c.ReadRegister(instruction.operand2)) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case SLTU:
		var b uint32 = 0
		if c.ReadRegister(instruction.operand1) < c.ReadRegister(instruction.operand2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
	case XOR:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)^c.ReadRegister(instruction.operand2))
	case SRL:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)>>c.ReadRegister(instruction.operand2))
	case SRA:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))>>int32(c.ReadRegister(instruction.operand2))))
	case OR:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)|c.ReadRegister(instruction.operand2))
	case AND:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)&c.ReadRegister(instruction.operand2))
	case NOP:
		fallthrough
	default:
	}
	return OK
}

func (c *CPU) ReadRegister(reg uint32) uint32 {
	if reg == 0 {
		return 0
	}
	return c.Registers[reg]
}

func (c *CPU) WriteRegister(reg, val uint32) {
	if reg == 0 {
		return
	}
	c.Registers[reg] = val
}
