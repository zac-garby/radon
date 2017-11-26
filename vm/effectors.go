package vm

import (
	"github.com/Zac-Garby/lang/bytecode"
)

type effector func(f *Frame, i bytecode.Instruction)

var effectors map[byte]effector

func init() {
	effectors = map[byte]effector{}
}
