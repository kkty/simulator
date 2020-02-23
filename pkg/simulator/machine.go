package simulator

import (
	"fmt"
	"sort"
)

type Label string

type ValueWithLabel struct {
	Label Label
	Value interface{}
}

// Instruction is an abstraction of real instructions that can be
// executed on machines. Operands can be string, int, float32 or Label
// depending on the Opcode.
type Instruction struct {
	Opcode   string
	Operands []interface{}
}

// Machine is an abstraction of physical machines.
type Machine struct {
	Registers         []uint32
	memory            []ValueWithLabel
	ProgramCounter    int32
	ConditionRegister bool
}

// NewMachine creates a Machine instance with empty registers/memory.
func NewMachine() Machine {
	return Machine{
		Registers: make([]uint32, 64),
		memory:    make([]ValueWithLabel, 100000000),
	}
}

// MemoryEntry groups a value, a label and an address.
type MemoryEntry struct {
	ValueWithLabel ValueWithLabel
	Address        int32
}

// Memory lists the machine's memory entries ordered by their addresses.
func (m *Machine) Memory() []MemoryEntry {
	memory := []MemoryEntry{}

	for address, valueWithLabel := range m.memory {
		memory = append(memory, MemoryEntry{valueWithLabel, int32(address)})
	}

	sort.Slice(memory, func(i, j int) bool { return memory[i].Address < memory[j].Address })

	return memory
}

// FindAddress iterates through the memory and returns the label's matching address.
func (m *Machine) FindAddress(label Label) (int32, error) {
	for i, d := range m.memory {
		if d.Label == label {
			return int32(i), nil
		}
	}

	return 0, fmt.Errorf("%v: undefined label", label)
}
