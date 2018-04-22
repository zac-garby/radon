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
	halted    bool

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
		halted:     false,
		Interrupts: make(chan Interrupt, InterruptQueueSize),
		Out:        os.Stdout,
	}
}

// makeFrame makes a Frame instance with the given parameters. offset is set to 0 and a new
// data stack is created. The variables in args, if it's non-nil, are declared in the Frame's
// store (note: declared, not assigned).
func (v *VM) makeFrame(code bytecode.Code, args, store *Store, constants []object.Object, names []string) *Frame {
	frame := &Frame{
		code:      code,
		store:     store,
		offset:    0,
		stack:     NewStack(),
		constants: constants,
		names:     names,
		vm:        v,
	}

	if args != nil {
		for k, v := range args.Data {
			frame.store.Set(k, v.Value, true)
		}
	}

	return frame
}

func (v *VM) pushFrame(frame *Frame) {
	v.frames = append(v.frames, frame)
}

func (v *VM) popFrame() *Frame {
	f := v.frames[len(v.frames)-1]
	v.frames = v.frames[:len(v.frames)-1]
	return f
}

// ExtractValue returns the top value from the top frame, if both those things exist.
func (v *VM) ExtractValue() (object.Object, error) {
	if len(v.frames) < 1 || v.frames[0].stack.Len() < 1 {
		return nil, nil
	}
	return v.frames[0].stack.Top()
}

// Error returns the VM's error, if one exists. nil if there is no error.
func (v *VM) Error() error {
	return v.err
}
