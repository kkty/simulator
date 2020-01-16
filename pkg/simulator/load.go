package simulator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	intRegisters = map[string]int32{
		"$zero": 0, "$tmp0": 1, "$tmp1": 2, "$tmp2": 3,
		"$hp": 29, "$sp": 30, "$ra": 31,
	}

	floatRegisters = map[string]int32{
		"$fzero": 0, "$ftmp": 1,
	}
)

func init() {
	for i := int32(0); i < 25; i++ {
		intRegisters[fmt.Sprintf("$i%d", i)] = i + 4
	}

	for i := int32(0); i < 30; i++ {
		floatRegisters[fmt.Sprintf("$f%d", i)] = i + 2
	}
}

func intRegister(s string) int32 {
	if _, exists := intRegisters[s]; !exists {
		log.Fatalf("unknown register: %s", s)
	}

	return intRegisters[s]
}

func floatRegister(s string) int32 {
	if _, exists := floatRegisters[s]; !exists {
		log.Fatalf("unknown register: %s", s)
	}

	return floatRegisters[s]
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
	case "SLLV":
		operands = append(operands,
			intRegister(fields[1]), intRegister(fields[2]), intRegister(fields[3]))
	case "ADDS":
		fallthrough
	case "SUBS":
		fallthrough
	case "MULS":
		fallthrough
	case "DIVS":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]), floatRegister(fields[3]))
	case "SLL":
		fallthrough
	case "ADDI":
		fallthrough
	case "LUI":
		fallthrough
	case "ORI":
		fallthrough
	case "BEQ":
		operands = append(operands,
			intRegister(fields[1]), intRegister(fields[2]), immediateOrLabel(fields[3]))
	case "BEQS":
		fallthrough
	case "BLS":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]), immediateOrLabel(fields[3]))
	case "BZS":
		operands = append(operands,
			floatRegister(fields[1]), immediateOrLabel(fields[2]))
	case "LW":
		fallthrough
	case "SW":
		operands = append(operands,
			intRegister(fields[1]), immediateOrLabel(fields[2]), intRegister(fields[3]))
	case "LWC1":
		fallthrough
	case "SWC1":
		operands = append(operands,
			floatRegister(fields[1]), immediateOrLabel(fields[2]), intRegister(fields[3]))
	case "SQRT":
		fallthrough
	case "FTOI":
		fallthrough
	case "ITOF":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]))
	case "J":
		fallthrough
	case "JAL":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "JR":
		fallthrough
	case "JALR":
		fallthrough
	case "OUT":
		fallthrough
	case "IN":
		operands = append(operands, intRegister(fields[1]))
	case "INF":
		operands = append(operands, floatRegister(fields[1]))
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
