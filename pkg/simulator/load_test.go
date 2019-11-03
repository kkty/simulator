package simulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseData(t *testing.T) {
	valueWithLabel, err := parseData([]string{"l1", ".float", "0.25"})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Label("l1"), valueWithLabel.Label)
	assert.Equal(t, float32(0.25), valueWithLabel.Value)
}
func TestParseInstruction(t *testing.T) {
	for _, c := range []struct {
		fields   []string
		expected Instruction
	}{
		{
			[]string{"add", "$i0", "$i1", "$i2"},
			Instruction{"add", []interface{}{"$i0", "$i1", "$i2"}},
		},
		{
			[]string{"addi", "$i0", "$i1", "1"},
			Instruction{"addi", []interface{}{"$i0", "$i1", int32(1)}},
		},
		{
			[]string{"addi", "$i0", "$i1", "l1"},
			Instruction{"addi", []interface{}{"$i0", "$i1", Label("l1")}},
		},
	} {
		actual, err := parseInstruction(c.fields)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, c.expected, actual)
	}
}
