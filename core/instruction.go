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
	opcode   uint32
	operand0 uint32
	operand1 uint32
	operand2 uint32
}

// TODO: Implement this function
func DecodeInstruction(inst uint32) instruction {
	return instruction{
		opcode:   0,
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}
}
