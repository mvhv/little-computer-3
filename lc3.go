//LC-3

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// registers
const (
	R_R0   = iota //general use (io buffer)
	R_R1          //general use
	R_R2          //general use
	R_R3          //general use
	R_R4          //general use
	R_R5          //general use
	R_R6          //general use
	R_R7          //general use (subroutine return)
	R_PC          //program counter
	R_COND        //condition register
)

// opcodes
const (
	O_BR   = iota //branch
	O_ADD         //add
	O_LD          //load
	O_ST          //store
	O_JSR         //jump register
	O_AND         //bitwise AND
	O_LDR         //load register
	O_STR         //store register
	O_RTI         //unused
	O_NOT         //bitwise NOT
	O_LDI         //store indirect
	O_JMP         //jump
	O_RES         //unused
	O_LEA         //load effective address
	O_TRAP        //execute trap
)

// flags
const (
	F_POS = 1 << iota //positive
	F_ZRO             //zero
	F_NEG             //negative
)

//traps
const (
	T_GETC  = 0x20 //get character from keyboard
	T_OUT   = 0x21 //output character
	T_PUTS  = 0x22 //output a word string
	T_IN    = 0x23 //get character from keybaord, echo to terminal
	T_PUTSP = 0x24 //output a byte string
	T_HALT  = 0x25 //halt
)

//memory-mapped registers
const (
	M_KBSR = 0xFE00 //keyboard status register
	M_KBDR = 0xFE02 //keyboard data register
	M_DSR  = 0xFE04 //display status register
	M_DDR  = 0xFE06 //display data register
	M_MCR  = 0xFFFE //machine control register
)

//vm constants
const (
	NUM_REGISTERS = 10
	MAX_UINT16    = 1<<16 - 1
)

type Machine struct {
	registers []uint16
	memory    []uint16
	ops       map[int]func(*Machine, uint16)
	traps     map[int]func(*Machine)
}

func (m *Machine) readMem(location uint16) uint16 {
	if location == M_KBSR {
		panic("Memory-mapped register not implemented: M_KBSR")
	}
	return m.memory[location]
}

func (m *Machine) writeMem(location, value uint16) {
	m.memory[location] = value
}

func (m *Machine) readReg(location uint16) uint16 {
	return m.registers[location]
}

func (m *Machine) writeReg(location, value uint16) {
	m.registers[location] = value
	m.updateFlags(value) //update condition flags after every register write
}

// func (machine *Machine) branch(instruction uint16) {
// 	destinationRegister := instruction
// }

// sign extends an integer with bitCount-bits to 16-bits
func signExtension(value, bitCount uint16) uint16 {
	if ((value >> (bitCount - 1)) & 1) != 0 {
		value |= (0xFFFF << bitCount)
	}
	return value
}

func (m *Machine) updateFlags(result uint16) {
	if result == 0 {
		m.writeReg(R_COND, F_ZRO)
	} else if result > 0 {
		m.writeReg(R_COND, F_POS)
	} else {
		m.writeReg(R_COND, F_NEG)
	}
}

//helper function to extract data from a 16-bit word
//from lsb (least significant bit), extracts n width word
func getBitRange(value, lsb, n uint16, signExtend bool) uint16 {

	temp := value

	if isBitSet(value, n) {
		temp = (value >> lsb) & ((1 << n) - 1) //THIS IS DEFINITELY BROKEN
	}

	if signExtend {
		return signExtension(temp, n)
	} else {
		return temp
	}
}

//helper function to test if bit n in value is 1
func isBitSet(value, n uint16) bool {
	return ((value >> n) & 1) == 1
}

func (m *Machine) getPC() uint16 {
	return m.registers[R_PC]
}

func (m *Machine) setPC(value uint16) {
	m.registers[R_PC] = value
}

func (m *Machine) incPC() {
	m.registers[R_PC]++
}

func (m *Machine) offsetPC(offset uint16) {
	m.registers[R_PC] += offset
}

func (m *Machine) isPosSet() bool {
	return isBitSet(m.readReg(R_COND), F_POS)
}

func (m *Machine) isNegSet() bool {
	return isBitSet(m.readReg(R_COND), F_NEG)
}

func (m *Machine) isZroSet() bool {
	return isBitSet(m.readReg(R_COND), F_ZRO)
}

