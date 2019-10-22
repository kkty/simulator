package simulator

import "fmt"

// Step fetches an instruction and executes it.
// Returns true if it has encountered "exit" call.
func (m *Machine) Step() (bool, error) {
	i := m.Memory[m.ProgramCounter].Value.(Instruction)

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
		m.IntRegisters["$ra"] = m.ProgramCounter + 1
		m.ProgramCounter = i.Operands[0].(int)
	case "jr":
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
			m.ProgramCounter += 1 + i.Operands[1].(int)
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
	case "exit":
		return true, nil
	default:
		return false, fmt.Errorf("%v: invalid opcode", opcode)
	}

	return false, nil
}

func (m *Machine) Run() error {
	for {
		done, err := m.Step()

		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}
