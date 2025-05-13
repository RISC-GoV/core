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

func (c *CPU) FetchInstruction(addr uint32) instruction {
	return DecodeInstruction(c.Memory.ReadWord(addr))
}

// ExecuteSingle decodes and executes a single instruction from memory at the program counter (PC).
// It updates CPU registers based on the decoded instruction and increments the PC where applicable.
// Returns True if Execution should be stopped
func (c *CPU) ExecuteSingle() bool {
	instruction := DecodeInstruction(c.Memory.ReadWord(c.PC))
	c.PC += 4
	switch instruction.value {
	case ADDI:
		c.Registers[instruction.operand0] = c.Registers[instruction.operand1] + instruction.operand2
	case LUI:
	case AUIPC:
	case JAL:
	case BEQ:
	case BNE:
	case BLT:
	case BGE:
	case BLTU:
	case BGEU:
	case JALR:
	case LB:
	case LH:
	case LW:
	case LBU:
	case LHU:
	case SB:
	case SH:
	case SW:
	case SLTI:
	case SLTIU:
	case XORI:
	case ORI:
	case ANDI:
	case SLLI:
	case SRLI:
	case SRAI:
	case EBREAK:
		return true
	case ECALL:
	case CALL:
	case ADD:
	case SUB:
	case SLL:
	case SLT:
	case SLTU:
	case XOR:
	case SRL:
	case SRA:
	case OR:
	case AND:
	case NOP:
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
