package machine

import (
	"fmt"
	"bufio"
	"os"
)

func (vm *LC3VM) trapPutString() {
	for i := vm.readReg(R_R0); vm.readMem(i) != 0x00000; i++ {
		fmt.Print(rune(vm.readMem(i)))
	}
	fmt.Print()
}

func (vm *LC3VM) trapGetChar() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	vm.writeReg(R_R0, uint16(input[0]))
}

func (vm *LC3VM) trapInput() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	vm.writeReg((R_R0), uint16(input[0]))
}

func (vm *LC3VM) trapOutput() {
	fmt.Print(rune(vm.readReg(R_R0)))
}

func (vm *LC3VM) trapPutStringP() {
	for i := vm.readReg(R_R0); vm.readMem(i) != 0x00000; i++ {
		fmt.Print(rune(vm.readMem(i)))
	}
}

func (vm *LC3VM) trapHalt() {
	panic("Execution Halted")
}

// Traps map for trap vector -> function lookup
var traps = map[int]func(*LC3VM){
	T_GETC: (*LC3VM).trapGetChar,
	T_OUT: (*LC3VM).trapOutput,
	T_PUTS: (*LC3VM).trapPutString,
	T_IN: (*LC3VM).trapInput,
	T_PUTSP: (*LC3VM).trapPutStringP,
	T_HALT: (*LC3VM).trapHalt,
}