package simulator

import (
	"errors"
	"fmt"
	"math"
	"os"
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

func uint32ToFloat32(i uint32) float32 {
	return math.Float32frombits(i)
}

func float32ToUint32(f float32) uint32 {
	return math.Float32bits(f)
}

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step(native bool) (bool, error) {
	i, ok := m.memory[m.ProgramCounter].Value.(Instruction)

	if !ok {
		return false, errors.New("no instruction on memory")
	}

	switch opcode := i.Opcode; opcode {
	case "ADD":
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(
			uint32ToInt32(m.Registers[i.Operands[1].(int32)]) + uint32ToInt32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "SUB":
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(
			uint32ToInt32(m.Registers[i.Operands[1].(int32)]) - uint32ToInt32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "SEQ":
		if m.Registers[i.Operands[1].(int32)] == m.Registers[i.Operands[2].(int32)] {
			m.Registers[i.Operands[0].(int32)] = 1
		} else {
			m.Registers[i.Operands[0].(int32)] = 0
		}
		m.ProgramCounter++
	case "SLT":
		if uint32ToInt32(m.Registers[i.Operands[1].(int32)]) < uint32ToInt32(m.Registers[i.Operands[2].(int32)]) {
			m.Registers[i.Operands[0].(int32)] = 1
		} else {
			m.Registers[i.Operands[0].(int32)] = 0
		}
		m.ProgramCounter++
	case "SLTS":
		if uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) < uint32ToFloat32(m.Registers[i.Operands[2].(int32)]) {
			m.Registers[i.Operands[0].(int32)] = 1
		} else {
			m.Registers[i.Operands[0].(int32)] = 0
		}
		m.ProgramCounter++
	case "SLL":
		m.Registers[i.Operands[0].(int32)] = m.Registers[i.Operands[1].(int32)] << i.Operands[2].(int32)
		m.ProgramCounter++
	case "SRL":
		m.Registers[i.Operands[0].(int32)] = m.Registers[i.Operands[1].(int32)] >> i.Operands[2].(int32)
		m.ProgramCounter++
	case "SLLV":
		m.Registers[i.Operands[0].(int32)] = m.Registers[i.Operands[1].(int32)] << m.Registers[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "SRLV":
		m.Registers[i.Operands[0].(int32)] = m.Registers[i.Operands[1].(int32)] >> m.Registers[i.Operands[2].(int32)]
		m.ProgramCounter++
	case "ADDS":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(
			uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) + uint32ToFloat32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "SUBS":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(
			uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) - uint32ToFloat32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "MULS":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(
			uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) * uint32ToFloat32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "DIVS":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(
			uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) / uint32ToFloat32(m.Registers[i.Operands[2].(int32)]))
		m.ProgramCounter++
	case "SQRT":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(
			float32(math.Sqrt(float64(uint32ToFloat32(m.Registers[i.Operands[1].(int32)])))))
		m.ProgramCounter++
	case "FTOI":
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(int32(math.Round(float64(uint32ToFloat32(m.Registers[i.Operands[1].(int32)])))))
		m.ProgramCounter++
	case "ITOF":
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(float32(uint32ToInt32(m.Registers[i.Operands[1].(int32)])))
		m.ProgramCounter++
	case "ADDI":
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(
			uint32ToInt32(m.Registers[i.Operands[1].(int32)]) + i.Operands[2].(int32))
		m.ProgramCounter++
	case "LUI":
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(i.Operands[2].(int32))<<16 | m.Registers[i.Operands[1].(int32)]
		m.ProgramCounter++
	case "ORI":
		m.Registers[i.Operands[0].(int32)] = m.Registers[i.Operands[1].(int32)] | int32ToUint32(i.Operands[2].(int32))
		m.ProgramCounter++
	case "BEQ":
		if m.Registers[i.Operands[0].(int32)] == m.Registers[i.Operands[1].(int32)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "BLT":
		if uint32ToInt32(m.Registers[i.Operands[0].(int32)]) < uint32ToInt32(m.Registers[i.Operands[1].(int32)]) {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "BLTS":
		if uint32ToFloat32(m.Registers[i.Operands[0].(int32)]) < uint32ToFloat32(m.Registers[i.Operands[1].(int32)]) {
			m.ProgramCounter += i.Operands[2].(int32) + 1
		} else {
			m.ProgramCounter++
		}
	case "J":
		m.ProgramCounter = i.Operands[0].(int32)
	case "JAL":
		m.Registers[registers["$ra"]] = int32ToUint32(m.ProgramCounter + 1)
		m.ProgramCounter = i.Operands[0].(int32)
	case "JR":
		m.ProgramCounter = uint32ToInt32(m.Registers[i.Operands[0].(int32)])
	case "LW":
		address := uint32(i.Operands[1].(int32)) + m.Registers[i.Operands[2].(int32)] + m.Registers[i.Operands[3].(int32)]
		m.Registers[i.Operands[0].(int32)] = m.memory[address].Value.(uint32)
		m.ProgramCounter++
	case "SW":
		address := uint32(i.Operands[1].(int32)) + m.Registers[i.Operands[2].(int32)] + m.Registers[i.Operands[3].(int32)]
		m.memory[address].Value = m.Registers[i.Operands[0].(int32)]
		m.ProgramCounter++
	case "OUT":
		os.Stdout.Write([]byte{byte(m.Registers[i.Operands[0].(int32)])})
		m.ProgramCounter++
	case "NOP":
		m.ProgramCounter++
	case "IN":
		var value int32
		fmt.Scan(&value)
		m.Registers[i.Operands[0].(int32)] = int32ToUint32(value)
		m.ProgramCounter++
	case "INF":
		var value float32
		fmt.Scan(&value)
		m.Registers[i.Operands[0].(int32)] = float32ToUint32(value)
		m.ProgramCounter++
	case "EXIT":
		return true, nil
	default:
		return false, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return false, nil
}

func (m *Machine) Run(native bool) (int, error) {
	executed := 0

	for {
		done, err := m.Step(native)

		if err != nil {
			return executed, err
		}

		executed++

		if done {
			return executed, nil
		}
	}
}
