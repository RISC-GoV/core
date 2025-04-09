package core

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

func (c *CPU) FetchNextInstruction() (Instruction, error) {
	return DecodeInstruction(c.Memory.ReadWord(c.PC))
}

func (c *CPU) FetchInstruction(addr uint32) (Instruction, error) {
	return DecodeInstruction(c.Memory.ReadWord(addr))
}

func (c *CPU) ExecuteSingle() error {
	instruction, err := DecodeInstruction(c.Memory.ReadWord(c.PC))
	if err != nil {
		return err
	}
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
	return nil
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
