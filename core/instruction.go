package core

import "errors"

type RISCVInstruction int

func DecodeInstruction(inst uint32) (Instruction, error) {
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
		return Instruction{}, errors.New("unknown opcode")
	}
}

func decodeBType(inst uint32, code OpCode) (Instruction, error) {
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
		return Instruction{
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

	return Instruction{
		value:    value,
		operand0: (inst >> 8) & 0x11111,
		operand1: (inst >> 13) & 0x11111,
		operand2: imm,
	}, nil
}

func decodeIType(inst uint32, code OpCode) (Instruction, error) {
	var value RISCVInstruction
	func3 := (inst >> 5) & 0b111
	func7 := (inst >> 18) & 0b1111111
	switch code {
	case 0b1100111:
		value = JALR
		break
	case 0b0000011:
		switch func3 {
		case 0x0:
			value = LB
			break
		case 0x1:
			value = LH
			break
		case 0x2:
			value = LW
			break
		case 0x4:
			value = LBU
			break
		case 0x5:
			value = LHU
			break
		default:
			return Instruction{}, errors.New("unknown function")
		}
		break
	case 0b0010011:
		switch func3 {
		case 0x0:
			value = ADDI
			break
		case 0x1:
			value = SLLI
			break
		case 0x2:
			value = SLTI
			break
		case 0x3:
			value = SLTIU
			break
		case 0x4:
			value = XORI
			break
		case 0x5:
			switch func7 {
			case 0x00:
				value = SRLI
				break
			case 0x32:
				value = SRAI
				break
			default:
				return Instruction{}, errors.New("unknown function")
			}
		case 0x6:
			value = ORI //and the will of the wisps
			break
		case 0x7:
			value = ANDI
		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0b1110011:
		switch func7 {
		case 0x1:
			value = EBREAK
			break
		case 0x0:
			value = ECALL
			break
		default:
			return Instruction{}, errors.New("unknown function")
		}
	default:
		return Instruction{}, errors.New("unknown opcode")
	}

	var result = Instruction{
		value:    value,
		operand0: inst & 0b11111,
		operand1: (inst >> 13) & 0b1111111111111111,
		operand2: (inst >> 8) & 0b11111,
	}

	switch value {
	case SLLI, SRLI, SRAI:
		tmp := result.operand1
		result.operand1 = result.operand2
		result.operand2 = tmp & 0b11111
		break
	default:
		break
	}

	return result, errors.New("todo")
}

func decodeJType(inst uint32, code OpCode) (Instruction, error) {
	return Instruction{
		value:    JAL,
		operand0: inst & 0b11111,
		operand1: (inst >> 5) & 0b11111111111111111111,
		operand2: 0,
	}, nil
}

func decodeRType(inst uint32, code OpCode) (Instruction, error) {
	if code == 0 {
		return Instruction{NOP, 0, 0, 0}, nil
	}
	var func7 = (inst >> 18) & 0b1111111
	var value RISCVInstruction
	switch (inst >> 5) & 0b111 {
	case 0x0:
		switch func7 {
		case 0x0:
			value = ADD
			break
		case 0x20:
			value = ADD
			break
		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0x1:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLL
		break
	case 0x2:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLT
		break
	case 0x3:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLTU
		break
	case 0x4:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = XOR
		break
	case 0x5:
		switch func7 {
		case 0x0:
			value = SRL
			break
		case 0x20:
			value = SRA
			break
		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0x6:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = OR
		break
	case 0x7:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = AND
		break
	}

	return Instruction{
		value:    value,
		operand0: inst & 0b11111,
		operand1: (inst >> 8) & 0b11111,
		operand2: (inst >> 13) & 0b11111,
	}, nil
}

func decodeSType(inst uint32, code OpCode) (Instruction, error) {
	var value RISCVInstruction
	switch (inst >> 5) & 0b111 {
	case 0x0:
		value = SB
		break
	case 0x1:
		value = SH
		break
	case 0x2:
		value = SW
		break
	default:
		return Instruction{}, errors.New("unknown function")
	}

	imm := (inst >> 18) & 0b1111111
	imm <<= 5
	imm |= inst & 0b11111

	return Instruction{
		value:    value,
		operand0: (inst >> 13) & 0b11111,
		operand1: imm,
		operand2: (inst >> 8) & 0b11111,
	}, nil
}

func decodeUType(inst uint32, code OpCode) (Instruction, error) {
	var value RISCVInstruction
	if code == 0b0110111 {
		value = LUI
	} else if code == 0b0010111 {
		value = AUIPC
	} else {
		return Instruction{}, errors.New("unknown opcode")
	}
	return Instruction{
		value:    value,
		operand0: inst & 0b11111,
		operand1: (inst >> 5) & 0b11111111111111111111,
		operand2: 0,
	}, nil
}
