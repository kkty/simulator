package simulator

import (
	"errors"
	"fmt"
	"math"
)

func uint32ToInt32(i uint32) int32 {
	if i >= (1 << 31) {
		return -int32((^i) + 1)
	} else {
		return int32(i)
	}
}

func int32ToUint32(i int32) uint32 {
	if i >= 0 {
		return uint32(i)
	} else {
		return (^uint32(-i)) + 1
	}
}

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i, ok := m.memory[m.ProgramCounter].Value.(Instruction)

	if !ok {
		return false, errors.New("no instruction on memory")
	}

	switch opcode := i.Opcode; opcode {
	case "ADD":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] + m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "SUB":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] - m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "ADDI":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] + i.Operands[2].(int32)
		m.ProgramCounter++
	case "SUBI":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] - i.Operands[2].(int32)
		m.ProgramCounter++
	case "LUI":
		m.IntRegisters[i.Operands[0].(int32)] = uint32ToInt32(int32ToUint32(i.Operands[2].(int32))<<16 | int32ToUint32(m.IntRegisters[i.Operands[1].(int32)]))
		m.ProgramCounter++
	case "ORI":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] | i.Operands[2].(int32)
		m.ProgramCounter++
	case "SLT":
		if m.IntRegisters[i.Operands[1].(int32)] < m.IntRegisters[i.Operands[2].(int32)] {
			m.IntRegisters[i.Operands[0].(int32)] = 1
		} else {
			m.IntRegisters[i.Operands[0].(int32)] = 0
		}
		m.ProgramCounter++
	case "SLL":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] << i.Operands[2].(int32)
		m.ProgramCounter++
	case "SLLV":
		m.IntRegisters[i.Operands[0].(int32)] = m.IntRegisters[i.Operands[1].(int32)] << m.IntRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "J":
		m.ProgramCounter = i.Operands[0].(int32)
	case "JAL":
		m.IntRegisters[31] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int32)
	case "JR":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(int32)]
	case "JALR":
		m.IntRegisters[31] = m.ProgramCounter + 1
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(int32)]
	case "BEQ":
		if m.IntRegisters[i.Operands[0].(int32)] == m.IntRegisters[i.Operands[1].(int32)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "BEQS":
		if m.FloatRegisters[i.Operands[0].(int32)] == m.FloatRegisters[i.Operands[1].(int32)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "BLS":
		if m.FloatRegisters[i.Operands[0].(int32)] < m.FloatRegisters[i.Operands[1].(int32)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "LW":
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
	case "LWC1":
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
	case "SW":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(int32)], m.IntRegisters[i.Operands[0].(int32)])
		m.ProgramCounter++
	case "SWC1":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(int32)], m.FloatRegisters[i.Operands[0].(int32)])
		m.ProgramCounter++
	case "ADDS":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] + m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "SUBS":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] - m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "MULS":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] * m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "DIVS":
		m.FloatRegisters[i.Operands[0].(int32)] = m.FloatRegisters[i.Operands[1].(int32)] / m.FloatRegisters[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "SQRT":
		m.FloatRegisters[i.Operands[0].(int32)] = float32(math.Sqrt(float64(m.FloatRegisters[i.Operands[1].(int32)])))
		m.ProgramCounter++
	case "FTOI":
		v := int32(math.Round(float64(m.FloatRegisters[i.Operands[1].(int32)])))

		if v >= 0 {
			m.FloatRegisters[i.Operands[0].(int32)] = math.Float32frombits(uint32(v))
		} else {
			m.FloatRegisters[i.Operands[0].(int32)] = math.Float32frombits((^uint32(-v)) + 1)
		}

		m.ProgramCounter++
	case "ITOF":
		u := math.Float32bits(m.FloatRegisters[i.Operands[1].(int32)])

		if u >= (1 << 31) {
			m.FloatRegisters[i.Operands[0].(int32)] = float32(-int32((^u) + 1))
		} else {
			m.FloatRegisters[i.Operands[0].(int32)] = float32(int32(u))
		}

		m.ProgramCounter++
	case "OUT":
		fmt.Print(string(m.IntRegisters[i.Operands[0].(int32)]))
		m.ProgramCounter++
	case "NOP":
		m.ProgramCounter++
	case "EXIT":
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
