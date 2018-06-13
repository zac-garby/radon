package runtime

// Interrupt defines the type of an interrupt which, when sent to the virtual machine through
// the Interrupts chan, stops the fetch-decode-execute cycle and does something.
type Interrupt int

const (
	_ Interrupt = iota

	// Stop stops the execution
	Stop

	// Pause pauses the execution until Resume is sent. All other interrupts are ignored until
	// Resume is sent.
	Pause

	// Resume resumes a paused virtual machine.
	Resume
)
