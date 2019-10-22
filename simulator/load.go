package simulator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	case "exit":
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
