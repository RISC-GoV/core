package core

type RISCVInstruction int

const (
	LUI RISCVInstruction = iota
	AUIPC
	JAL
	BEQ
	BNE
	BLT
	BGE
	BLTU
	BGEU
	JALR
	LB
	LH
	LW
	LBU
	LHU
	SB
	SH
	SW
	ADDI
	SLTI
	SLTIU
	XORI
	ORI
	ANDI
	SLLI
	SRLI
	SRAI
	EBREAK
	ECALL
	CALL
	ADD
	SUB
	SLL
	SLT
	SLTU
	XOR
	SRL
	SRA
	OR
	AND
	NOP
)

type instruction struct {
	opcode   byte
	operand0 uint32
	operand1 uint32
	operand2 uint32
}

func DecodeInstruction(inst uint32) instruction {
	var opcode OpCode = OpCode((inst >> 24) & 0xFF)
	switch OpToType[opcode] {
	case B:
		return decodeBType(inst, opcode)
	case I:
		return decodeIType(inst, opcode)
	case J:
		return decodeJType(inst, opcode)
	case R:
		return decodeRType(inst, opcode)
	case S:
		return decodeSType(inst, opcode)
	case U:
		return decodeUType(inst, opcode)
	default:
		return instruction{
			opcode:   0,
			operand0: 0,
			operand1: 0,
			operand2: 0,
		}
	}
}

func decodeBType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}

func decodeIType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}

func decodeJType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}

func decodeRType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}

func decodeSType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}

func decodeUType(inst uint32, code OpCode) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}
