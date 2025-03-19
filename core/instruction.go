package core

const (
	ADDI = 0 // TODO:Change this value this is just a placeholder
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
