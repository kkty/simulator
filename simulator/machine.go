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

type Instruction struct {
	Opcode   string
	Operands []interface{}
}

// Machine is an abstraction of physical machines.
type Machine struct {
	IntRegisters      map[string]int
	FloatRegisters    map[string]float32
	Memory            map[int]ValueWithLabel
	ProgramCounter    int
	ConditionRegister bool
	findAddressCache  map[Label]int
}

// NewMachine creates a Machine instance with memory of the specified size.
func NewMachine(entrypoint string) Machine {
	return Machine{
		IntRegisters:     make(map[string]int),
		FloatRegisters:   make(map[string]float32),
		Memory:           make(map[int]ValueWithLabel),
		findAddressCache: make(map[Label]int),
	}
}

func (m *Machine) setValueToMemory(address int, value interface{}) {
	m.Memory[address] = ValueWithLabel{m.Memory[address].Label, value}
}

// FindAddress iterates through the memory and returns the label's matching address.
func (m *Machine) FindAddress(label Label) (int, error) {
	if address, ok := m.findAddressCache[label]; ok {
		return address, nil
	}

	for i, d := range m.Memory {
		if d.Label == label {
			m.findAddressCache[label] = i
			return i, nil
		}
	}

	return 0, fmt.Errorf("%v: undefined label", label)
}
