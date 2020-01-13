package float

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const eps = 1e-6
const rep = 10000

func generate() float32 {
	return rand.Float32()*1000 - 500
}

func TestAdd(t *testing.T) {
	for i := 0; i < rep; i++ {
		x1, x2 := generate(), generate()
		t.Run(fmt.Sprintf("x1=%.2f,x2=%.2f", x1, x2), func(t *testing.T) {
			assert.InEpsilon(t, x1+x2, Add(NewFromFloat32(x1), NewFromFloat32(x2)).Float32(), eps)
		})
	}
}

func TestDiv(t *testing.T) {
	for i := 0; i < rep; i++ {
		x1, x2 := generate(), generate()
		t.Run(fmt.Sprintf("x1=%.2f,x2=%.2f", x1, x2), func(t *testing.T) {
			assert.InEpsilon(t, x1/x2, Div(NewFromFloat32(x1), NewFromFloat32(x2)).Float32(), eps)
		})
	}
}

func TestMul(t *testing.T) {
	for i := 0; i < rep; i++ {
		x1, x2 := generate(), generate()
		t.Run(fmt.Sprintf("x1=%.2f,x2=%.2f", x1, x2), func(t *testing.T) {
			assert.InEpsilon(t, x1*x2, Mul(NewFromFloat32(x1), NewFromFloat32(x2)).Float32(), eps)
		})
	}
}

func TestSub(t *testing.T) {
	for i := 0; i < rep; i++ {
		x1, x2 := generate(), generate()
		t.Run(fmt.Sprintf("x1=%.2f,x2=%.2f", x1, x2), func(t *testing.T) {
			assert.InEpsilon(t, x1-x2, Sub(NewFromFloat32(x1), NewFromFloat32(x2)).Float32(), eps)
		})
	}
}

func TestFloat(t *testing.T) {
	assert.Equal(t, -3, NewFromInt(6, -3).Int())
	assert.Equal(t, 3, NewFromInt(6, 3).Int())
	assert.Equal(t, uint(3), NewFromUint(6, 3).Uint())
	assert.Equal(t, float32(3.14), NewFromFloat32(3.14).Float32())
}
