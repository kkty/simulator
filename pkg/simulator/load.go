package simulator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	registers = map[string]int32{
		"$zero": 60, "$hp": 61, "$sp": 62, "$ra": 63,
	}
)

func init() {
	for i := int32(0); i < 60; i++ {
		registers[fmt.Sprintf("$r%d", i)] = i
	}
}

func register(s string) int32 {
	if i, exists := registers[s]; exists {
		return i
	}

	panic(fmt.Sprintf("unknown register: %s", s))
}

func parseInstruction(fields []string) (Instruction, error) {
	immediateOrLabel := func(s string) interface{} {
		i, err := strconv.Atoi(s)

		if err != nil {
			return Label(s)
		}

		return int32(i)
	}

	opcode := fields[0]
	var operands []interface{}

	switch opcode {
	case "ADD":
		fallthrough
	case "SUB":
		fallthrough
	case "SLT":
		fallthrough
	case "SLTS":
		fallthrough
	case "SEQ":
		fallthrough
	case "SLLV":
		fallthrough
	case "SRLV":
		fallthrough
	case "ADDS":
		fallthrough
	case "SUBS":
		fallthrough
	case "MULS":
		fallthrough
	case "DIVS":
		operands = append(operands,
			register(fields[1]), register(fields[2]), register(fields[3]))
	case "SLL":
		fallthrough
	case "SRL":
		fallthrough
	case "ADDI":
		fallthrough
	case "LUI":
		fallthrough
	case "ORI":
		fallthrough
	case "BLT":
		fallthrough
	case "BLTS":
		fallthrough
	case "BEQ":
		operands = append(operands,
			register(fields[1]), register(fields[2]), immediateOrLabel(fields[3]))
	case "LW":
		fallthrough
	case "SW":
		operands = append(operands,
			register(fields[1]), immediateOrLabel(fields[2]), register(fields[3]))
	case "SQRT":
		fallthrough
	case "FTOI":
		fallthrough
	case "ITOF":
		operands = append(operands,
			register(fields[1]), register(fields[2]))
	case "J":
		fallthrough
	case "JAL":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "JR":
		fallthrough
	case "OUT":
		fallthrough
	case "IN":
		operands = append(operands, register(fields[1]))
	case "INF":
		operands = append(operands, register(fields[1]))
	case "EXIT":
	case "NOP":
	default:
		return Instruction{}, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return Instruction{
		Opcode:   opcode,
		Operands: operands,
	}, nil
}

// TODO: remove this
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
			float32(f),
		}, nil
	default:
		return ValueWithLabel{}, errors.New("invalid data type")
	}
}

// Load loads a program onto the memory.
func (m *Machine) Load(program string) error {
	var nextLabel Label
	var nextAddress int32

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

		if strings.HasSuffix(fields[0], ":") {
			nextLabel = Label(strings.TrimSuffix(fields[0], ":"))
			continue
		}

		instruction, err := parseInstruction(fields)

		if err != nil {
			return wrapError(err)
		}

		m.memory[nextAddress] = ValueWithLabel{nextLabel, instruction}
		nextAddress++

		nextLabel = ""
	}

	// Iterates thorough the memory and replaces labels with address values.
	for i := int32(0); i < nextAddress; i++ {
		if instruction, ok := m.memory[i].Value.(Instruction); ok {
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
