package simulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAddress(t *testing.T) {
	m := Machine{
		Memory: []ValueWithLabel{
			{Label: ""},
			{Label: "l1"},
			{Label: "l2"},
		},
	}

	for _, c := range []struct {
		label    Label
		expected int
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
