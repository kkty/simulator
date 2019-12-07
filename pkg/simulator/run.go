package simulator

import (
	"errors"
	"fmt"
	"math"
)

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i, ok := m.memory[m.ProgramCounter].Value.(Instruction)

	if !ok {
		return false, errors.New("no instruction on memory")
	}

	switch opcode := i.Opcode; opcode {
	case "add":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] + m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "sub":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] - m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "addi":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] + i.Operands[2].(int32)
		m.ProgramCounter++
	case "subi":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] - i.Operands[2].(int32)
		m.ProgramCounter++
	case "lui":
		m.IntRegisters[i.Operands[0].(int32)] = i.Operands[1].(int32) << 16
		m.ProgramCounter++
	case "ori":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] | i.Operands[2].(int32)
		m.ProgramCounter++
	case "slt":
		if m.IntRegisters[i.Operands[1].(int32)] < m.IntRegisters[i.Operands[2].(int32)] {
			m.IntRegisters[i.Operands[0].(int32)] = 1
		} else {
			m.IntRegisters[i.Operands[0].(int32)] = 0
		}
		m.ProgramCounter++
	case "sll":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] << i.Operands[2].(int32)
		m.ProgramCounter++
	case "sllv":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] << m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "j":
		m.ProgramCounter = i.Operands[0].(int32)
	case "jal":
		m.IntRegisters[31] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int32)
	case "jr":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(int32)]
	case "jalr":
		m.IntRegisters[31] = m.ProgramCounter + 1
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(int32)]
	case "beq":
		if m.IntRegisters[i.Operands[0].(int32)] == m.IntRegisters[i.Operands[1].(int32)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "c.eq.s":
		if m.FloatRegisters[i.Operands[0].(int32)] == m.FloatRegisters[i.Operands[1].(int32)] {
			m.ConditionRegister = true
		} else {
			m.ConditionRegister = false
		}

		m.ProgramCounter++
	case "c.le.s":
		if m.FloatRegisters[i.Operands[0].(int32)] <= m.FloatRegisters[i.Operands[1].(int32)] {
			m.ConditionRegister = true
		} else {
			m.ConditionRegister = false
		}

		m.ProgramCounter++
	case "bc1t":
		if m.ConditionRegister {
			m.ProgramCounter += 1 + i.Operands[0].(int32)
		} else {
			m.ProgramCounter++
		}
	case "lw":
		address := i.Operands[1].(int32) + m.IntRegisters[i.Operands[2].(int32)]
		value := m.memory[address].Value

		if v, ok := value.(int32); ok {
			m.IntRegisters[i.Operands[0].(int32)] = v
		} else if v, ok := value.(float32); ok {
			u := math.Float32bits(v)

			if u >= (1 << 31) {
				m.IntRegisters[i.Operands[0].(int32)] = -int32((^u) + 1)
			} else {
				m.IntRegisters[i.Operands[0].(int32)] = int32(u)
			}
		} else {
			return false, errors.New("invalid data on memory")
		}
		m.ProgramCounter++
	case "lwc1":
		address := i.Operands[1].(int32) + m.IntRegisters[i.Operands[2].(int32)]
		value := m.memory[address].Value

		if v, ok := value.(float32); ok {
			m.FloatRegisters[i.Operands[0].(int32)] = v
		} else if v, ok := value.(int32); ok {
			if v >= 0 {
				m.FloatRegisters[i.Operands[0].(int32)] = math.Float32frombits(uint32(v))
			} else {
				m.FloatRegisters[i.Operands[0].(int32)] = math.Float32frombits((^uint32(-v)) + 1)
			}
		} else {
			return false, errors.New("invalid data on memory")
		}
		m.ProgramCounter++
	case "sw":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(int32)], m.IntRegisters[i.Operands[0].(int32)])
		m.ProgramCounter++
	case "swc1":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(int32)], m.FloatRegisters[i.Operands[0].(int32)])
		m.ProgramCounter++
	case "add.s":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] + m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "sub.s":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] - m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "mul.s":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] * m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "div.s":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] / m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "sqrt":
		m.FloatRegisters[i.Operands[0].(int32)] = float32(math.Sqrt(float64(m.FloatRegisters[i.Operands[1].(int32)])))
		m.ProgramCounter++
	case "ftoi":
		m.IntRegisters[i.Operands[0].(int32)] = int32(math.Round(float64(m.FloatRegisters[i.Operands[1].(int32)])))
		m.ProgramCounter++
	case "itof":
		m.FloatRegisters[i.Operands[0].(int32)] = float32(m.IntRegisters[i.Operands[1].(int32)])
		m.ProgramCounter++
	case "read_i":
		var v int32
		_, err := fmt.Scanf("%d", &v)
		if err != nil {
			return false, err
		}
		m.IntRegisters[i.Operands[0].(int32)] = v
		m.ProgramCounter++
	case "read_f":
		var v float32
		_, err := fmt.Scanf("%f", &v)
		if err != nil {
			return false, err
		}
		m.FloatRegisters[i.Operands[0].(int32)] = v
		m.ProgramCounter++
	case "out_i":
		fmt.Print(m.IntRegisters[i.Operands[0].(int32)])
		m.ProgramCounter++
	case "out_c":
		fmt.Print(string(m.IntRegisters[i.Operands[0].(int32)]))
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

func (m *Machine) Run() (int, error) {
	executed := 0

	for {
		done, err := m.Step()

		if err != nil {
			return executed, err
		}

		executed++

		if done {
			return executed, nil
		}
	}
}
