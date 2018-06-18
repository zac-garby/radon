package runtime

import (
	"io"
	"os"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
)

var (
	// InterruptQueueSize specifies the buffer size of the interrupt queue.
	InterruptQueueSize = 64
)

// A VM (for Virtual Machine), interprets bytecode. It's a stack machine, so operates by
// pushing to and popping from data stacks, instead of using registers like a computer.
type VM struct {
	frames    []*Frame
	frame     *Frame
	returnVal object.Object
	err       error
	storePool *StorePool

	// Interrupts is a queue (in the form of a channel) which, when an interrupt is added,
	// will stop the virtual machine and do something based on the nature of the interrupt.
	Interrupts chan Interrupt

	// Out is the io.Writer to which the virtual machine outputs to. This includes functions
	// like print, but also errors and various messages.
	Out io.Writer
}

// New creates a new virtual machine.
func New() *VM {
	return &VM{
		frames:     make([]*Frame, 0),
		frame:      nil,
		returnVal:  nil,
		err:        nil,
		storePool:  NewStorePool(),
		Interrupts: make(chan Interrupt, InterruptQueueSize),
		Out:        os.Stdout,
	}
}

// MakeFrame makes a Frame instance with the given parameters. offset is set to 0 and a new
// data stack is created. The variables in args, if it's non-nil, are declared in the Frame's
// store (note: declared, not assigned).
func (v *VM) MakeFrame(code bytecode.Code, args, store *Store, constants []object.Object, names []string, jumps []int) *Frame {
	frame := &Frame{
		code:        code,
		stores:      []*Store{store},
		offset:      0,
		stack:       NewStack(),
		constants:   constants,
		names:       names,
		jumps:       jumps,
		vm:          v,
		matchInputs: make([]object.Object, 0),
	}

	if args != nil {
		for k, v := range args.Data {
			frame.store().Set(k, v.Value, true)
		}
	}

	return frame
}

// PushFrame pushes a frame to the top of the call stack.
func (v *VM) PushFrame(frame *Frame) {
	v.frames = append(v.frames, frame)
}

// PopFrame pops the frame from the top of the call stack.
func (v *VM) PopFrame() *Frame {
	f := v.frames[len(v.frames)-1]
	v.frames = v.frames[:len(v.frames)-1]
	return f
}

// ExtractValue returns the top value from the top frame, if both those things exist.
func (v *VM) ExtractValue() object.Object {
	if len(v.frames) < 1 || v.frames[0].stack.Len() < 1 {
		return nil
	}

	top, err := v.frames[0].stack.Top()
	if err != nil {
		return nil
	}

	return top
}

// Error returns the VM's error, if one exists. nil if there is no error.
func (v *VM) Error() error {
	return v.err
}

// Run executes a virtual machine, starting from the most recently pushed frame. If, after
// execution, any values are left in the top frame, the top one will be returned. It will
// also return, if any, a runtime error.
func (v *VM) Run() (object.Object, error) {
main:
	for {
		// Handle interrupts first
		for len(v.Interrupts) > 0 {
			interrupt := <-v.Interrupts

			switch interrupt {
			case Stop:
				break main

			case Pause:
				for {
					if <-v.Interrupts == Resume {
						break
					}
				}
			}
		}

		if len(v.frames) == 0 {
			break
		}

		// This may be able to be optimized slightly by putting it outside the loop
		top := v.frames[len(v.frames)-1]

		if top.offset >= len(top.code) {
			if len(v.frames) == 1 {
				break
			} else {
				v.PopFrame()

				if top.stack.Len() > 0 {
					next := v.frames[len(v.frames)-1]

					ret, err := top.stack.Pop()
					if err != nil {
						v.err = err
						break
					}

					if err := next.stack.Push(ret); err != nil {
						v.err = err
						break
					}
				}
			}

			continue
		}

		// Fetch
		instr := top.code[top.offset]
		top.offset++

		// Decode
		eff := Effectors[instr.Code]
		if eff == nil {
			v.err = makeError(InternalError, "instruction %s not yet implemented", instr.Name)
			break
		}

		// Execute :)
		if err := eff(v, top, instr.Arg); err != nil {
			v.err = err
			break
		}
	}

	return v.ExtractValue(), v.err
}
