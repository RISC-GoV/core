package core

import "fmt"

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
	stop := false
	for !stop {
		stop = c.ExecuteSingle()
		var command string
		switch _, err = fmt.Scan(command); command {
		case "print":
			c.PrintRegisters()
		case "stop":
			stop = true
		default:
		}
	}
	return nil
}

func (c *CPU) PrintRegisters() {
	for i := 0; i < 32; i++ {
		fmt.Printf("x%d: %x\n", i, c.Registers[i])
	}
}

func (c *CPU) ExecuteFile(path string) error {
	err := c.LoadFile(path)
	if err != nil {
		return err
	}
	for !c.ExecuteSingle() {
	}
	return nil
}

func (c *CPU) FetchNextInstruction() instruction {
	return DecodeInstruction(c.Memory.ReadWord(c.PC))
}

func (c *CPU) FetchInstruction(addr uint32) (Instruction, error) {
	return DecodeInstruction(c.Memory.ReadWord(addr))
}

// ExecuteSingle decodes and executes a single instruction from memory at the program counter (PC).
// It updates CPU registers based on the decoded instruction and increments the PC where applicable.
// Returns True if Execution should be stopped
func (c *CPU) ExecuteSingle() bool {
	instruction, err := DecodeInstruction(c.Memory.ReadWord(c.PC))
	if err != nil {
		return false
	}
	c.PC += 4
	switch instruction.value {
	case LUI:
		c.WriteRegister(instruction.operand0, c.Registers[instruction.operand1])
		break
	case AUIPC:
		c.WriteRegister(instruction.operand0, uint32(int32(c.PC)+int32(instruction.operand1)))
		break
	case JAL:
		c.WriteRegister(instruction.operand0, c.PC+4)
		c.PC = uint32(int32(c.PC) + int32(instruction.operand1))
		break
	case BEQ:
		if instruction.operand0 == instruction.operand1 {
			c.PC = uint32(int32(c.PC) + int32(instruction.operand2))
		}
		break
	case BNE:
		if instruction.operand0 != instruction.operand1 {
			c.PC = uint32(int32(c.PC) + int32(instruction.operand2))
		}
		break
	case BLT:
		if instruction.operand0 < instruction.operand1 {
			c.PC = uint32(int32(c.PC) + int32(instruction.operand2))
		}
		break
	case BGE:
		if instruction.operand0 >= instruction.operand1 {
			c.PC = uint32(int32(c.PC) + int32(instruction.operand2))
		}
		break
	case BLTU:
		if instruction.operand0 < instruction.operand1 {
			c.PC += instruction.operand2
		}
		break
	case BGEU:
		if instruction.operand0 >= instruction.operand1 {
			c.PC += instruction.operand2
		}
		break
	case JALR:
		c.WriteRegister(instruction.operand0, c.PC+4)
		c.PC = uint32(int32(c.ReadRegister(instruction.operand2)) + int32(instruction.operand1))
		break
	case LB:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadSingleByte(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
		break
	case LH:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadHalfWord(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
		break
	case LW:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadWord(
				uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1))))
		break
	case LBU:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadSingleByte(
				c.ReadRegister(instruction.operand2)+instruction.operand1))
		break
	case LHU:
		c.WriteRegister(
			instruction.operand0,
			c.Memory.ReadHalfWord(
				c.ReadRegister(instruction.operand2)+instruction.operand1))
		break
	case SB:
		c.Memory.WriteSingleByte(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
		break
	case SH:
		c.Memory.WriteHalfWord(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
		break
	case SW:
		c.Memory.WriteWord(
			uint32(int32(c.ReadRegister(instruction.operand2))+int32(instruction.operand1)),
			c.ReadRegister(instruction.operand0))
		break
	case ADDI:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))+int32(instruction.operand2)))
		break
	case SLTI:
		var b uint32 = 0
		if int32(c.ReadRegister(instruction.operand1)) < int32(instruction.operand2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
		break
	case SLTIU:
		var b uint32 = 0
		if c.ReadRegister(instruction.operand1) < instruction.operand2 {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
		break
	case XORI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)^instruction.operand2)
		break
	case ORI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)|instruction.operand2)
		break //and the blind forest
	case ANDI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)&instruction.operand2)
		break
	case SLLI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)<<instruction.operand2)
		break
	case SRLI:
		c.WriteRegister(instruction.operand0, c.ReadRegister(instruction.operand1)>>instruction.operand2)
		break
	case SRAI:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))>>int32(instruction.operand2)))
		break
	case EBREAK:
		return true
		break //todo
	case ECALL:
		break //todo
	case CALL:
		break //todo ecall
	case ADD:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)+c.ReadRegister(instruction.operand2))
		break
	case SUB:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)-c.ReadRegister(instruction.operand2))
		break
	case SLL:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)<<c.ReadRegister(instruction.operand2))
		break
	case SLT:
		var b uint32 = 0
		if int32(c.ReadRegister(instruction.operand1)) < int32(c.ReadRegister(instruction.operand2)) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
		break
	case SLTU:
		var b uint32 = 0
		if c.ReadRegister(instruction.operand1) < c.ReadRegister(instruction.operand2) {
			b = 1
		}
		c.WriteRegister(instruction.operand0, b)
		break
	case XOR:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)^c.ReadRegister(instruction.operand2))
		break
	case SRL:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)>>c.ReadRegister(instruction.operand2))
		break
	case SRA:
		c.WriteRegister(
			instruction.operand0,
			uint32(int32(c.ReadRegister(instruction.operand1))>>int32(c.ReadRegister(instruction.operand2))))
		break
	case OR:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)|c.ReadRegister(instruction.operand2))
		break
	case AND:
		c.WriteRegister(
			instruction.operand0,
			c.ReadRegister(instruction.operand1)&c.ReadRegister(instruction.operand2))
		break
	case NOP:
		break
	}
	return false
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
