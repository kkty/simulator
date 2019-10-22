package simulator

import "fmt"

type Label string

type ValueWithLabel struct {
	Label Label
	Value interface{}
}

type Instruction struct {
	Opcode   string
	Operands []interface{}
}

// Machine is an abstraction of physical machines.
type Machine struct {
	IntRegisters      map[string]int
	FloatRegisters    map[string]float32
	Memory            []ValueWithLabel
	ProgramCounter    int
	ConditionRegister bool
}

// NewMachine creates a Machine instance with memory of the specified size.
func NewMachine(initialMemorySize int) Machine {
	return Machine{
		IntRegisters:   make(map[string]int),
		FloatRegisters: make(map[string]float32),
		Memory:         make([]ValueWithLabel, initialMemorySize),
	}
}

// FindAddress iterates through the memory and returns the label's matching address.
func (m *Machine) FindAddress(label Label) (int, error) {
	for i, d := range m.Memory {
		if d.Label == label {
			return i, nil
		}
	}

	return 0, fmt.Errorf("%v: undefined label", label)
}
