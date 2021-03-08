package machine

// sign extends an integer with n-bits to 16-bits
func signExtension(value, n uint16) uint16 {
	temp := value
	if isBitSet(value, n) {
		temp |= (MAX_UINT16 << n)
	}
	return temp
}

// extracts and returns n-bit word starting from lsb (least significant bit)
func getBitRange(value, lsb, n uint16, signExtend bool) uint16 {
	temp := (value >> lsb) & ((1 << n) - 1)
	if signExtend {
		return signExtension(temp, n)
	}

	return temp
}

// tests if the bit at poistion n is set, and returns true if so
func isBitSet(value, n uint16) bool {
	return ((value >> n) & 1) == 1
}