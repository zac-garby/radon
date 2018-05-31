package runtime

import "github.com/Zac-Garby/radon/bytecode"

// An Effector is the function which is called for a particular instruction.
type Effector func(v *VM, f *Frame, arg rune) error

// Effectors stores the effector function for each instruction. The effector
// for an instruction n can be accessed via Effectors[n].
var Effectors [256]Effector

func init() {
	Effectors[bytecode.Nop] = func(v *VM, f *Frame, arg rune) error { return nil }
	Effectors[bytecode.NopArg] = func(v *VM, f *Frame, arg rune) error { return nil }
}
