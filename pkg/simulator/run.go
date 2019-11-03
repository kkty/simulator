package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i, ok := m.Memory[m.ProgramCounter].Value.(Instruction)

	if !ok {
		return false, errors.New("no instruction on memory")
	}

	defer func() {
		// The zero register should always be zero.
		m.IntRegisters[zeroRegisterName] = 0
	}()

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
		m.IntRegisters[raRegisterName] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int)
	case "jr":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(string)]
	case "jalr":
		m.IntRegisters[raRegisterName] = m.ProgramCounter + 1
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
			m.ProgramCounter += 1 + i.Operands[0].(int)
		} else {
			m.ProgramCounter++
		}
	case "lw":
		m.IntRegisters[i.Operands[0].(string)] = m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value.(int)
		m.ProgramCounter++
	case "lwc1":
		m.FloatRegisters[i.Operands[0].(string)] = m.Memory[i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)]].Value.(float32)
		m.ProgramCounter++
	case "sw":
		m.setValueToMemory(i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)], m.IntRegisters[i.Operands[0].(string)])
		m.ProgramCounter++
	case "swc1":
		m.setValueToMemory(i.Operands[1].(int)+m.IntRegisters[i.Operands[2].(string)], m.FloatRegisters[i.Operands[0].(string)])
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
	case "read_i":
		var v int
		_, err := fmt.Scanf("%d", &v)
		if err != nil {
			return false, err
		}
		m.IntRegisters[i.Operands[0].(string)] = v
		m.ProgramCounter++
	case "read_f":
		var v float32
		_, err := fmt.Scanf("%f", &v)
		if err != nil {
			return false, err
		}
		m.FloatRegisters[i.Operands[0].(string)] = v
		m.ProgramCounter++
	case "out_i":
		fmt.Print(m.IntRegisters[i.Operands[0].(string)])
		m.ProgramCounter++
	case "out_c":
		fmt.Print(string(m.IntRegisters[i.Operands[0].(string)]))
		m.ProgramCounter++
	case "nop":
		m.ProgramCounter++
	case "exit":
		return true, nil
	default:
		return false, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return false, nil
}

func (m *Machine) Run(debug bool) error {
	if debug {
		// Prints the initial memory state.

		memory := []struct {
			Address int
			Value   interface{}
			Label   Label
		}{}

		for address, valueWithLabel := range m.Memory {
			memory = append(memory, struct {
				Address int
				Value   interface{}
				Label   Label
			}{address, valueWithLabel.Value, valueWithLabel.Label})
		}

		sort.Slice(memory, func(i, j int) bool { return memory[i].Address < memory[j].Address })

		b, err := json.Marshal(memory)

		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%s\n", string(b))
	}

	for {
		if debug {
			// Prints the current machine state.
			b, err := json.Marshal(map[string]interface{}{
				"programCounter": m.ProgramCounter,
				"intRegisters":   m.IntRegisters,
				"floatRegisters": m.FloatRegisters,
				"instruction":    m.Memory[m.ProgramCounter],
			})

			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(b))
		}

		done, err := m.Step()

		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}
