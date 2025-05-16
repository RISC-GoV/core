package core

import "errors"

type RISCVInstruction int

func DecodeInstruction(inst uint32) (Instruction, error) {
	var opcode OpCode = OpCode(inst & 0b1111111)
	inst >>= 7
	switch OpToType[opcode] {
	case B:
		return decodeBType(inst)
	case I:
		return decodeIType(inst, opcode)
	case J:
		return decodeJType(inst)
	case R:
		return decodeRType(inst, opcode)
	case S:
		return decodeSType(inst)
	case U:
		return decodeUType(inst, opcode)
	default:
		return Instruction{}, errors.New("unknown opcode")
	}
}

func decodeBType(inst uint32) (Instruction, error) {
	var value RISCVInstruction
	tmp := (inst >> 5) & 0x7
	switch tmp {
	case 0x0:
		value = BEQ

	case 0x1:
		value = BNE

	case 0x4:
		value = BLT

	case 0x5:
		value = BGE

	case 0x6:
		value = BLTU

	case 0x7:
		value = BGEU

	default:
		return Instruction{
			operand0: 0,
			operand1: 0,
			operand2: 0,
		}, errors.New("unkown function")
	}

	var twelfth uint32 = (inst & 0x01000000) >> 13
	var eleventh uint32 = (inst & 0x00000001) << 10
	var leftpart uint32 = (inst & 0x00FC0000) >> 14
	var rightpart uint32 = (inst & 0x0000001E) >> 1
	var tmpimm uint32 = twelfth | eleventh | leftpart | rightpart
	var imm = uint32(int32(tmpimm<<20) >> 20) // sign extend trick
	return Instruction{
		value:    value,
		operand0: (inst >> 8) & 0b11111,
		operand1: (inst >> 13) & 0b11111,
		operand2: imm << 1,
	}, nil
}

func decodeIType(inst uint32, code OpCode) (Instruction, error) {
	var value RISCVInstruction
	func3 := (inst >> 5) & 0b111
	func7 := (inst >> 18) & 0b1111111
	switch code {
	case 0b1100111:
		value = JALR

	case 0b0000011:
		switch func3 {
		case 0x0:
			value = LB

		case 0x1:
			value = LH

		case 0x2:
			value = LW

		case 0x4:
			value = LBU

		case 0x5:
			value = LHU

		default:
			return Instruction{}, errors.New("unknown function")
		}

	case 0b0010011:
		switch func3 {
		case 0x0:
			value = ADDI

		case 0x1:
			value = SLLI

		case 0x2:
			value = SLTI

		case 0x3:
			value = SLTIU

		case 0x4:
			value = XORI

		case 0x5:
			switch func7 {
			case 0x00:
				value = SRLI

			case 0x20:
				value = SRAI

			default:
				return Instruction{}, errors.New("unknown function")
			}
		case 0x6:
			value = ORI //and the will of the wisps

		case 0x7:
			value = ANDI
		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0b1110011:
		switch (inst >> 13) & 0xFFF {
		case 0x1:
			value = EBREAK

		case 0x0:
			value = ECALL

		default:
			return Instruction{}, errors.New("unknown function")
		}
	default:
		return Instruction{}, errors.New("unknown opcode")
	}

	var result = Instruction{
		value:    value,
		operand0: inst & 0x1F,
		operand1: (inst >> 8) & 0x1F,
		operand2: uint32(int32(((inst>>13)&0xFFF)<<20) >> 20),
	}

	switch value {
	case SLLI, SRLI, SRAI:
		result.operand1 = result.operand1
		result.operand2 = result.operand2 & 0x1F

	default:
	}

	return result, nil
}

func decodeJType(inst uint32) (Instruction, error) {
	twenty := (inst & 0x1000000) >> 5
	twelve_nineteen := (inst & 0x1FE0) << 6
	eleven := (inst & 0x2000) >> 3
	one_ten := (inst & 0xFFC000) >> 14
	var op1 uint32 = twenty | twelve_nineteen | eleven | one_ten
	op1 = uint32(int32(op1<<12) >> 12) // sign extend trick

	return Instruction{
		value:    JAL,
		operand0: inst & 0b11111,
		operand1: op1 << 1,
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

		case 0x20:
			value = SUB

		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0x1:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLL

	case 0x2:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLT

	case 0x3:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = SLTU

	case 0x4:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = XOR

	case 0x5:
		switch func7 {
		case 0x0:
			value = SRL

		case 0x20:
			value = SRA

		default:
			return Instruction{}, errors.New("unknown function")
		}
	case 0x6:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = OR

	case 0x7:
		if func7 != 0 {
			return Instruction{}, errors.New("unknown function")
		}
		value = AND

	}

	return Instruction{
		value:    value,
		operand0: inst & 0b11111,
		operand1: (inst >> 8) & 0b11111,
		operand2: (inst >> 13) & 0b11111,
	}, nil
}

func decodeSType(inst uint32) (Instruction, error) {
	var value RISCVInstruction
	switch (inst >> 5) & 0b111 {
	case 0x0:
		value = SB

	case 0x1:
		value = SH

	case 0x2:
		value = SW

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
