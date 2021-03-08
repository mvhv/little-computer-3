package machine

import (
	"encoding/binary"
	"fmt"
	"os"
)

// LC3VM is a virtual machine implementation of the LC3 architecture
type LC3VM struct {
	registers []uint16
	memory    []uint16
	ops       map[int]func(*LC3VM, uint16)
	traps     map[int]func(*LC3VM)
}

func (vm *LC3VM) readMem(location uint16) uint16 {
	if location == M_KBSR {
		panic("Memory-mapped register not implemented: M_KBSR")
	}
	return vm.memory[location]
}

func (vm *LC3VM) writeMem(location, value uint16) {
	vm.memory[location] = value
}

func (vm *LC3VM) readReg(location uint16) uint16 {
	return vm.registers[location]
}

func (vm *LC3VM) writeReg(location, value uint16) {
	vm.registers[location] = value
	vm.updateFlags(value) //update condition flags after every register write
}

func (vm *LC3VM) updateFlags(result uint16) {
	if result == 0 {
		vm.writeReg(R_COND, F_ZRO)
	} else if result > 0 {
		vm.writeReg(R_COND, F_POS)
	} else {
		vm.writeReg(R_COND, F_NEG)
	}
}

func (vm *LC3VM) getPC() uint16 {
	return vm.registers[R_PC]
}

func (vm *LC3VM) setPC(value uint16) {
	vm.registers[R_PC] = value
}

func (vm *LC3VM) incPC() {
	vm.registers[R_PC]++
}

func (vm *LC3VM) offsetPC(offset uint16) {
	vm.registers[R_PC] += offset
}

func (vm *LC3VM) isPosSet() bool {
	return isBitSet(vm.readReg(R_COND), F_POS)
}

func (vm *LC3VM) isNegSet() bool {
	return isBitSet(vm.readReg(R_COND), F_NEG)
}

func (vm *LC3VM) isZroSet() bool {
	return isBitSet(vm.readReg(R_COND), F_ZRO)
}

func (vm *LC3VM) dumpRegisters() {
	fmt.Printf("R_R0: %16b - 0x%x\n", vm.readReg(R_R0), vm.readReg(R_R0))
	fmt.Printf("R_R1: %16b - 0x%x\n", vm.readReg(R_R1), vm.readReg(R_R1))
	fmt.Printf("R_R2: %16b - 0x%x\n", vm.readReg(R_R2), vm.readReg(R_R2))
	fmt.Printf("R_R3: %16b - 0x%x\n", vm.readReg(R_R3), vm.readReg(R_R3))
	fmt.Printf("R_R4: %16b - 0x%x\n", vm.readReg(R_R4), vm.readReg(R_R4))
	fmt.Printf("R_R5: %16b - 0x%x\n", vm.readReg(R_R5), vm.readReg(R_R5))
	fmt.Printf("R_R6: %16b - 0x%x\n", vm.readReg(R_R6), vm.readReg(R_R6))
	fmt.Printf("R_R7: %16b - 0x%x\n", vm.readReg(R_R7), vm.readReg(R_R7))
	fmt.Printf("R_PC: %16b - 0x%x\n", vm.readReg(R_PC), vm.readReg(R_PC))
	fmt.Printf("R_COND: %16b - 0x%x\n", vm.readReg(R_COND), vm.readReg(R_COND))
}

func NewLC3VM() *LC3VM {
	vm := LC3VM{
		registers: make([]uint16, 10),
		memory:    make([]uint16, MAX_UINT16),
	}

	// default PC start location is 0x3000
	vm.setPC(0x3000)
	return &vm
}

func (vm *LC3VM) LoadImage(filePath string) {
	var origin uint16

	image, _ := os.Open(filePath)
	binary.Read(image, binary.LittleEndian, &origin)
	binary.Read(image, binary.LittleEndian, vm.memory[origin:])
}

func (vm *LC3VM) Run() {
	for {
		// fetch
		instruction := vm.readMem(vm.getPC())
		opcode := instruction >> 12 //right shift 12 to extract highest 4 bits
		ops[int(opcode)](vm, instruction)
	}
}
