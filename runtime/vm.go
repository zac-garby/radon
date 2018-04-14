package runtime

import (
	"io"
	"os"

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
