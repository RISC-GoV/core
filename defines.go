package core

import (
	"fmt"
)

type OpCode byte
type OpType int

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
	ORI //and the blind forest
	ANDI
	SLLI
	SRLI
	SRAI
	EBREAK
	ECALL
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

type Instruction struct {
	value    RISCVInstruction
	operand0 uint32
	operand1 uint32
	operand2 uint32
}

const (
	R OpType = iota
	I
	S
	B
	U
	J
)

var OpToType = map[OpCode]OpType{
	0b0110111: U,
	0b0010111: U,
	0b1100011: B,
	0b1100111: I,
	0b0000011: I,
	0b0010011: I,
	0b1110011: I,
	0b0100011: S,
	0b0110011: R,
	0b0000000: R,
	0b1101111: J,
}

func RISCVInstructionToString(instruction RISCVInstruction) string {
	switch instruction {
	case LUI:
		return "LUI"
	case AUIPC:
		return "AUIPC"
	case JAL:
		return "JAL"
	case BEQ:
		return "BEQ"
	case BNE:
		return "BNE"
	case BLT:
		return "BLT"
	case BGE:
		return "BGE"
	case BLTU:
		return "BLTU"
	case BGEU:
		return "BGEU"
	case JALR:
		return "JALR"
	case LB:
		return "LB"
	case LH:
		return "LH"
	case LW:
		return "LW"
	case LBU:
		return "LBU"
	case LHU:
		return "LHU"
	case SB:
		return "SB"
	case SH:
		return "SH"
	case SW:
		return "SW"
	case ADDI:
		return "ADDI"
	case SLTI:
		return "SLTI"
	case SLTIU:
		return "SLTIU"
	case XORI:
		return "XORI"
	case ORI:
		return "ORI" //and the blind forest
	case ANDI:
		return "ANDI"
	case SLLI:
		return "SLLI"
	case SRLI:
		return "SRLI"
	case SRAI:
		return "SRAI"
	case EBREAK:
		return "EBREAK"
	case ECALL:
		return "ECALL"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case SLL:
		return "SLL"
	case SLT:
		return "SLT"
	case SLTU:
		return "SLTU"
	case XOR:
		return "XOR"
	case SRL:
		return "SRL"
	case SRA:
		return "SRA"
	case OR:
		return "OR"
	case AND:
		return "AND"
	default:
	}
	return "Unknown instruction"
}

func RegisterToString(register uint32) string {
	switch register {
	case 0:
		return "x0"
	case 1:
		return "ra"
	case 2:
		return "sp"
	case 3:
		return "gp"
	case 4:
		return "tp"
	case 5:
		return "t0"
	case 6:
		return "t1"
	case 7:
		return "t2"
	case 8:
		return "s0"
	case 9:
		return "s1"
	case 10:
		return "a0"
	case 11:
		return "a1"
	case 12:
		return "a2"
	case 13:
		return "a3"
	case 14:
		return "a4"
	case 15:
		return "a5"
	case 16:
		return "a6"
	case 17:
		return "a7"
	case 18:
		return "s2"
	case 19:
		return "s3"
	case 20:
		return "s4"
	case 21:
		return "s5"
	case 22:
		return "s6"
	case 23:
		return "s7"
	case 24:
		return "s8"
	case 25:
		return "s9"
	case 26:
		return "s10"
	case 27:
		return "s11"
	case 28:
		return "t3"
	case 29:
		return "t4"
	case 30:
		return "t5"
	case 31:
		return "t6"
	default:
		return fmt.Sprintf("0x%08x", register)
	}
}
