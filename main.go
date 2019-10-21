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

const ra = "$ra"
const sp = "$sp"
const hp = "$hp"

type Label string

type Machine struct {
	IntRegisters   map[string]int
	FloatRegisters map[string]float32
	Memory         []struct {
		Label Label
		Value interface{}
	}
	ProgramCounter    int
	ConditionRegister bool
}

func NewMachine() Machine {
	return Machine{
		IntRegisters:   make(map[string]int),
		FloatRegisters: make(map[string]float32),
	}
}

var ErrUndefinedLabel = errors.New("undefined label")
var ErrUndefinedOpcode = errors.New("undefined opcode")

func (m *Machine) FindAddress(label Label) (int, error) {
	for i, d := range m.Memory {
		if d.Label == label {
			return i, nil
		}
	}

	return 0, fmt.Errorf("%v: %w", label, ErrUndefinedLabel)
}

func (m *Machine) Step() error {
	i := m.Memory[m.ProgramCounter].Value.(Instruction)

	switch i.Opcode {
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
		m.IntRegisters[ra] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int)
	case "jr":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(string)]
	case "beq":
		if m.IntRegisters[i.Operands[0].(string)] == m.IntRegisters[i.Operands[1].(string)] {
			m.ProgramCounter += i.Operands[2].(int) + 1
		} else {
			m.ProgramCounter++
		}
	case "ceqs":
		if m.FloatRegisters[i.Operands[0].(string)] == m.FloatRegisters[i.Operands[1].(string)] {
			m.ConditionRegister = true
		} else {
			m.ConditionRegister = false
		}

		m.ProgramCounter++
	case "cles":
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

	default:
		log.Fatal()
	}

	return nil
}

type Instruction struct {
	Opcode   string
	Operands []interface{}
}

func immediateOrLabel(s string) interface{} {
	i, err := strconv.Atoi(s)

	if err != nil {
		return Label(s)
	}

	return i
}

func parseInstruction(fields []string) (Instruction, error) {
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
		return Instruction{}, ErrUndefinedOpcode
	}

	return Instruction{
		Opcode:   opcode,
		Operands: operands,
	}, nil
}

func (m *Machine) Load(program string) error {
	section := "text"
	var nextLabel Label

	for _, line := range strings.Split(program, "\n") {
		for _, c := range []string{"(", ")", ","} {
			line = strings.ReplaceAll(line, c, " ")
		}

		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		if strings.HasPrefix(fields[0], ".") {
			section = strings.TrimPrefix(fields[0], ".")
			continue
		}

		switch section {
		case "data":
			if len(fields) != 3 {
				log.Fatal("parse error")
			}

			switch fields[1] {
			case ".float":
				f, err := strconv.ParseFloat(fields[2], 32)
				if err != nil {
					log.Fatal("invalid data")
				}
				m.Memory = append(m.Memory, struct {
					Label Label
					Value interface{}
				}{Label(strings.TrimSuffix(fields[0], ":")), f})
			default:
				log.Fatal("invalid data type")
			}
		case "text":
			if strings.HasSuffix(fields[0], ":") {
				nextLabel = Label(strings.TrimSuffix(fields[0], ":"))
				continue
			}

			instruction, err := parseInstruction(fields)

			if err != nil {
				return err
			}

			m.Memory = append(m.Memory, struct {
				Label Label
				Value interface{}
			}{nextLabel, instruction})

			nextLabel = ""
		default:
			log.Fatal("invalid section")
		}
	}

	// Replaces labels with address values.
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

func (m *Machine) AllocateMemory(heapSize, stackSize int) {
	m.IntRegisters[hp] = len(m.Memory)
	for i := 0; i < heapSize; i++ {
		m.Memory = append(m.Memory, struct {
			Label Label
			Value interface{}
		}{})
	}
	m.IntRegisters[sp] = len(m.Memory)
	for i := 0; i < heapSize; i++ {
		m.Memory = append(m.Memory, struct {
			Label Label
			Value interface{}
		}{})
	}
}

func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := NewMachine()
	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}
	m.AllocateMemory(1000, 1000)
	m.ProgramCounter, err = m.FindAddress("min_caml_start")
	if err != nil {
		log.Fatal(err)
	}
}