func (m *Machine) dumpRegisters() {
	fmt.Printf("R_R0: %16b - 0x%x\n", m.readReg(R_R0), m.readReg(R_R0))
	fmt.Printf("R_R1: %16b - 0x%x\n", m.readReg(R_R1), m.readReg(R_R1))
	fmt.Printf("R_R2: %16b - 0x%x\n", m.readReg(R_R2), m.readReg(R_R2))
	fmt.Printf("R_R3: %16b - 0x%x\n", m.readReg(R_R3), m.readReg(R_R3))
	fmt.Printf("R_R4: %16b - 0x%x\n", m.readReg(R_R4), m.readReg(R_R4))
	fmt.Printf("R_R5: %16b - 0x%x\n", m.readReg(R_R5), m.readReg(R_R5))
	fmt.Printf("R_R6: %16b - 0x%x\n", m.readReg(R_R6), m.readReg(R_R6))
	fmt.Printf("R_R7: %16b - 0x%x\n", m.readReg(R_R7), m.readReg(R_R7))
	fmt.Printf("R_PC: %16b - 0x%x\n", m.readReg(R_PC), m.readReg(R_PC))
	fmt.Printf("R_COND: %16b - 0x%x\n", m.readReg(R_COND), m.readReg(R_COND))
}

func (m *Machine) add(instruction uint16) {
	var dest, v1, v2 uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	v1 = m.readReg(getBitRange(instruction, 6, 3, false))
	if isBitSet(instruction, 5) { //immediate mode
		v2 = getBitRange(instruction, 0, 5, true)
	} else { //register mode
		v2 = m.readReg(getBitRange(instruction, 0, 3, false))
	}

	//calculate
	res := v1 + v2
	m.writeReg(dest, res)
}

func (m *Machine) and(instruction uint16) {
	var dest, v1, v2 uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	v1 = m.readReg(getBitRange(instruction, 6, 3, false))
	if isBitSet(instruction, 5) { //immediate mode
		v2 = getBitRange(instruction, 0, 5, true)
	} else { //register mode
		v2 = m.readReg(getBitRange(instruction, 0, 3, false))
	}

	//calculate
	res := v1 & v2
	m.writeReg(dest, res)

}

func (m *Machine) branch(instruction uint16) {
	var dest, offset uint16
	var n, z, p bool

	//increment PC
	m.incPC()

	//decode
	offset = getBitRange(instruction, 0, 9, true)
	dest = m.getPC() + offset
	n = isBitSet(instruction, 11) && m.isNegSet()
	z = isBitSet(instruction, 10) && m.isZroSet()
	p = isBitSet(instruction, 9) && m.isPosSet()

	//jump if any checked condition true
	if n || z || p {
		m.setPC(dest)
	}
}

func (m *Machine) jump(instruction uint16) {
	var dest, base uint16

	//increment not used, but kept for consistency
	m.incPC()

	//decode
	base = getBitRange(instruction, 6, 3, false)
	dest = m.readReg(base)

	//jump unconditionally to address in vase register
	m.setPC(dest)
}

func (m *Machine) jumpSubroutine(instruction uint16) {
	var dest, base, offset uint16

	//increment and save PC
	m.incPC()
	m.writeMem(R_R7, m.getPC())

	//decode
	if isBitSet(instruction, 11) { //offset mode
		offset = getBitRange(instruction, 0, 11, true)
		dest = offset + m.getPC()
	} else { //base register mode
		base = getBitRange(instruction, 6, 3, false)
		dest = m.readReg(base)
	}

	//jump to subroutine
	m.setPC(dest)
}

func (m *Machine) load(instruction uint16) {
	var dest, offset, source uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	source = m.getPC() + offset

	//load data from memory
	m.writeReg(dest, m.readMem(source))
}

func (m *Machine) loadIndirect(instruction uint16) {
	var dest, offset, source uint16

	//incrment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	source = m.readMem(m.getPC() + offset)

	//load data from memory
	m.writeReg(dest, m.readMem(source))
}

func (m *Machine) loadBaseOffset(instruction uint16) {
	var dest, base, offset uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	base = getBitRange(instruction, 6, 3, false)
	offset = getBitRange(instruction, 0, 6, true)

	m.writeReg(dest, m.readMem(m.readReg(base)+offset))
}

func (m *Machine) loadEffectiveAddress(instruction uint16) {
	var dest, offset uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)

	m.writeReg(dest, m.getPC()+offset)
}

func (m *Machine) not(instruction uint16) {
	var dest, source uint16

	//increment PC
	m.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	source = getBitRange(instruction, 6, 3, false)

	m.writeReg(dest, ^m.readReg(source))
}

func (m *Machine) store(instruction uint16) {
	var source, offset, dest uint16

	//increment PC
	m.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	dest = m.getPC() + offset

	m.writeMem(dest, m.readReg(source))
}

func (m *Machine) storeIndirect(instruction uint16) {
	var source, offset, dest uint16

	//increment PC
	m.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	dest = m.readMem(m.getPC() + offset)

	m.writeMem(dest, m.readReg(source))
}

