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
	case "add":
		fallthrough
	case "sub":
		fallthrough
	case "slt":
		fallthrough
	case "sllv":
		operands = append(operands,
			intRegister(fields[1]), intRegister(fields[2]), intRegister(fields[3]))
	case "add.s":
		fallthrough
	case "sub.s":
		fallthrough
	case "mul.s":
		fallthrough
	case "div.s":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]), floatRegister(fields[3]))
	case "sll":
		fallthrough
	case "addi":
		fallthrough
	case "lui":
		fallthrough
	case "ori":
		fallthrough
	case "beq":
		operands = append(operands,
			intRegister(fields[1]), intRegister(fields[2]), immediateOrLabel(fields[3]))
	case "beqs":
		fallthrough
	case "bls":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]), immediateOrLabel(fields[3]))
	case "lw":
		fallthrough
	case "sw":
		operands = append(operands,
			intRegister(fields[1]), immediateOrLabel(fields[2]), intRegister(fields[3]))
	case "lwc1":
		fallthrough
	case "swc1":
		operands = append(operands,
			floatRegister(fields[1]), immediateOrLabel(fields[2]), intRegister(fields[3]))
	case "bc1t":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "c.eq.s":
		fallthrough
	case "c.le.s":
		fallthrough
	case "sqrt":
		fallthrough
	case "ftoi":
		fallthrough
	case "itof":
		operands = append(operands,
			floatRegister(fields[1]), floatRegister(fields[2]))
	case "j":
		fallthrough
	case "jal":
		operands = append(operands, immediateOrLabel(fields[1]))
	case "jr":
		fallthrough
	case "jalr":
		fallthrough
	case "out_i":
		fallthrough
	case "out_c":
		fallthrough
	case "read_i":
		operands = append(operands, intRegister(fields[1]))
	case "read_f":
		operands = append(operands, floatRegister(fields[1]))
	case "exit":
	case "nop":
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
func (m *Machine) Load(program string, mappedData []byte, mappedAddress int32) error {
	// The default section is "text".
	section := "text"

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

		// Changes section.
		if strings.HasPrefix(fields[0], ".") {
			section = strings.TrimPrefix(fields[0], ".")
			continue
		}

		switch section {
		case "data":
			var err error

			m.memory[nextAddress], err = parseData(fields)
			nextAddress++

			if err != nil {
				return wrapError(err)
			}
		case "text":
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
		default:
			return wrapError(fmt.Errorf("%v: invalid section", section))
		}
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

	for i, c := range mappedData {
		i := int32(i)
		if i%4 == 0 {
			m.memory[mappedAddress+(i/4)*4].Value = int32(0)
		}

		m.memory[mappedAddress+(i/4)*4].Value = uint32ToInt32(int32ToUint32(m.memory[mappedAddress+(i/4)*4].Value.(int32)) + uint32(c)<<(24-(i%4)*8))
	}

	return nil
}
