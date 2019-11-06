package simulator

import "fmt"

const zeroRegisterName = "$zero"
const raRegisterName = "$ra"
const spRegisterName = "$sp"
const hpRegisterName = "$hp"

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
	IntRegisters      map[string]int32
	FloatRegisters    map[string]float32
	memory            map[int32]ValueWithLabel
	ProgramCounter    int32
	ConditionRegister bool
	findAddressCache  map[Label]int32
}

// NewMachine creates a Machine instance with empty registers/memory.
func NewMachine() Machine {
	return Machine{
		IntRegisters:     make(map[string]int32),
		FloatRegisters:   make(map[string]float32),
		memory:           make(map[int32]ValueWithLabel),
		findAddressCache: make(map[Label]int32),
	}
}

// setValueToMemory sets a value to the memory without modifying label values.
func (m *Machine) setValueToMemory(address int32, value interface{}) {
	m.memory[address] = ValueWithLabel{m.memory[address].Label, value}
}

// FindAddress iterates through the memory and returns the label's matching address.
func (m *Machine) FindAddress(label Label) (int32, error) {
	if address, ok := m.findAddressCache[label]; ok {
		return address, nil
	}

	for i, d := range m.memory {
		if d.Label == label {
			m.findAddressCache[label] = i
			return i, nil
		}
	}

	return 0, fmt.Errorf("%v: undefined label", label)
}
