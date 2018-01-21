package vm

// An Interrupt, when sent to the virtual machine, interrupts
// the fetch-decode-execute cycle and does something, such as
// halting the execution.
type Interrupt int

const (
	_ Interrupt = iota

	// Halts execution
	InterruptHalt

	// Pauses execution, until Resume is sent
	InterruptPause

	// Resumes paused execution
	InterruptResume
)