func (m *Machine) storeBaseOffset(instruction uint16) {
	var source, base, offset, dest uint16

	//increment PC
	m.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	base = getBitRange(instruction, 6, 3, false)
	offset = getBitRange(instruction, 0, 6, true)
	dest = m.readReg(base) + offset

	m.writeMem(dest, m.readReg(source))
}

func (m *Machine) trapPutString() {
	for i := m.readReg(R_R0); m.readMem(i) != 0x00000; i++ {
		fmt.Print(rune(m.readMem(i)))
	}
	fmt.Print()
}

func (m *Machine) trapGetChar() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	m.writeReg(R_R0, uint16(input[0]))
}

func (m *Machine) trapInput() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	m.writeReg((R_R0), uint16(input[0]))
}

func (m *Machine) trapOutput() {
	fmt.Print(rune(m.readReg(R_R0)))
}

func (m *Machine) trapPutStringP() {
	for i := m.readReg(R_R0); m.readMem(i) != 0x00000; i++ {
		fmt.Print(rune(m.readMem(i)))
	}
}

func (m *Machine) trapHalt() {
	panic("Execution Halted")
}

func (m *Machine) trap(instruction uint16) {
	var trapVector uint16

	//increment PC and store in R7
	m.incPC()
	m.writeReg(R_R7, m.getPC())

	//decode
	trapVector = getBitRange(instruction, 0, 8, false)

	//traps implemented in go for simplicity
	switch trapVector {
	case T_GETC:
		trapGetChar()
	case T_OUT:
		trapOutput()
	case T_PUTS:
		trapPutString()
	case T_IN:
		trapInput()
	case T_PUTSP:
		trapPutStringP()
	case T_HALT:
		trapHalt()
	}

	//set PC to trap location for system call
	//m.setPC(trapVector)
}

func NewMachine() Machine {
	m := Machine{
		registers: make([]uint16, 10),
		memory:    make([]uint16, MAX_UINT16),
		ops: map[int]func(*Machine, uint16){
			O_BR:   (*Machine).branch,
			O_ADD:  (*Machine).add,
			O_LD:   (*Machine).load,
			O_ST:   (*Machine).store,                                    //store
			O_JSR:  (*Machine).jumpSubroutine,                           //jump to subroutine
			O_AND:  (*Machine).and,                                      //bitwise AND
			O_LDR:  (*Machine).loadBaseOffset,                           //load base and offset
			O_STR:  (*Machine).storeBaseOffset,                          //store base and offset
			O_RTI:  func() { panic("Bad opcode: RTI not implemented") }, //unused
			O_NOT:  (*Machine).not,                                      //bitwise NOT
			O_LDI:  (*Machine).loadIndirect,                             //load indirect
			O_JMP:  (*Machine).jump,                                     //jump
			O_RES:  func() { panic("Bad opcode: RES not implemeted") },  //reserved
			O_LEA:  (*Machine).loadEffectiveAddress,                     //load effective address
			O_TRAP: (*Machine).trap,                                     //execute trap
		},
		traps: map[int]func(*Machine){
			T_GETC: (*Machine).trapGetChar,
		},
	}

	// default PC start location is 0x3000
	m.setPC(0x3000)
	return m
}

func (m *Machine) LoadImage(filePath string) {
	var origin uint16

	image, _ := os.Open(filePath)
	binary.Read(image, binary.LittleEndian, &origin)
	binary.Read(image, binary.LittleEndian, m.memory[origin:])
}

func (m *Machine) Run() {
	for {
		// fetch
		instruction := m.readMem(m.getPC())
		opcode := instruction >> 12 //right shift 12 to extract highest 4 bits
		switch opcode {
		case O_BR: //branch
			m.branch(instruction)
		case O_ADD: //add
			m.add(instruction)
		case O_LD: //load
			m.load(instruction)
		case O_ST: //store
			m.store(instruction)
		case O_JSR: //jump to subroutine
			m.jumpSubroutine(instruction)
		case O_AND: //bitwise AND
			m.and(instruction)
		case O_LDR: //load base and offset
			m.loadBaseOffset(instruction)
		case O_STR: //store base and offset
			m.storeBaseOffset(instruction)
		case O_RTI: //unused
			panic("Bad opcode: RTI not implemented")
		case NOT: //bitwise NOT
			m.not(instruction)
		case O_LDI: //store indirect
			m.loadIndirect(instruction)
		case O_JMP: //jump
			m.jump(instruction)
		case O_RES: //unused
			panic("Bad opcode: RES not implemeted")
		case LEA: //load effective address
			m.loadEffectiveAddress(instruction)
		case O_TRAP: //execute trap
			m.trap(instruction)
		default:
			panic("Bad opcode: Other")
		}
	}
}
