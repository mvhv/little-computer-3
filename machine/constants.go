package machine

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