package core

import "errors"

type RISCVInstruction int

func DecodeInstruction(inst uint32) instruction {
	var opcode OpCode = OpCode(inst & 0b1111111)
	inst >>= 7
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
			operand0: 0,
			operand1: 0,
			operand2: 0,
		}
	}
}

func decodeBType(inst uint32, code OpCode) (instruction, error) {
	var value RISCVInstruction
	switch (inst >> 5) & 0xFFF {
	case 0x0:
		value = BEQ
		break
	case 0x1:
		value = BNE
		break
	case 0x4:
		value = BLT
		break
	case 0x5:
		value = BGE
		break
	case 0x6:
		value = BLTU
		break
	case 0x7:
		value = BGEU
		break
	default:
		return instruction{
			operand0: 0,
			operand1: 0,
			operand2: 0,
		}, errors.New("unkown function")
	}

	var imm uint32 = (inst >> 24) & 0b1
	imm <<= 1
	imm |= inst & 0b1
	imm <<= 1
	imm |= (inst >> 18) & 0b111111
	imm <<= 6
	imm |= (inst >> 1) & 0b1111
	imm <<= 1

	return instruction{
		value:    value,
		operand0: (inst >> 8) & 0x11111,
		operand1: (inst >> 13) & 0x11111,
		operand2: imm,
	}, nil
}

func decodeIType(inst uint32, code OpCode) (instruction, error) { //todo
	return instruction{
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}, errors.New("todo")
}

func decodeJType(inst uint32, code OpCode) (instruction, error) {
	return instruction{
		value:    JAL,
		operand0: inst & 0b11111,
		operand1: (inst >> 5) & 0b11111111111111111111,
		operand2: 0,
	}, nil
}

func decodeRType(inst uint32, code OpCode) (instruction, error) { //todo
	return instruction{
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}, errors.New("todo")
}

func decodeSType(inst uint32, code OpCode) (instruction, error) { //todo
	return instruction{
		operand0: 0,
		operand1: 0,
		operand2: 0,
	}, errors.New("todo")
}

func decodeUType(inst uint32, code OpCode) (instruction, error) {
	var value RISCVInstruction
	if code == 0b0110111 {
		value = LUI
	} else if code == 0b0010111 {
		value = AUIPC
	} else {
		return instruction{
			operand0: 0,
			operand1: 0,
			operand2: 0,
		}, errors.New("unknown opcode")
	}
	return instruction{
		value:    value,
		operand0: inst & 0b11111,
		operand1: (inst >> 5) & 0b11111111111111111111,
		operand2: 0,
	}, nil
}
