package simulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAddress(t *testing.T) {
	m := NewMachine()
	m.Memory = map[int32]ValueWithLabel{
		1: ValueWithLabel{Label("l1"), 10},
		2: ValueWithLabel{Label("l2"), 20},
	}

	for _, c := range []struct {
		label    Label
		expected int32
	}{
		{Label("l1"), 1},
		{Label("l2"), 2},
	} {
		actual, err := m.FindAddress(c.label)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, c.expected, actual)
	}

}
