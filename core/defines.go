package core

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
}
