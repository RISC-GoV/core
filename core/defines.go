package core

type OpCode byte
type OpType int

type Instruction struct {
	opType OpCode
	opByte []byte
	name   string
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

var OpTypeToInstruction = map[OpCode]Instruction{

	// U TYPE
	0b0110111: {U, nil, "lui"},
	0b0010111: {U, nil, "auipc"},

	// J TYPE
	0b1101111: {J, nil, "jal"},

	// B TYPE
	0b1100011: {B, []byte{0x0}, "beq"},
	0b1100011: {B, []byte{0x1}, "bne"},
	0b1100011: {B, []byte{0x4}, "blt"},
	0b1100011: {B, []byte{0x5}, "bge"},
	0b1100011: {B, []byte{0x6}, "bltu"},
	0b1100011: {B, []byte{0x7}, "bgeu"},

	// I TYPE: opbyte = opcode, func3, imm[11:5]
	"jalr":   {I, []byte{, 0x0, 0x0}},
	"lb":     {I, []byte{, 0x0, 0x0}},
	"lh":     {I, []byte{, 0x1, 0x0}},
	"lw":     {I, []byte{, 0x2, 0x0}},
	"lbu":    {I, []byte{, 0x4, 0x0}},
	"lhu":    {I, []byte{, 0x5, 0x0}},
	"sb":     {S, []byte{, 0x0, 0x0}},
	"sh":     {S, []byte{, 0x1, 0x0}},
	"sw":     {S, []byte{, 0x2, 0x0}},
	"addi":   {I, []byte{, 0x0, 0x0}},
	"slti":   {I, []byte{, 0x2, 0x0}},
	"sltiu":  {I, []byte{, 0x3, 0x0}},
	"xori":   {I, []byte{, 0x4, 0x0}},
	"ori":    {I, []byte{, 0x6, 0x0}},
	"andi":   {I, []byte{, 0x7, 0x0}},
	"slli":   {I, []byte{, 0x1, 0x00}},
	"srli":   {I, []byte{, 0x5, 0x00}},
	"srai":   {I, []byte{, 0x5, 0x32}},
	"ebreak": {I, []byte{, 0x0, 0x1}}, // Adjust as needed
	"ecall":  {I, []byte{, 0x0, 0x0}}, // Adjust as needed
	"call":   {I, []byte{}},                    // Adjust as needed

	// R TYPE
	"add":  {R, []byte{, 0x0, 0x00}},
	"sub":  {R, []byte{, 0x0, 0x20}},
	"sll":  {R, []byte{, 0x1, 0x00}},
	"slt":  {R, []byte{, 0x2, 0x00}},
	"sltu": {R, []byte{, 0x3, 0x00}},
	"xor":  {R, []byte{, 0x4, 0x00}},
	"srl":  {R, []byte{, 0x5, 0x00}},
	"sra":  {R, []byte{, 0x5, 0x20}},
	"or":   {R, []byte{, 0x6, 0x00}},
	"and":  {R, []byte{, 0x7, 0x00}},
	"nop":  {R, []byte{, 0x0, 0x0}}, // Representing NOP, adjust as necessary
}
