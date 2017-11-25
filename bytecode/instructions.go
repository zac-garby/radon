package bytecode

// 0-9: stack operations
const (
	// Pop pops the stack
	Pop byte = iota

	// Dup duplicates the top item, so [x, y, z] -> [x, y, z, z]
	Dup

	// Rot rotates the top two items, so [x, y, z] -> [x, z, y]
	Rot
)

// 10-19: load/store
const (
	// LoadConst loads a constant by index
	LoadConst byte = iota + 10

	// LoadName loads a name by name index
	LoadName

	// StoreName stores the top item
	StoreName

	// LoadField pops two items, essentially does second[top]
	LoadField

	// StoreField pops three items, essentially does second[top] = third
	StoreField
)

// 20-39: operators
const (
	// Unary operators pop one item and do something with it
	UnaryInvert byte = iota + 20
	UnaryNegate

	// Binary operators pop two items and do something with them
	BinaryAdd byte = iota + 25
	BinarySubtract
	BinaryMultiply
	BinaryDivide
	BinaryExponent
	BinaryFloorDiv
	BinaryMod
	BinaryOr
	BinaryAnd
	BinaryBitOr
	BinaryBitAnd
	BinaryEquals
	BinaryNotEqual
	BinaryLessThan
	BinaryMoreThan
	BinaryLessEq
	BinaryMoreEq
)

// 50-59: using functions/blocks
const (
	// CallFn calls the function at the top of the stack,
	// popping arguments as necessary
	CallFn byte = iota + 50

	// Return skips to the end of the context
	Return
)

// 60-89: builtin functions
const (
	// Print prints the item at the top of the stack
	Print byte = iota + 60

	// Println prints the item at the top of the stack,
	// with a trailing new line
	Println

	// Length pushes the length of the collection at
	// the top of the stack
	Length
)

// 90-99: control flow
const (
	// Jump unconditionally jumps to the given offset
	Jump byte = iota + 90

	// JumpIfTrue jumps to the given offset if the top item is truthy
	JumpIfTrue

	// JumpIfFalse jumps to the given offset if the top item is falsey
	JumpIfFalse

	// Break jumps to the LoopEnd instruction of the innermost loop
	Break

	// Next jumps to the LoopStart instruction of the innermost loop
	Next

	// LoopStart pushes the start and end positions for the loop
	LoopStart

	// LoopEnd pops the start and end positions
	LoopEnd
)

// 100-109: data constructors
const (
	// MakeList makes an array object from the top n items
	MakeList byte = iota + 100

	// MakeTuple makes a tuple from the top n items
	MakeTuple

	// MakeMap makes a map from the top n * 2 items.
	// The top n*2 items should be in key, val, ..., key, val order
	MakeMap
)
