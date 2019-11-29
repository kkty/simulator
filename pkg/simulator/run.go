package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
)

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i, ok := m.memory[m.ProgramCounter].Value.(Instruction)

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
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] + i.Operands[2].(int32)
		m.ProgramCounter++
	case "subi":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] - i.Operands[2].(int32)
		m.ProgramCounter++
	case "lui":
		m.IntRegisters[i.Operands[0].(string)] = i.Operands[1].(int32) << 16
		m.ProgramCounter++
	case "ori":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] | i.Operands[2].(int32)
		m.ProgramCounter++
	case "slt":
		if m.IntRegisters[i.Operands[1].(string)] < m.IntRegisters[i.Operands[2].(string)] {
			m.IntRegisters[i.Operands[0].(string)] = 1
		} else {
			m.IntRegisters[i.Operands[0].(string)] = 0
		}
		m.ProgramCounter++
	case "sll":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] << i.Operands[2].(int32)
		m.ProgramCounter++
	case "sllv":
		m.IntRegisters[i.Operands[0].(string)] = m.IntRegisters[i.Operands[1].(string)] << m.IntRegisters[i.Operands[2].(string)]
		m.ProgramCounter++
	case "j":
		m.ProgramCounter = i.Operands[0].(int32)
	case "jal":
		m.IntRegisters[raRegisterName] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int32)
	case "jr":
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(string)]
	case "jalr":
		m.IntRegisters[raRegisterName] = m.ProgramCounter + 1
		m.ProgramCounter = m.IntRegisters[i.Operands[0].(string)]
	case "beq":
		if m.IntRegisters[i.Operands[0].(string)] == m.IntRegisters[i.Operands[1].(string)] {
			m.ProgramCounter += i.Operands[2].(int32) + 1
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
			m.ProgramCounter += 1 + i.Operands[0].(int32)
		} else {
			m.ProgramCounter++
		}
	case "lw":
		value := m.memory[i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(string)]].Value

		if v, ok := value.(int32); ok {
			m.IntRegisters[i.Operands[0].(string)] = v
		} else if v, ok := value.(float32); ok {
			u := math.Float32bits(v)

			if u >= (1 << 31) {
				m.IntRegisters[i.Operands[0].(string)] = -int32((^u) + 1)
			} else {
				m.IntRegisters[i.Operands[0].(string)] = int32(u)
			}
		} else {
			return false, errors.New("invalid data on memory")
		}
		m.ProgramCounter++
	case "lwc1":
		value := m.memory[i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(string)]].Value

		if v, ok := value.(float32); ok {
			m.FloatRegisters[i.Operands[0].(string)] = v
		} else if v, ok := value.(int32); ok {
			if v >= 0 {
				m.FloatRegisters[i.Operands[0].(string)] = math.Float32frombits(uint32(v))
			} else {
				m.FloatRegisters[i.Operands[0].(string)] = math.Float32frombits((^uint32(-v)) + 1)
			}
		} else {
			return false, errors.New("invalid data on memory")
		}
		m.ProgramCounter++
	case "sw":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(string)], m.IntRegisters[i.Operands[0].(string)])
		m.ProgramCounter++
	case "swc1":
		m.setValueToMemory(i.Operands[1].(int32)+m.IntRegisters[i.Operands[2].(string)], m.FloatRegisters[i.Operands[0].(string)])
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
	case "sqrt":
		m.FloatRegisters[i.Operands[0].(string)] = float32(math.Sqrt(float64(m.FloatRegisters[i.Operands[1].(string)])))
		m.ProgramCounter++
	case "ftoi":
		m.IntRegisters[i.Operands[0].(string)] = int32(math.Round(float64(m.FloatRegisters[i.Operands[1].(string)])))
		m.ProgramCounter++
	case "itof":
		m.FloatRegisters[i.Operands[0].(string)] = float32(m.IntRegisters[i.Operands[1].(string)])
		m.ProgramCounter++
	case "read_i":
		var v int32
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

// Stats can be used to count the number of executed instructions and
// the number of jumps to each label.
type Stats struct {
	Executed int
	jumps    map[Label]int
}

// LabelWithCount is used to store a label and the number of times it was called.
type LabelWithCount struct {
	Label Label
	Count int
}

// Jumps fetches the labels that were most frequently jumped to.
func (s *Stats) Jumps(count int) []LabelWithCount {
	l := []LabelWithCount{}

	for label, count := range s.jumps {
		l = append(l, LabelWithCount{label, count})
	}

	sort.Slice(l, func(i, j int) bool { return l[i].Count > l[j].Count })

	if len(l) <= count {
		return l
	} else {
		return l[:count]
	}
}

func (m *Machine) Run(debug bool) (Stats, error) {
	stats := Stats{
		jumps: make(map[Label]int),
	}

	if debug {
		// Prints the initial memory state.
		memory := m.Memory()

		b, err := json.Marshal(memory)

		if err != nil {
			return stats, err
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
				"instruction":    m.memory[m.ProgramCounter],
			})

			if err != nil {
				return stats, err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(b))
		}

		if label := m.memory[m.ProgramCounter].Label; label != "" {
			stats.jumps[label] += 1
		}

		done, err := m.Step()

		if err != nil {
			return stats, err
		}

		stats.Executed++

		if done {
			return stats, nil
		}
	}
}
