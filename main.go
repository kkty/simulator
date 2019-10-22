package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

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

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i := m.Memory[m.ProgramCounter].Value.(Instruction)

	switch opcode := i.Opcode; opcode {
	case "add":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] + m.IntRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "sub":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] - m.IntRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "addi":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] + i.Operands[2].(int)
		m.ProgramCounter++
	case "subi":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] - i.Operands[2].(int)
		m.ProgramCounter++
	case "lui":
		m.IntRegisters[i.Operands[0].(string)] = i.Operands[1].(int) << 16
		m.ProgramCounter++
	case "ori":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] | i.Operands[2].(int)
		m.ProgramCounter++
	case "slt":
		if m.IntRegisters[i.Operands[1].(string)] < m.IntRegisters[i.Operands[2].(string)] {
			m.IntRegisters[i.Operands[0].(string)] = 1
		} else {
			m.IntRegisters[i.Operands[0].(string)] = 0
		}
		m.ProgramCounter++
	case "sll":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] << i.Operands[2].(int)
		m.ProgramCounter++
	case "sllv":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] << m.IntRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "j":
		m.ProgramCounter = i.Operands[0].(int)
	case "jal":
		m.IntRegisters["$ra"] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int)
	case "jr":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(string)]
	case "beq":
		if m.IntRegisters[i.Operands[0].(string)] == m.IntRegisters[i.Operands[1].(string)] {
			m.ProgramCounter += i.Operands[2].(int) + 1
		} else {
			m.ProgramCounter++
		}
	case "c.eq.s":
		if m.FloatRegisters[i.Operands[0].(string)] == m.FloatRegisters[i.Operands[1].(string)] {
			m.ConditionRegister = true
		} else {
			m.ConditionRegister = false
		}

		m.ProgramCounter++
	case "c.le.s":
		if m.FloatRegisters[i.Operands[0].(string)] <= m.FloatRegisters[i.Operands[1].(string)] {
			m.ConditionRegister = true
		} else {
			m.ConditionRegister = false
		}

		m.ProgramCounter++
	case "bc1t":
		if m.ConditionRegister {
			m.ProgramCounter += 1 + i.Operands[1].(int)
		} else {
			m.ProgramCounter++
		}
	case "lw":
		m.IntRegisters[i.Operands[0].(string)] = m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value.(int)
		m.ProgramCounter++
	case "lwc1":
		m.FloatRegisters[i.Operands[0].(string)] = m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value.(float32)
		m.FloatRegisters[i.Operands[0].(string)] = m.FloatRegisters[i.Operands[1].(string)] * m.FloatRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "sw":
		m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value = m.IntRegisters[i.Operands[0].(string)]
		m.ProgramCounter++
	case "swc1":
		m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value = m.FloatRegisters[i.Operands[0].(string)]
		m.ProgramCounter++
	case "add.s":
		m.FloatRegisters[i.Operands[0].(string)] = m.FloatRegisters[i.Operands[1].(string)] + m.FloatRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "sub.s":
		m.FloatRegisters[i.Operands[0].(string)] = m.FloatRegisters[i.Operands[1].(string)] - m.FloatRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "mul.s":
		m.FloatRegisters[i.Operands[0].(string)] = m.FloatRegisters[i.Operands[1].(string)] * m.FloatRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "div.s":
		m.FloatRegisters[i.Operands[0].(string)] = m.FloatRegisters[i.Operands[1].(string)] / m.FloatRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "out":
		fmt.Println(m.IntRegisters[i.Operands[0].(string)])
		m.ProgramCounter++
	case "exit":
		return true, nil
	default:
		return false, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return false, nil
}

func parseInstruction(fields []string) (Instruction, error) {
	immediateOrLabel := func(s string) interface{} {
		i, err := strconv.Atoi(s)

		if err != nil {
			return Label(s)
		}

		return i
	}

	opcode := fields[0]
	var operands []interface{}

	switch opcode {
	case "add":
		fallthrough
	case "sub":
		fallthrough
	case "slt":
		fallthrough
	case "sllv":
		fallthrough
	case "add.s":
		fallthrough
	case "sub.s":
		fallthrough
	case "mul.s":
		fallthrough
	case "div.s":
		operands = append(operands, fields[1], fields[2], fields[3])
	case "sll":
		fallthrough
	case "addi":
		fallthrough
	case "lui":
		fallthrough
	case "ori":
		fallthrough
	case "beq":
		operands = append(operands, fields[1], fields[2], immediateOrLabel(fields[3]))
	case "lw":
		fallthrough
	case "lwc1":
		fallthrough
	case "sw":
		fallthrough
	case "swc1":
		operands = append(operands, fields[1], immediateOrLabel(fields[2]), fields[3])
	case "bc1t":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "c.eq.s":
		fallthrough
	case "c.le.s":
		operands = append(operands, fields[1], fields[2])
	case "j":
		fallthrough
	case "jal":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "jr":
		fallthrough
	case "out":
		operands = append(operands, fields[1])
	default:
		return Instruction{}, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return Instruction{
		Opcode:   opcode,
		Operands: operands,
	}, nil
}

func parseData(fields []string) (ValueWithLabel, error) {
	if len(fields) != 3 {
		return ValueWithLabel{}, errors.New("invalid syntax")
	}

	switch t := fields[1]; t {
	case ".float":
		f, err := strconv.ParseFloat(fields[2], 32)

		if err != nil {
			return ValueWithLabel{}, err
		}

		return ValueWithLabel{
			Label(strings.TrimSuffix(fields[0], ":")),
			f,
		}, nil
	default:
		return ValueWithLabel{}, errors.New("invalid data type")
	}
}

// Load loads a program onto the memory.
func (m *Machine) Load(program string) error {
	// The default section is "text".
	section := "text"

	var nextLabel Label

	for _, line := range strings.Split(program, "\n") {
		// Replaces "(", ")" or "," with whitespaces.
		for _, c := range []string{"(", ")", ","} {
			line = strings.ReplaceAll(line, c, " ")
		}

		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		wrapError := func(err error) error {
			return fmt.Errorf("parse %v: %w", fields, err)
		}

		// Changes section.
		if strings.HasPrefix(fields[0], ".") {
			section = strings.TrimPrefix(fields[0], ".")
			continue
		}

		switch section {
		case "data":
			valueWithLabel, err := parseData(fields)

			if err != nil {
				return wrapError(err)
			}

			m.Memory = append(m.Memory, valueWithLabel)
		case "text":
			if strings.HasSuffix(fields[0], ":") {
				nextLabel = Label(strings.TrimSuffix(fields[0], ":"))
				continue
			}

			instruction, err := parseInstruction(fields)

			if err != nil {
				return wrapError(err)
			}

			m.Memory = append(m.Memory, ValueWithLabel{nextLabel, instruction})

			nextLabel = ""
		default:
			return wrapError(fmt.Errorf("%v: invalid section", section))
		}
	}

	// Iterates thorough the memory and replaces labels with address values.
	for i := 0; i < len(m.Memory); i++ {
		if instruction, ok := m.Memory[i].Value.(Instruction); ok {
			for j := 0; j < len(instruction.Operands); j++ {
				if label, ok := instruction.Operands[j].(Label); ok {
					var err error
					instruction.Operands[j], err = m.FindAddress(label)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func main() {
	b, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	m := NewMachine(1000)

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress("min_caml_start")

	if err != nil {
		log.Fatal(err)
	}
}
