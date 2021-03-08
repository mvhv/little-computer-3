package machine

func (vm *LC3VM) opAdd(instruction uint16) {
	var dest, v1, v2 uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	v1 = vm.readReg(getBitRange(instruction, 6, 3, false))
	if isBitSet(instruction, 5) { //immediate mode
		v2 = getBitRange(instruction, 0, 5, true)
	} else { //register mode
		v2 = vm.readReg(getBitRange(instruction, 0, 3, false))
	}

	//calculate
	res := v1 + v2
	vm.writeReg(dest, res)
}

func (vm *LC3VM) opAnd(instruction uint16) {
	var dest, v1, v2 uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	v1 = vm.readReg(getBitRange(instruction, 6, 3, false))
	if isBitSet(instruction, 5) { //immediate mode
		v2 = getBitRange(instruction, 0, 5, true)
	} else { //register mode
		v2 = vm.readReg(getBitRange(instruction, 0, 3, false))
	}

	//calculate
	res := v1 & v2
	vm.writeReg(dest, res)

}

func (vm *LC3VM) opBranch(instruction uint16) {
	var dest, offset uint16
	var n, z, p bool

	//increment PC
	vm.incPC()

	//decode
	offset = getBitRange(instruction, 0, 9, true)
	dest = vm.getPC() + offset
	n = isBitSet(instruction, 11) && vm.isNegSet()
	z = isBitSet(instruction, 10) && vm.isZroSet()
	p = isBitSet(instruction, 9) && vm.isPosSet()

	//jump if any checked condition true
	if n || z || p {
		vm.setPC(dest)
	}
}

func (vm *LC3VM) opJump(instruction uint16) {
	var dest, base uint16

	//increment not used, but kept for consistency
	vm.incPC()

	//decode
	base = getBitRange(instruction, 6, 3, false)
	dest = vm.readReg(base)

	//jump unconditionally to address in base register
	vm.setPC(dest)
}

func (vm *LC3VM) opJumpSubroutine(instruction uint16) {
	var dest, base, offset uint16

	//increment and save PC
	vm.incPC()
	vm.writeMem(R_R7, vm.getPC())

	//decode
	if isBitSet(instruction, 11) { //offset mode
		offset = getBitRange(instruction, 0, 11, true)
		dest = offset + vm.getPC()
	} else { //base register mode
		base = getBitRange(instruction, 6, 3, false)
		dest = vm.readReg(base)
	}

	//jump to subroutine
	vm.setPC(dest)
}

func (vm *LC3VM) opLoad(instruction uint16) {
	var dest, offset, source uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	source = vm.getPC() + offset

	//load data from memory
	vm.writeReg(dest, vm.readMem(source))
}

func (vm *LC3VM) opLoadIndirect(instruction uint16) {
	var dest, offset, source uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	source = vm.readMem(vm.getPC() + offset)

	//load data from memory
	vm.writeReg(dest, vm.readMem(source))
}

func (vm *LC3VM) opLoadBaseOffset(instruction uint16) {
	var dest, base, offset uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	base = getBitRange(instruction, 6, 3, false)
	offset = getBitRange(instruction, 0, 6, true)

	vm.writeReg(dest, vm.readMem(vm.readReg(base)+offset))
}

func (vm *LC3VM) opLoadEffectiveAddress(instruction uint16) {
	var dest, offset uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)

	vm.writeReg(dest, vm.getPC()+offset)
}

func (vm *LC3VM) opNot(instruction uint16) {
	var dest, source uint16

	//increment PC
	vm.incPC()

	//decode
	dest = getBitRange(instruction, 9, 3, false)
	source = getBitRange(instruction, 6, 3, false)

	vm.writeReg(dest, ^vm.readReg(source))
}

func (vm *LC3VM) opStore(instruction uint16) {
	var source, offset, dest uint16

	//increment PC
	vm.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	dest = vm.getPC() + offset

	vm.writeMem(dest, vm.readReg(source))
}

func (vm *LC3VM) opStoreIndirect(instruction uint16) {
	var source, offset, dest uint16

	//increment PC
	vm.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	offset = getBitRange(instruction, 0, 9, true)
	dest = vm.readMem(vm.getPC() + offset)

	vm.writeMem(dest, vm.readReg(source))
}

func (vm *LC3VM) opStoreBaseOffset(instruction uint16) {
	var source, base, offset, dest uint16

	//increment PC
	vm.incPC()

	//decode
	source = getBitRange(instruction, 9, 3, false)
	base = getBitRange(instruction, 6, 3, false)
	offset = getBitRange(instruction, 0, 6, true)
	dest = vm.readReg(base) + offset

	vm.writeMem(dest, vm.readReg(source))
}

func (vm *LC3VM) opSystemCall(instruction uint16) {
	var trapVector uint16

	//increment PC and store in R7
	vm.incPC()
	vm.writeReg(R_R7, vm.getPC())

	//decode
	trapVector = getBitRange(instruction, 0, 8, false)

	//traps implemented in Go for simplicity
	traps[int(trapVector)](vm)

	//set PC to trap location for system call
	//vm.setPC(trapVector)
}

func (vm *LC3VM) opNotImplemented(instruction uint16) {
	panic("Bad opcode")
}

// ops map for opcode -> function lookup
var ops = map[int]func(*LC3VM, uint16){
	O_BR:   (*LC3VM).opBranch, //branch
	O_ADD:  (*LC3VM).opAdd, //add
	O_LD:   (*LC3VM).opLoad,
	O_ST:   (*LC3VM).opStore,                                    //store
	O_JSR:  (*LC3VM).opJumpSubroutine,                           //jump to subroutine
	O_AND:  (*LC3VM).opAnd,                                      //bitwise AND
	O_LDR:  (*LC3VM).opLoadBaseOffset,                           //load base and offset
	O_STR:  (*LC3VM).opStoreBaseOffset,                          //store base and offset
	O_RTI:  (*LC3VM).opNotImplemented, //unused
	O_NOT:  (*LC3VM).opNot,                                      //bitwise NOT
	O_LDI:  (*LC3VM).opLoadIndirect,                             //load indirect
	O_JMP:  (*LC3VM).opJump,                                     //unconditional jump
	O_RES:  (*LC3VM).opNotImplemented,  //reserved
	O_LEA:  (*LC3VM).opLoadEffectiveAddress,                     //load effective address
	O_TRAP: (*LC3VM).opSystemCall,                               //execute system call
}