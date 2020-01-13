package float

import "math"

type Bits []bool

func NewFromFloat32(data float32) Bits {
	return NewFromUint(32, uint(math.Float32bits(data)))
}

func (b Bits) Float32() float32 {
	return math.Float32frombits(uint32(b.Uint()))
}

func NewFromUint(length int, data uint) Bits {
	ret := []bool{}
	for i := length - 1; i >= 0; i-- {
		ret = append(ret, (data&(uint(1)<<i)) > 0)
	}
	return ret
}

func NewFromInt(length int, data int) Bits {
	if data >= 0 {
		return NewFromUint(length, uint(data))
	}

	return NewFromUint(length, NewFromUint(length, uint(-data)).Reverse().Uint()+1)
}

func (b Bits) Uint() uint {
	data := uint(0)
	for i := 0; i < len(b); i++ {
		if b[len(b)-1-i] {
			data |= uint(1) << i
		}
	}
	return data
}

func (b Bits) Int() int {
	if b[0] == false {
		return int(Bits(b[1:]).Uint())
	}

	return -int(Bits(b[1:]).Reverse().Uint() + 1)
}

func (b Bits) Reverse() Bits {
	ret := []bool{}
	for _, i := range b {
		ret = append(ret, !i)
	}
	return ret
}

func (b Bits) Slice(h, l int) Bits {
	return b[len(b)-1-h : len(b)-l]
}

func (b Bits) Append(i Bits) Bits {
	return append(b, i...)
}

func Add(x1, x2 Bits) Bits {
	s1 := x1.Slice(31, 31)
	e1 := x1.Slice(30, 23)
	m1 := x1.Slice(22, 0)
	s2 := x2.Slice(31, 31)
	e2 := x2.Slice(30, 23)
	m2 := x2.Slice(22, 0)

	var m1a Bits
	if e1.Uint() == 0 {
		m1a = NewFromUint(2, 0).Append(m1)
	} else {
		m1a = NewFromUint(2, 1).Append(m1)
	}

	var m2a Bits
	if e2.Uint() == 0 {
		m2a = NewFromUint(2, 0).Append(m2)
	} else {
		m2a = NewFromUint(2, 1).Append(m2)
	}

	var e1a, e2a Bits
	if e1.Uint() > 0 {
		e1a = e1
	} else {
		e1a = NewFromUint(8, 1)
	}
	if e2.Uint() > 0 {
		e2a = e2
	} else {
		e2a = NewFromUint(8, 1)
	}

	e2ai := e2a.Reverse()
	te := NewFromUint(9, e1a.Uint()+e2ai.Uint())

	var ce, hog, tde Bits
	if te.Slice(8, 8).Uint() > 0 {
		ce = NewFromUint(1, 0)
		hog = NewFromUint(10, te.Uint()+1)
		tde = hog.Slice(7, 0)
	} else {
		ce = NewFromUint(1, 1)
		hog = NewFromUint(1, 0).Append(te.Reverse())
		tde = te.Slice(7, 0).Reverse()
	}

	var de Bits
	if tde.Uint() > 31 {
		de = NewFromUint(5, 31)
	} else {
		de = tde.Slice(4, 0)
	}

	var sel Bits
	if de.Uint() == 0 {
		if m1a.Uint() > m2a.Uint() {
			sel = NewFromUint(1, 0)
		} else {
			sel = NewFromUint(1, 1)
		}
	} else {
		sel = ce
	}

	var ms, mi, es, ss Bits
	if sel.Uint() == 0 {
		ms = m1a
		mi = m2a
		es = e1a
		ss = s1
	} else {
		ms = m2a
		mi = m1a
		es = e2a
		ss = s2
	}

	mie := mi.Append(NewFromUint(31, 0))

	mia := NewFromUint(56, mie.Uint()>>de.Uint())

	var tstck Bits
	if mia.Slice(28, 0).Uint() > 0 {
		tstck = NewFromUint(1, 1)
	} else {
		tstck = NewFromUint(1, 0)
	}

	var mye Bits
	if s1.Uint() == s2.Uint() {
		mye = NewFromUint(27, ms.Append(NewFromUint(2, 0)).Uint()+mia.Slice(55, 29).Uint())
	} else {
		mye = NewFromUint(27, ms.Append(NewFromUint(2, 0)).Uint()-mia.Slice(55, 29).Uint())
	}

	esi := NewFromUint(8, es.Uint()+1)

	var eyd, myd, stck Bits

	if mye.Slice(26, 26).Uint() > 0 {
		if esi.Uint() == 255 {
			eyd = NewFromUint(8, 255)
			myd = NewFromUint(2, 1).Append(NewFromUint(25, 0))
			stck = NewFromUint(1, 0)
		} else {
			eyd = esi
			myd = NewFromUint(27, mye.Uint()>>1)
			if (tstck.Uint() > 0) || (mye.Slice(0, 0).Uint() > 0) {
				stck = NewFromUint(1, 1)
			} else {
				stck = NewFromUint(1, 0)
			}
		}
	} else {
		eyd = es
		myd = mye
		stck = tstck
	}

	var se Bits
	for i := 25; i >= 0; i-- {
		if myd.Slice(i, i).Uint() == 1 {
			se = NewFromUint(5, 25-uint(i))
			break
		}
		if i == 0 {
			se = NewFromUint(5, 26)
		}
	}

	eyf := NewFromUint(9, eyd.Uint()-se.Uint())

	var eyr, myf Bits

	if eyf.Int() > 0 {
		eyr = eyf.Slice(7, 0)
		myf = NewFromUint(27, myd.Uint()<<se.Uint())
	} else {
		eyr = NewFromUint(8, 0)
		myf = NewFromUint(27, myd.Uint()<<(eyd.Slice(4, 0).Uint()-1))
	}

	var myr Bits
	if (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 0 && stck.Slice(0, 0).Uint() == 0 && myf.Slice(2, 2).Uint() == 1) || (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 0 && s1.Uint() == s2.Uint() && stck.Slice(0, 0).Uint() == 1) || (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 1) {
		myr = NewFromUint(25, myf.Slice(26, 2).Uint()+1)
	} else {
		myr = myf.Slice(26, 2)
	}

	eyri := NewFromUint(8, eyr.Uint()+1)

	var my, ey Bits
	if myr.Slice(24, 24).Uint() == 1 {
		my = NewFromUint(23, 0)
		ey = eyri
	} else if myr.Slice(23, 0).Uint() == 0 {
		my = NewFromUint(23, 0)
		ey = NewFromUint(8, 0)
	} else {
		my = myr.Slice(22, 0)
		ey = eyr
	}

	var sy Bits
	if ey.Uint() == 0 && my.Uint() == 0 {
		if s1.Uint() > 0 && s2.Uint() > 0 {
			sy = NewFromUint(1, 1)
		} else {
			sy = NewFromUint(1, 0)
		}
	} else {
		sy = ss
	}

	var nzm1 Bits
	if m1.Slice(22, 0).Uint() > 0 {
		nzm1 = NewFromUint(1, 1)
	} else {
		nzm1 = NewFromUint(1, 0)
	}

	var nzm2 Bits
	if m2.Slice(22, 0).Uint() > 0 {
		nzm2 = NewFromUint(1, 1)
	} else {
		nzm2 = NewFromUint(1, 0)
	}

	var y Bits
	if e1.Uint() == 255 && e2.Uint() != 255 {
		y = s1.Append(NewFromUint(8, 255)).Append(nzm1).Append(m1.Slice(21, 0))
	} else if e2.Uint() == 255 && e1.Uint() != 255 {
		y = s2.Append(NewFromUint(8, 255)).Append(nzm2).Append(m2.Slice(21, 0))
	} else if e2.Uint() == 255 && nzm2.Uint() > 0 {
		y = s2.Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(m2.Slice(21, 0))
	} else if e1.Uint() == 255 && nzm1.Uint() > 0 {
		y = s1.Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(m1.Slice(21, 0))
	} else if e1.Uint() == 255 && e2.Uint() == 255 && s1.Uint() == s2.Uint() {
		y = s1.Append(NewFromUint(8, 255)).Append(NewFromUint(23, 0))
	} else if e1.Uint() == 255 && e2.Uint() == 255 {
		y = NewFromUint(1, 1).Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(NewFromUint(22, 0))
	} else {
		y = sy.Append(ey).Append(my)
	}

	return y
}

func Div(x1, x2 Bits) Bits {
	var s Bits
	if (x1.Slice(31, 31).Uint() > 0) != (x2.Slice(31, 31).Uint() > 0) {
		s = NewFromUint(1, 1)
	} else {
		s = NewFromUint(1, 0)
	}
	e1 := x1.Slice(30, 23)
	e2 := x2.Slice(30, 23)
	eReal := NewFromInt(10, int(e1.Uint())-int(e2.Uint())+127)
	eMoge := eReal.Slice(7, 0)

	m1 := x1.Slice(22, 0)
	m2 := x2.Slice(22, 0)

	var r0 Bits
	switch m2.Slice(22, 14).Uint() {
	case 0b000000000:
		r0 = NewFromUint(9, 0b000000000)
	case 0b000000001:
		r0 = NewFromUint(9, 0b111111111)
	case 0b000000010:
		r0 = NewFromUint(9, 0b111111110)
	case 0b000000011:
		r0 = NewFromUint(9, 0b111111101)
	case 0b000000100:
		r0 = NewFromUint(9, 0b111111100)
	case 0b000000101:
		r0 = NewFromUint(9, 0b111111011)
	case 0b000000110:
		r0 = NewFromUint(9, 0b111111010)
	case 0b000000111:
		r0 = NewFromUint(9, 0b111111001)
	case 0b000001000:
		r0 = NewFromUint(9, 0b111111000)
	case 0b000001001:
		r0 = NewFromUint(9, 0b111110111)
	case 0b000001010:
		r0 = NewFromUint(9, 0b111110110)
	case 0b000001011:
		r0 = NewFromUint(9, 0b111110101)
	case 0b000001100:
		r0 = NewFromUint(9, 0b111110100)
	case 0b000001101:
		r0 = NewFromUint(9, 0b111110011)
	case 0b000001110:
		r0 = NewFromUint(9, 0b111110010)
	case 0b000001111:
		r0 = NewFromUint(9, 0b111110001)
	case 0b000010000:
		r0 = NewFromUint(9, 0b111110000)
	case 0b000010001:
		r0 = NewFromUint(9, 0b111101111)
	case 0b000010010:
		r0 = NewFromUint(9, 0b111101110)
	case 0b000010011:
		r0 = NewFromUint(9, 0b111101101)
	case 0b000010100:
		r0 = NewFromUint(9, 0b111101100)
	case 0b000010101:
		r0 = NewFromUint(9, 0b111101011)
	case 0b000010110:
		r0 = NewFromUint(9, 0b111101010)
	case 0b000010111:
		r0 = NewFromUint(9, 0b111101001)
	case 0b000011000:
		r0 = NewFromUint(9, 0b111101001)
	case 0b000011001:
		r0 = NewFromUint(9, 0b111101000)
	case 0b000011010:
		r0 = NewFromUint(9, 0b111100111)
	case 0b000011011:
		r0 = NewFromUint(9, 0b111100110)
	case 0b000011100:
		r0 = NewFromUint(9, 0b111100101)
	case 0b000011101:
		r0 = NewFromUint(9, 0b111100100)
	case 0b000011110:
		r0 = NewFromUint(9, 0b111100011)
	case 0b000011111:
		r0 = NewFromUint(9, 0b111100010)
	case 0b000100000:
		r0 = NewFromUint(9, 0b111100001)
	case 0b000100001:
		r0 = NewFromUint(9, 0b111100000)
	case 0b000100010:
		r0 = NewFromUint(9, 0b111100000)
	case 0b000100011:
		r0 = NewFromUint(9, 0b111011111)
	case 0b000100100:
		r0 = NewFromUint(9, 0b111011110)
	case 0b000100101:
		r0 = NewFromUint(9, 0b111011101)
	case 0b000100110:
		r0 = NewFromUint(9, 0b111011100)
	case 0b000100111:
		r0 = NewFromUint(9, 0b111011011)
	case 0b000101000:
		r0 = NewFromUint(9, 0b111011010)
	case 0b000101001:
		r0 = NewFromUint(9, 0b111011010)
	case 0b000101010:
		r0 = NewFromUint(9, 0b111011001)
	case 0b000101011:
		r0 = NewFromUint(9, 0b111011000)
	case 0b000101100:
		r0 = NewFromUint(9, 0b111010111)
	case 0b000101101:
		r0 = NewFromUint(9, 0b111010110)
	case 0b000101110:
		r0 = NewFromUint(9, 0b111010101)
	case 0b000101111:
		r0 = NewFromUint(9, 0b111010100)
	case 0b000110000:
		r0 = NewFromUint(9, 0b111010100)
	case 0b000110001:
		r0 = NewFromUint(9, 0b111010011)
	case 0b000110010:
		r0 = NewFromUint(9, 0b111010010)
	case 0b000110011:
		r0 = NewFromUint(9, 0b111010001)
	case 0b000110100:
		r0 = NewFromUint(9, 0b111010000)
	case 0b000110101:
		r0 = NewFromUint(9, 0b111001111)
	case 0b000110110:
		r0 = NewFromUint(9, 0b111001111)
	case 0b000110111:
		r0 = NewFromUint(9, 0b111001110)
	case 0b000111000:
		r0 = NewFromUint(9, 0b111001101)
	case 0b000111001:
		r0 = NewFromUint(9, 0b111001100)
	case 0b000111010:
		r0 = NewFromUint(9, 0b111001011)
	case 0b000111011:
		r0 = NewFromUint(9, 0b111001011)
	case 0b000111100:
		r0 = NewFromUint(9, 0b111001010)
	case 0b000111101:
		r0 = NewFromUint(9, 0b111001001)
	case 0b000111110:
		r0 = NewFromUint(9, 0b111001000)
	case 0b000111111:
		r0 = NewFromUint(9, 0b111000111)
	case 0b001000000:
		r0 = NewFromUint(9, 0b111000111)
	case 0b001000001:
		r0 = NewFromUint(9, 0b111000110)
	case 0b001000010:
		r0 = NewFromUint(9, 0b111000101)
	case 0b001000011:
		r0 = NewFromUint(9, 0b111000100)
	case 0b001000100:
		r0 = NewFromUint(9, 0b111000011)
	case 0b001000101:
		r0 = NewFromUint(9, 0b111000011)
	case 0b001000110:
		r0 = NewFromUint(9, 0b111000010)
	case 0b001000111:
		r0 = NewFromUint(9, 0b111000001)
	case 0b001001000:
		r0 = NewFromUint(9, 0b111000000)
	case 0b001001001:
		r0 = NewFromUint(9, 0b111000000)
	case 0b001001010:
		r0 = NewFromUint(9, 0b110111111)
	case 0b001001011:
		r0 = NewFromUint(9, 0b110111110)
	case 0b001001100:
		r0 = NewFromUint(9, 0b110111101)
	case 0b001001101:
		r0 = NewFromUint(9, 0b110111101)
	case 0b001001110:
		r0 = NewFromUint(9, 0b110111100)
	case 0b001001111:
		r0 = NewFromUint(9, 0b110111011)
	case 0b001010000:
		r0 = NewFromUint(9, 0b110111010)
	case 0b001010001:
		r0 = NewFromUint(9, 0b110111010)
	case 0b001010010:
		r0 = NewFromUint(9, 0b110111001)
	case 0b001010011:
		r0 = NewFromUint(9, 0b110111000)
	case 0b001010100:
		r0 = NewFromUint(9, 0b110110111)
	case 0b001010101:
		r0 = NewFromUint(9, 0b110110111)
	case 0b001010110:
		r0 = NewFromUint(9, 0b110110110)
	case 0b001010111:
		r0 = NewFromUint(9, 0b110110101)
	case 0b001011000:
		r0 = NewFromUint(9, 0b110110100)
	case 0b001011001:
		r0 = NewFromUint(9, 0b110110100)
	case 0b001011010:
		r0 = NewFromUint(9, 0b110110011)
	case 0b001011011:
		r0 = NewFromUint(9, 0b110110010)
	case 0b001011100:
		r0 = NewFromUint(9, 0b110110010)
	case 0b001011101:
		r0 = NewFromUint(9, 0b110110001)
	case 0b001011110:
		r0 = NewFromUint(9, 0b110110000)
	case 0b001011111:
		r0 = NewFromUint(9, 0b110101111)
	case 0b001100000:
		r0 = NewFromUint(9, 0b110101111)
	case 0b001100001:
		r0 = NewFromUint(9, 0b110101110)
	case 0b001100010:
		r0 = NewFromUint(9, 0b110101101)
	case 0b001100011:
		r0 = NewFromUint(9, 0b110101101)
	case 0b001100100:
		r0 = NewFromUint(9, 0b110101100)
	case 0b001100101:
		r0 = NewFromUint(9, 0b110101011)
	case 0b001100110:
		r0 = NewFromUint(9, 0b110101010)
	case 0b001100111:
		r0 = NewFromUint(9, 0b110101010)
	case 0b001101000:
		r0 = NewFromUint(9, 0b110101001)
	case 0b001101001:
		r0 = NewFromUint(9, 0b110101000)
	case 0b001101010:
		r0 = NewFromUint(9, 0b110101000)
	case 0b001101011:
		r0 = NewFromUint(9, 0b110100111)
	case 0b001101100:
		r0 = NewFromUint(9, 0b110100110)
	case 0b001101101:
		r0 = NewFromUint(9, 0b110100110)
	case 0b001101110:
		r0 = NewFromUint(9, 0b110100101)
	case 0b001101111:
		r0 = NewFromUint(9, 0b110100100)
	case 0b001110000:
		r0 = NewFromUint(9, 0b110100100)
	case 0b001110001:
		r0 = NewFromUint(9, 0b110100011)
	case 0b001110010:
		r0 = NewFromUint(9, 0b110100010)
	case 0b001110011:
		r0 = NewFromUint(9, 0b110100010)
	case 0b001110100:
		r0 = NewFromUint(9, 0b110100001)
	case 0b001110101:
		r0 = NewFromUint(9, 0b110100000)
	case 0b001110110:
		r0 = NewFromUint(9, 0b110100000)
	case 0b001110111:
		r0 = NewFromUint(9, 0b110011111)
	case 0b001111000:
		r0 = NewFromUint(9, 0b110011110)
	case 0b001111001:
		r0 = NewFromUint(9, 0b110011110)
	case 0b001111010:
		r0 = NewFromUint(9, 0b110011101)
	case 0b001111011:
		r0 = NewFromUint(9, 0b110011100)
	case 0b001111100:
		r0 = NewFromUint(9, 0b110011100)
	case 0b001111101:
		r0 = NewFromUint(9, 0b110011011)
	case 0b001111110:
		r0 = NewFromUint(9, 0b110011010)
	case 0b001111111:
		r0 = NewFromUint(9, 0b110011010)
	case 0b010000000:
		r0 = NewFromUint(9, 0b110011001)
	case 0b010000001:
		r0 = NewFromUint(9, 0b110011000)
	case 0b010000010:
		r0 = NewFromUint(9, 0b110011000)
	case 0b010000011:
		r0 = NewFromUint(9, 0b110010111)
	case 0b010000100:
		r0 = NewFromUint(9, 0b110010111)
	case 0b010000101:
		r0 = NewFromUint(9, 0b110010110)
	case 0b010000110:
		r0 = NewFromUint(9, 0b110010101)
	case 0b010000111:
		r0 = NewFromUint(9, 0b110010101)
	case 0b010001000:
		r0 = NewFromUint(9, 0b110010100)
	case 0b010001001:
		r0 = NewFromUint(9, 0b110010011)
	case 0b010001010:
		r0 = NewFromUint(9, 0b110010011)
	case 0b010001011:
		r0 = NewFromUint(9, 0b110010010)
	case 0b010001100:
		r0 = NewFromUint(9, 0b110010010)
	case 0b010001101:
		r0 = NewFromUint(9, 0b110010001)
	case 0b010001110:
		r0 = NewFromUint(9, 0b110010000)
	case 0b010001111:
		r0 = NewFromUint(9, 0b110010000)
	case 0b010010000:
		r0 = NewFromUint(9, 0b110001111)
	case 0b010010001:
		r0 = NewFromUint(9, 0b110001111)
	case 0b010010010:
		r0 = NewFromUint(9, 0b110001110)
	case 0b010010011:
		r0 = NewFromUint(9, 0b110001101)
	case 0b010010100:
		r0 = NewFromUint(9, 0b110001101)
	case 0b010010101:
		r0 = NewFromUint(9, 0b110001100)
	case 0b010010110:
		r0 = NewFromUint(9, 0b110001011)
	case 0b010010111:
		r0 = NewFromUint(9, 0b110001011)
	case 0b010011000:
		r0 = NewFromUint(9, 0b110001010)
	case 0b010011001:
		r0 = NewFromUint(9, 0b110001010)
	case 0b010011010:
		r0 = NewFromUint(9, 0b110001001)
	case 0b010011011:
		r0 = NewFromUint(9, 0b110001001)
	case 0b010011100:
		r0 = NewFromUint(9, 0b110001000)
	case 0b010011101:
		r0 = NewFromUint(9, 0b110000111)
	case 0b010011110:
		r0 = NewFromUint(9, 0b110000111)
	case 0b010011111:
		r0 = NewFromUint(9, 0b110000110)
	case 0b010100000:
		r0 = NewFromUint(9, 0b110000110)
	case 0b010100001:
		r0 = NewFromUint(9, 0b110000101)
	case 0b010100010:
		r0 = NewFromUint(9, 0b110000100)
	case 0b010100011:
		r0 = NewFromUint(9, 0b110000100)
	case 0b010100100:
		r0 = NewFromUint(9, 0b110000011)
	case 0b010100101:
		r0 = NewFromUint(9, 0b110000011)
	case 0b010100110:
		r0 = NewFromUint(9, 0b110000010)
	case 0b010100111:
		r0 = NewFromUint(9, 0b110000010)
	case 0b010101000:
		r0 = NewFromUint(9, 0b110000001)
	case 0b010101001:
		r0 = NewFromUint(9, 0b110000000)
	case 0b010101010:
		r0 = NewFromUint(9, 0b110000000)
	case 0b010101011:
		r0 = NewFromUint(9, 0b101111111)
	case 0b010101100:
		r0 = NewFromUint(9, 0b101111111)
	case 0b010101101:
		r0 = NewFromUint(9, 0b101111110)
	case 0b010101110:
		r0 = NewFromUint(9, 0b101111110)
	case 0b010101111:
		r0 = NewFromUint(9, 0b101111101)
	case 0b010110000:
		r0 = NewFromUint(9, 0b101111101)
	case 0b010110001:
		r0 = NewFromUint(9, 0b101111100)
	case 0b010110010:
		r0 = NewFromUint(9, 0b101111011)
	case 0b010110011:
		r0 = NewFromUint(9, 0b101111011)
	case 0b010110100:
		r0 = NewFromUint(9, 0b101111010)
	case 0b010110101:
		r0 = NewFromUint(9, 0b101111010)
	case 0b010110110:
		r0 = NewFromUint(9, 0b101111001)
	case 0b010110111:
		r0 = NewFromUint(9, 0b101111001)
	case 0b010111000:
		r0 = NewFromUint(9, 0b101111000)
	case 0b010111001:
		r0 = NewFromUint(9, 0b101111000)
	case 0b010111010:
		r0 = NewFromUint(9, 0b101110111)
	case 0b010111011:
		r0 = NewFromUint(9, 0b101110111)
	case 0b010111100:
		r0 = NewFromUint(9, 0b101110110)
	case 0b010111101:
		r0 = NewFromUint(9, 0b101110101)
	case 0b010111110:
		r0 = NewFromUint(9, 0b101110101)
	case 0b010111111:
		r0 = NewFromUint(9, 0b101110100)
	case 0b011000000:
		r0 = NewFromUint(9, 0b101110100)
	case 0b011000001:
		r0 = NewFromUint(9, 0b101110011)
	case 0b011000010:
		r0 = NewFromUint(9, 0b101110011)
	case 0b011000011:
		r0 = NewFromUint(9, 0b101110010)
	case 0b011000100:
		r0 = NewFromUint(9, 0b101110010)
	case 0b011000101:
		r0 = NewFromUint(9, 0b101110001)
	case 0b011000110:
		r0 = NewFromUint(9, 0b101110001)
	case 0b011000111:
		r0 = NewFromUint(9, 0b101110000)
	case 0b011001000:
		r0 = NewFromUint(9, 0b101110000)
	case 0b011001001:
		r0 = NewFromUint(9, 0b101101111)
	case 0b011001010:
		r0 = NewFromUint(9, 0b101101111)
	case 0b011001011:
		r0 = NewFromUint(9, 0b101101110)
	case 0b011001100:
		r0 = NewFromUint(9, 0b101101110)
	case 0b011001101:
		r0 = NewFromUint(9, 0b101101101)
	case 0b011001110:
		r0 = NewFromUint(9, 0b101101101)
	case 0b011001111:
		r0 = NewFromUint(9, 0b101101100)
	case 0b011010000:
		r0 = NewFromUint(9, 0b101101100)
	case 0b011010001:
		r0 = NewFromUint(9, 0b101101011)
	case 0b011010010:
		r0 = NewFromUint(9, 0b101101011)
	case 0b011010011:
		r0 = NewFromUint(9, 0b101101010)
	case 0b011010100:
		r0 = NewFromUint(9, 0b101101010)
	case 0b011010101:
		r0 = NewFromUint(9, 0b101101001)
	case 0b011010110:
		r0 = NewFromUint(9, 0b101101001)
	case 0b011010111:
		r0 = NewFromUint(9, 0b101101000)
	case 0b011011000:
		r0 = NewFromUint(9, 0b101101000)
	case 0b011011001:
		r0 = NewFromUint(9, 0b101100111)
	case 0b011011010:
		r0 = NewFromUint(9, 0b101100111)
	case 0b011011011:
		r0 = NewFromUint(9, 0b101100110)
	case 0b011011100:
		r0 = NewFromUint(9, 0b101100110)
	case 0b011011101:
		r0 = NewFromUint(9, 0b101100101)
	case 0b011011110:
		r0 = NewFromUint(9, 0b101100101)
	case 0b011011111:
		r0 = NewFromUint(9, 0b101100100)
	case 0b011100000:
		r0 = NewFromUint(9, 0b101100100)
	case 0b011100001:
		r0 = NewFromUint(9, 0b101100011)
	case 0b011100010:
		r0 = NewFromUint(9, 0b101100011)
	case 0b011100011:
		r0 = NewFromUint(9, 0b101100010)
	case 0b011100100:
		r0 = NewFromUint(9, 0b101100010)
	case 0b011100101:
		r0 = NewFromUint(9, 0b101100001)
	case 0b011100110:
		r0 = NewFromUint(9, 0b101100001)
	case 0b011100111:
		r0 = NewFromUint(9, 0b101100000)
	case 0b011101000:
		r0 = NewFromUint(9, 0b101100000)
	case 0b011101001:
		r0 = NewFromUint(9, 0b101011111)
	case 0b011101010:
		r0 = NewFromUint(9, 0b101011111)
	case 0b011101011:
		r0 = NewFromUint(9, 0b101011110)
	case 0b011101100:
		r0 = NewFromUint(9, 0b101011110)
	case 0b011101101:
		r0 = NewFromUint(9, 0b101011101)
	case 0b011101110:
		r0 = NewFromUint(9, 0b101011101)
	case 0b011101111:
		r0 = NewFromUint(9, 0b101011101)
	case 0b011110000:
		r0 = NewFromUint(9, 0b101011100)
	case 0b011110001:
		r0 = NewFromUint(9, 0b101011100)
	case 0b011110010:
		r0 = NewFromUint(9, 0b101011011)
	case 0b011110011:
		r0 = NewFromUint(9, 0b101011011)
	case 0b011110100:
		r0 = NewFromUint(9, 0b101011010)
	case 0b011110101:
		r0 = NewFromUint(9, 0b101011010)
	case 0b011110110:
		r0 = NewFromUint(9, 0b101011001)
	case 0b011110111:
		r0 = NewFromUint(9, 0b101011001)
	case 0b011111000:
		r0 = NewFromUint(9, 0b101011000)
	case 0b011111001:
		r0 = NewFromUint(9, 0b101011000)
	case 0b011111010:
		r0 = NewFromUint(9, 0b101011000)
	case 0b011111011:
		r0 = NewFromUint(9, 0b101010111)
	case 0b011111100:
		r0 = NewFromUint(9, 0b101010111)
	case 0b011111101:
		r0 = NewFromUint(9, 0b101010110)
	case 0b011111110:
		r0 = NewFromUint(9, 0b101010110)
	case 0b011111111:
		r0 = NewFromUint(9, 0b101010101)
	case 0b100000000:
		r0 = NewFromUint(9, 0b101010101)
	case 0b100000001:
		r0 = NewFromUint(9, 0b101010100)
	case 0b100000010:
		r0 = NewFromUint(9, 0b101010100)
	case 0b100000011:
		r0 = NewFromUint(9, 0b101010100)
	case 0b100000100:
		r0 = NewFromUint(9, 0b101010011)
	case 0b100000101:
		r0 = NewFromUint(9, 0b101010011)
	case 0b100000110:
		r0 = NewFromUint(9, 0b101010010)
	case 0b100000111:
		r0 = NewFromUint(9, 0b101010010)
	case 0b100001000:
		r0 = NewFromUint(9, 0b101010001)
	case 0b100001001:
		r0 = NewFromUint(9, 0b101010001)
	case 0b100001010:
		r0 = NewFromUint(9, 0b101010000)
	case 0b100001011:
		r0 = NewFromUint(9, 0b101010000)
	case 0b100001100:
		r0 = NewFromUint(9, 0b101010000)
	case 0b100001101:
		r0 = NewFromUint(9, 0b101001111)
	case 0b100001110:
		r0 = NewFromUint(9, 0b101001111)
	case 0b100001111:
		r0 = NewFromUint(9, 0b101001110)
	case 0b100010000:
		r0 = NewFromUint(9, 0b101001110)
	case 0b100010001:
		r0 = NewFromUint(9, 0b101001101)
	case 0b100010010:
		r0 = NewFromUint(9, 0b101001101)
	case 0b100010011:
		r0 = NewFromUint(9, 0b101001101)
	case 0b100010100:
		r0 = NewFromUint(9, 0b101001100)
	case 0b100010101:
		r0 = NewFromUint(9, 0b101001100)
	case 0b100010110:
		r0 = NewFromUint(9, 0b101001011)
	case 0b100010111:
		r0 = NewFromUint(9, 0b101001011)
	case 0b100011000:
		r0 = NewFromUint(9, 0b101001010)
	case 0b100011001:
		r0 = NewFromUint(9, 0b101001010)
	case 0b100011010:
		r0 = NewFromUint(9, 0b101001010)
	case 0b100011011:
		r0 = NewFromUint(9, 0b101001001)
	case 0b100011100:
		r0 = NewFromUint(9, 0b101001001)
	case 0b100011101:
		r0 = NewFromUint(9, 0b101001000)
	case 0b100011110:
		r0 = NewFromUint(9, 0b101001000)
	case 0b100011111:
		r0 = NewFromUint(9, 0b101001000)
	case 0b100100000:
		r0 = NewFromUint(9, 0b101000111)
	case 0b100100001:
		r0 = NewFromUint(9, 0b101000111)
	case 0b100100010:
		r0 = NewFromUint(9, 0b101000110)
	case 0b100100011:
		r0 = NewFromUint(9, 0b101000110)
	case 0b100100100:
		r0 = NewFromUint(9, 0b101000110)
	case 0b100100101:
		r0 = NewFromUint(9, 0b101000101)
	case 0b100100110:
		r0 = NewFromUint(9, 0b101000101)
	case 0b100100111:
		r0 = NewFromUint(9, 0b101000100)
	case 0b100101000:
		r0 = NewFromUint(9, 0b101000100)
	case 0b100101001:
		r0 = NewFromUint(9, 0b101000100)
	case 0b100101010:
		r0 = NewFromUint(9, 0b101000011)
	case 0b100101011:
		r0 = NewFromUint(9, 0b101000011)
	case 0b100101100:
		r0 = NewFromUint(9, 0b101000010)
	case 0b100101101:
		r0 = NewFromUint(9, 0b101000010)
	case 0b100101110:
		r0 = NewFromUint(9, 0b101000010)
	case 0b100101111:
		r0 = NewFromUint(9, 0b101000001)
	case 0b100110000:
		r0 = NewFromUint(9, 0b101000001)
	case 0b100110001:
		r0 = NewFromUint(9, 0b101000000)
	case 0b100110010:
		r0 = NewFromUint(9, 0b101000000)
	case 0b100110011:
		r0 = NewFromUint(9, 0b101000000)
	case 0b100110100:
		r0 = NewFromUint(9, 0b100111111)
	case 0b100110101:
		r0 = NewFromUint(9, 0b100111111)
	case 0b100110110:
		r0 = NewFromUint(9, 0b100111110)
	case 0b100110111:
		r0 = NewFromUint(9, 0b100111110)
	case 0b100111000:
		r0 = NewFromUint(9, 0b100111110)
	case 0b100111001:
		r0 = NewFromUint(9, 0b100111101)
	case 0b100111010:
		r0 = NewFromUint(9, 0b100111101)
	case 0b100111011:
		r0 = NewFromUint(9, 0b100111100)
	case 0b100111100:
		r0 = NewFromUint(9, 0b100111100)
	case 0b100111101:
		r0 = NewFromUint(9, 0b100111100)
	case 0b100111110:
		r0 = NewFromUint(9, 0b100111011)
	case 0b100111111:
		r0 = NewFromUint(9, 0b100111011)
	case 0b101000000:
		r0 = NewFromUint(9, 0b100111011)
	case 0b101000001:
		r0 = NewFromUint(9, 0b100111010)
	case 0b101000010:
		r0 = NewFromUint(9, 0b100111010)
	case 0b101000011:
		r0 = NewFromUint(9, 0b100111001)
	case 0b101000100:
		r0 = NewFromUint(9, 0b100111001)
	case 0b101000101:
		r0 = NewFromUint(9, 0b100111001)
	case 0b101000110:
		r0 = NewFromUint(9, 0b100111000)
	case 0b101000111:
		r0 = NewFromUint(9, 0b100111000)
	case 0b101001000:
		r0 = NewFromUint(9, 0b100111000)
	case 0b101001001:
		r0 = NewFromUint(9, 0b100110111)
	case 0b101001010:
		r0 = NewFromUint(9, 0b100110111)
	case 0b101001011:
		r0 = NewFromUint(9, 0b100110110)
	case 0b101001100:
		r0 = NewFromUint(9, 0b100110110)
	case 0b101001101:
		r0 = NewFromUint(9, 0b100110110)
	case 0b101001110:
		r0 = NewFromUint(9, 0b100110101)
	case 0b101001111:
		r0 = NewFromUint(9, 0b100110101)
	case 0b101010000:
		r0 = NewFromUint(9, 0b100110101)
	case 0b101010001:
		r0 = NewFromUint(9, 0b100110100)
	case 0b101010010:
		r0 = NewFromUint(9, 0b100110100)
	case 0b101010011:
		r0 = NewFromUint(9, 0b100110100)
	case 0b101010100:
		r0 = NewFromUint(9, 0b100110011)
	case 0b101010101:
		r0 = NewFromUint(9, 0b100110011)
	case 0b101010110:
		r0 = NewFromUint(9, 0b100110010)
	case 0b101010111:
		r0 = NewFromUint(9, 0b100110010)
	case 0b101011000:
		r0 = NewFromUint(9, 0b100110010)
	case 0b101011001:
		r0 = NewFromUint(9, 0b100110001)
	case 0b101011010:
		r0 = NewFromUint(9, 0b100110001)
	case 0b101011011:
		r0 = NewFromUint(9, 0b100110001)
	case 0b101011100:
		r0 = NewFromUint(9, 0b100110000)
	case 0b101011101:
		r0 = NewFromUint(9, 0b100110000)
	case 0b101011110:
		r0 = NewFromUint(9, 0b100110000)
	case 0b101011111:
		r0 = NewFromUint(9, 0b100101111)
	case 0b101100000:
		r0 = NewFromUint(9, 0b100101111)
	case 0b101100001:
		r0 = NewFromUint(9, 0b100101111)
	case 0b101100010:
		r0 = NewFromUint(9, 0b100101110)
	case 0b101100011:
		r0 = NewFromUint(9, 0b100101110)
	case 0b101100100:
		r0 = NewFromUint(9, 0b100101110)
	case 0b101100101:
		r0 = NewFromUint(9, 0b100101101)
	case 0b101100110:
		r0 = NewFromUint(9, 0b100101101)
	case 0b101100111:
		r0 = NewFromUint(9, 0b100101100)
	case 0b101101000:
		r0 = NewFromUint(9, 0b100101100)
	case 0b101101001:
		r0 = NewFromUint(9, 0b100101100)
	case 0b101101010:
		r0 = NewFromUint(9, 0b100101011)
	case 0b101101011:
		r0 = NewFromUint(9, 0b100101011)
	case 0b101101100:
		r0 = NewFromUint(9, 0b100101011)
	case 0b101101101:
		r0 = NewFromUint(9, 0b100101010)
	case 0b101101110:
		r0 = NewFromUint(9, 0b100101010)
	case 0b101101111:
		r0 = NewFromUint(9, 0b100101010)
	case 0b101110000:
		r0 = NewFromUint(9, 0b100101001)
	case 0b101110001:
		r0 = NewFromUint(9, 0b100101001)
	case 0b101110010:
		r0 = NewFromUint(9, 0b100101001)
	case 0b101110011:
		r0 = NewFromUint(9, 0b100101000)
	case 0b101110100:
		r0 = NewFromUint(9, 0b100101000)
	case 0b101110101:
		r0 = NewFromUint(9, 0b100101000)
	case 0b101110110:
		r0 = NewFromUint(9, 0b100100111)
	case 0b101110111:
		r0 = NewFromUint(9, 0b100100111)
	case 0b101111000:
		r0 = NewFromUint(9, 0b100100111)
	case 0b101111001:
		r0 = NewFromUint(9, 0b100100110)
	case 0b101111010:
		r0 = NewFromUint(9, 0b100100110)
	case 0b101111011:
		r0 = NewFromUint(9, 0b100100110)
	case 0b101111100:
		r0 = NewFromUint(9, 0b100100101)
	case 0b101111101:
		r0 = NewFromUint(9, 0b100100101)
	case 0b101111110:
		r0 = NewFromUint(9, 0b100100101)
	case 0b101111111:
		r0 = NewFromUint(9, 0b100100100)
	case 0b110000000:
		r0 = NewFromUint(9, 0b100100100)
	case 0b110000001:
		r0 = NewFromUint(9, 0b100100100)
	case 0b110000010:
		r0 = NewFromUint(9, 0b100100011)
	case 0b110000011:
		r0 = NewFromUint(9, 0b100100011)
	case 0b110000100:
		r0 = NewFromUint(9, 0b100100011)
	case 0b110000101:
		r0 = NewFromUint(9, 0b100100010)
	case 0b110000110:
		r0 = NewFromUint(9, 0b100100010)
	case 0b110000111:
		r0 = NewFromUint(9, 0b100100010)
	case 0b110001000:
		r0 = NewFromUint(9, 0b100100001)
	case 0b110001001:
		r0 = NewFromUint(9, 0b100100001)
	case 0b110001010:
		r0 = NewFromUint(9, 0b100100001)
	case 0b110001011:
		r0 = NewFromUint(9, 0b100100001)
	case 0b110001100:
		r0 = NewFromUint(9, 0b100100000)
	case 0b110001101:
		r0 = NewFromUint(9, 0b100100000)
	case 0b110001110:
		r0 = NewFromUint(9, 0b100100000)
	case 0b110001111:
		r0 = NewFromUint(9, 0b100011111)
	case 0b110010000:
		r0 = NewFromUint(9, 0b100011111)
	case 0b110010001:
		r0 = NewFromUint(9, 0b100011111)
	case 0b110010010:
		r0 = NewFromUint(9, 0b100011110)
	case 0b110010011:
		r0 = NewFromUint(9, 0b100011110)
	case 0b110010100:
		r0 = NewFromUint(9, 0b100011110)
	case 0b110010101:
		r0 = NewFromUint(9, 0b100011101)
	case 0b110010110:
		r0 = NewFromUint(9, 0b100011101)
	case 0b110010111:
		r0 = NewFromUint(9, 0b100011101)
	case 0b110011000:
		r0 = NewFromUint(9, 0b100011100)
	case 0b110011001:
		r0 = NewFromUint(9, 0b100011100)
	case 0b110011010:
		r0 = NewFromUint(9, 0b100011100)
	case 0b110011011:
		r0 = NewFromUint(9, 0b100011100)
	case 0b110011100:
		r0 = NewFromUint(9, 0b100011011)
	case 0b110011101:
		r0 = NewFromUint(9, 0b100011011)
	case 0b110011110:
		r0 = NewFromUint(9, 0b100011011)
	case 0b110011111:
		r0 = NewFromUint(9, 0b100011010)
	case 0b110100000:
		r0 = NewFromUint(9, 0b100011010)
	case 0b110100001:
		r0 = NewFromUint(9, 0b100011010)
	case 0b110100010:
		r0 = NewFromUint(9, 0b100011001)
	case 0b110100011:
		r0 = NewFromUint(9, 0b100011001)
	case 0b110100100:
		r0 = NewFromUint(9, 0b100011001)
	case 0b110100101:
		r0 = NewFromUint(9, 0b100011000)
	case 0b110100110:
		r0 = NewFromUint(9, 0b100011000)
	case 0b110100111:
		r0 = NewFromUint(9, 0b100011000)
	case 0b110101000:
		r0 = NewFromUint(9, 0b100011000)
	case 0b110101001:
		r0 = NewFromUint(9, 0b100010111)
	case 0b110101010:
		r0 = NewFromUint(9, 0b100010111)
	case 0b110101011:
		r0 = NewFromUint(9, 0b100010111)
	case 0b110101100:
		r0 = NewFromUint(9, 0b100010110)
	case 0b110101101:
		r0 = NewFromUint(9, 0b100010110)
	case 0b110101110:
		r0 = NewFromUint(9, 0b100010110)
	case 0b110101111:
		r0 = NewFromUint(9, 0b100010101)
	case 0b110110000:
		r0 = NewFromUint(9, 0b100010101)
	case 0b110110001:
		r0 = NewFromUint(9, 0b100010101)
	case 0b110110010:
		r0 = NewFromUint(9, 0b100010101)
	case 0b110110011:
		r0 = NewFromUint(9, 0b100010100)
	case 0b110110100:
		r0 = NewFromUint(9, 0b100010100)
	case 0b110110101:
		r0 = NewFromUint(9, 0b100010100)
	case 0b110110110:
		r0 = NewFromUint(9, 0b100010011)
	case 0b110110111:
		r0 = NewFromUint(9, 0b100010011)
	case 0b110111000:
		r0 = NewFromUint(9, 0b100010011)
	case 0b110111001:
		r0 = NewFromUint(9, 0b100010011)
	case 0b110111010:
		r0 = NewFromUint(9, 0b100010010)
	case 0b110111011:
		r0 = NewFromUint(9, 0b100010010)
	case 0b110111100:
		r0 = NewFromUint(9, 0b100010010)
	case 0b110111101:
		r0 = NewFromUint(9, 0b100010001)
	case 0b110111110:
		r0 = NewFromUint(9, 0b100010001)
	case 0b110111111:
		r0 = NewFromUint(9, 0b100010001)
	case 0b111000000:
		r0 = NewFromUint(9, 0b100010001)
	case 0b111000001:
		r0 = NewFromUint(9, 0b100010000)
	case 0b111000010:
		r0 = NewFromUint(9, 0b100010000)
	case 0b111000011:
		r0 = NewFromUint(9, 0b100010000)
	case 0b111000100:
		r0 = NewFromUint(9, 0b100001111)
	case 0b111000101:
		r0 = NewFromUint(9, 0b100001111)
	case 0b111000110:
		r0 = NewFromUint(9, 0b100001111)
	case 0b111000111:
		r0 = NewFromUint(9, 0b100001111)
	case 0b111001000:
		r0 = NewFromUint(9, 0b100001110)
	case 0b111001001:
		r0 = NewFromUint(9, 0b100001110)
	case 0b111001010:
		r0 = NewFromUint(9, 0b100001110)
	case 0b111001011:
		r0 = NewFromUint(9, 0b100001101)
	case 0b111001100:
		r0 = NewFromUint(9, 0b100001101)
	case 0b111001101:
		r0 = NewFromUint(9, 0b100001101)
	case 0b111001110:
		r0 = NewFromUint(9, 0b100001101)
	case 0b111001111:
		r0 = NewFromUint(9, 0b100001100)
	case 0b111010000:
		r0 = NewFromUint(9, 0b100001100)
	case 0b111010001:
		r0 = NewFromUint(9, 0b100001100)
	case 0b111010010:
		r0 = NewFromUint(9, 0b100001100)
	case 0b111010011:
		r0 = NewFromUint(9, 0b100001011)
	case 0b111010100:
		r0 = NewFromUint(9, 0b100001011)
	case 0b111010101:
		r0 = NewFromUint(9, 0b100001011)
	case 0b111010110:
		r0 = NewFromUint(9, 0b100001010)
	case 0b111010111:
		r0 = NewFromUint(9, 0b100001010)
	case 0b111011000:
		r0 = NewFromUint(9, 0b100001010)
	case 0b111011001:
		r0 = NewFromUint(9, 0b100001010)
	case 0b111011010:
		r0 = NewFromUint(9, 0b100001001)
	case 0b111011011:
		r0 = NewFromUint(9, 0b100001001)
	case 0b111011100:
		r0 = NewFromUint(9, 0b100001001)
	case 0b111011101:
		r0 = NewFromUint(9, 0b100001001)
	case 0b111011110:
		r0 = NewFromUint(9, 0b100001000)
	case 0b111011111:
		r0 = NewFromUint(9, 0b100001000)
	case 0b111100000:
		r0 = NewFromUint(9, 0b100001000)
	case 0b111100001:
		r0 = NewFromUint(9, 0b100000111)
	case 0b111100010:
		r0 = NewFromUint(9, 0b100000111)
	case 0b111100011:
		r0 = NewFromUint(9, 0b100000111)
	case 0b111100100:
		r0 = NewFromUint(9, 0b100000111)
	case 0b111100101:
		r0 = NewFromUint(9, 0b100000110)
	case 0b111100110:
		r0 = NewFromUint(9, 0b100000110)
	case 0b111100111:
		r0 = NewFromUint(9, 0b100000110)
	case 0b111101000:
		r0 = NewFromUint(9, 0b100000110)
	case 0b111101001:
		r0 = NewFromUint(9, 0b100000101)
	case 0b111101010:
		r0 = NewFromUint(9, 0b100000101)
	case 0b111101011:
		r0 = NewFromUint(9, 0b100000101)
	case 0b111101100:
		r0 = NewFromUint(9, 0b100000101)
	case 0b111101101:
		r0 = NewFromUint(9, 0b100000100)
	case 0b111101110:
		r0 = NewFromUint(9, 0b100000100)
	case 0b111101111:
		r0 = NewFromUint(9, 0b100000100)
	case 0b111110000:
		r0 = NewFromUint(9, 0b100000100)
	case 0b111110001:
		r0 = NewFromUint(9, 0b100000011)
	case 0b111110010:
		r0 = NewFromUint(9, 0b100000011)
	case 0b111110011:
		r0 = NewFromUint(9, 0b100000011)
	case 0b111110100:
		r0 = NewFromUint(9, 0b100000011)
	case 0b111110101:
		r0 = NewFromUint(9, 0b100000010)
	case 0b111110110:
		r0 = NewFromUint(9, 0b100000010)
	case 0b111110111:
		r0 = NewFromUint(9, 0b100000010)
	case 0b111111000:
		r0 = NewFromUint(9, 0b100000010)
	case 0b111111001:
		r0 = NewFromUint(9, 0b100000001)
	case 0b111111010:
		r0 = NewFromUint(9, 0b100000001)
	case 0b111111011:
		r0 = NewFromUint(9, 0b100000001)
	case 0b111111100:
		r0 = NewFromUint(9, 0b100000001)
	case 0b111111101:
		r0 = NewFromUint(9, 0b100000000)
	case 0b111111110:
		r0 = NewFromUint(9, 0b100000000)
	case 0b111111111:
		r0 = NewFromUint(9, 0b100000000)
	default:
		r0 = NewFromUint(9, 0b100000000)

	}

	var a1Pre, b1Pre Bits
	if m2.Slice(22, 14).Uint() == 0 {
		a1Pre = NewFromUint(48,
			NewFromUint(1, 1).Append(m1).Uint()*NewFromUint(1, 1).Append(r0).Append(NewFromUint(14, 0)).Uint())
		b1Pre = NewFromUint(48,
			NewFromUint(1, 1).Append(m2).Uint()*NewFromUint(1, 1).Append(r0).Append(NewFromUint(14, 0)).Uint())
	} else {
		a1Pre = NewFromUint(48,
			NewFromUint(1, 1).Append(m1).Uint()*NewFromUint(1, 0).Append(r0).Append(NewFromUint(14, 0)).Uint())
		b1Pre = NewFromUint(48,
			NewFromUint(1, 1).Append(m2).Uint()*NewFromUint(1, 0).Append(r0).Append(NewFromUint(14, 0)).Uint())
	}

	a1 := a1Pre.Slice(46, 21)
	b1 := b1Pre.Slice(46, 21)
	r1 := b1.Reverse()
	a2Pre := NewFromUint(52, a1.Uint()*r1.Uint())
	b2Pre := NewFromUint(52, b1.Uint()*r1.Uint())
	a2 := a2Pre.Slice(50, 25)
	b2 := b2Pre.Slice(50, 25)
	r2 := b2.Reverse()
	mPre := NewFromUint(52, a2.Uint()*r2.Uint())

	var m, e Bits
	if mPre.Slice(50, 50).Uint() > 0 {
		m = mPre.Slice(49, 27)
		e = eMoge
	} else {
		m = mPre.Slice(48, 26)
		e = NewFromUint(8, eMoge.Uint()-1)
	}

	if x1.Uint() == 0 {
		return NewFromUint(32, 0)
	} else {
		return s.Append(e).Append(m)
	}
}

func Mul(x1, x2 Bits) Bits {
	var s Bits
	if (x1.Slice(31, 31).Uint() > 0) != (x2.Slice(31, 31).Uint() > 0) {
		s = NewFromUint(1, 1)
	} else {
		s = NewFromUint(1, 0)
	}

	e1 := x1.Slice(30, 23)
	e2 := x2.Slice(30, 23)
	m1 := x1.Slice(22, 0)
	m2 := x2.Slice(22, 0)

	mReal := NewFromUint(48,
		NewFromUint(1, 1).Append(m1).Uint()*NewFromUint(1, 1).Append(m2).Uint())
	eReal := NewFromInt(10,
		int(e1.Uint())+int(e2.Uint())-127)

	var mMid Bits
	if eReal.Slice(8, 0).Uint() == 0 {
		mMid = NewFromUint(48, mReal.Uint()>>1)
	} else if eReal.Slice(9, 9).Uint() == 1 {
		mMid = NewFromUint(48, mReal.Uint()>>(-eReal.Uint()+1))
	} else {
		mMid = mReal
	}

	var eMoge Bits
	if mReal.Slice(47, 47).Uint() == 1 || (eReal.Slice(8, 0).Uint() == 0 && mReal.Slice(46, 24).Reverse().Uint() == 0) {
		eMoge = NewFromUint(8, eReal.Slice(7, 0).Uint()+1)
	} else {
		eMoge = eReal.Slice(7, 0)
	}

	var m Bits
	if mMid.Slice(47, 47).Uint() == 1 {
		if mMid.Slice(23, 23).Uint() == 1 && (mMid.Slice(24, 24).Uint() == 1 || mMid.Slice(22, 0).Uint() > 0) {
			m = NewFromUint(23, mMid.Slice(46, 24).Uint()+1)
		} else {
			m = mMid.Slice(46, 24)
		}
	} else {
		if mMid.Slice(22, 22).Uint() == 1 && (mMid.Slice(23, 23).Uint() == 1 || mMid.Slice(21, 0).Uint() > 0) {
			m = NewFromUint(23, mMid.Slice(45, 23).Uint()+1)
		} else {
			m = mMid.Slice(45, 23)
		}
	}

	var e Bits
	if eMoge.Slice(7, 0).Uint() == 0 {
		e = NewFromUint(8, 0)
	} else {
		e = eMoge
	}

	if x1.Uint() == 0 || x2.Uint() == 0 {
		return NewFromUint(32, 0)
	} else {
		return s.Append(e).Append(m)
	}
}

func Sub(x1, x2 Bits) Bits {
	s1 := x1.Slice(31, 31)
	e1 := x1.Slice(30, 23)
	m1 := x1.Slice(22, 0)
	s2 := x2.Slice(31, 31).Reverse()
	e2 := x2.Slice(30, 23)
	m2 := x2.Slice(22, 0)

	var m1a Bits
	if e1.Uint() == 0 {
		m1a = NewFromUint(2, 0).Append(m1)
	} else {
		m1a = NewFromUint(2, 1).Append(m1)
	}

	var m2a Bits
	if e2.Uint() == 0 {
		m2a = NewFromUint(2, 0).Append(m2)
	} else {
		m2a = NewFromUint(2, 1).Append(m2)
	}
	var e1a Bits
	if e1.Uint() > 0 {
		e1a = e1
	} else {
		e1a = NewFromUint(8, 1)
	}

	var e2a Bits
	if e2.Uint() > 0 {
		e2a = e2
	} else {
		e2a = NewFromUint(8, 1)
	}

	e2ai := e2a.Reverse()
	te := NewFromUint(9, e1a.Uint()+e2ai.Uint())

	var ce, hog, tde Bits
	if te.Slice(8, 8).Uint() == 1 {
		ce = NewFromUint(1, 0)
		hog = NewFromUint(10, te.Uint()+1)
		tde = hog.Slice(7, 0)
	} else {
		ce = NewFromUint(1, 1)
		hog = NewFromUint(1, 0).Append(te.Reverse())
		tde = te.Slice(7, 0).Reverse()
	}

	var de Bits
	if tde.Uint() > 31 {
		de = NewFromUint(5, 31)
	} else {
		de = tde.Slice(4, 0)
	}

	var sel Bits
	if de.Uint() == 0 {
		if m1a.Uint() > m2a.Uint() {
			sel = NewFromUint(1, 0)
		} else {
			sel = NewFromUint(1, 1)
		}
	} else {
		sel = ce
	}

	var ms, mi, es, ss Bits
	if sel.Uint() == 0 {
		ms = m1a
		mi = m2a
		es = e1a
		ss = s1
	} else {
		ms = m2a
		mi = m1a
		es = e2a
		ss = s2
	}

	mie := mi.Append(NewFromUint(31, 0))
	mia := NewFromUint(56, mie.Uint()>>de.Uint())

	var tstck Bits
	if mia.Slice(28, 0).Uint() > 0 {
		tstck = NewFromUint(1, 1)
	} else {
		tstck = NewFromUint(1, 0)
	}

	var mye Bits
	if s1.Uint() == s2.Uint() {
		mye = NewFromUint(27,
			ms.Append(NewFromUint(2, 0)).Uint()+mia.Slice(55, 29).Uint())
	} else {
		mye = NewFromUint(27,
			ms.Append(NewFromUint(2, 0)).Uint()-mia.Slice(55, 29).Uint())
	}

	esi := NewFromUint(8, es.Uint()+1)

	var eyd, myd, stck Bits
	if mye.Slice(26, 26).Uint() == 1 {
		if esi.Uint() == 255 {
			eyd = NewFromUint(8, 255)
			myd = NewFromUint(2, 1).Append(NewFromUint(25, 0))
			stck = NewFromUint(1, 0)
		} else {
			eyd = esi
			myd = NewFromUint(27, mye.Uint()>>1)
			if tstck.Uint() > 0 || mye.Slice(0, 0).Uint() > 0 {
				stck = NewFromUint(1, 1)
			} else {
				stck = NewFromUint(1, 0)
			}
		}
	} else {
		eyd = es
		myd = mye
		stck = tstck
	}

	var se Bits

	for i := 25; i >= 0; i-- {
		if myd.Slice(i, i).Uint() == 1 {
			se = NewFromUint(5, 25-uint(i))
			break
		}
		if i == 0 {
			se = NewFromUint(5, 26)
		}
	}

	eyf := NewFromUint(9, eyd.Uint()-se.Uint())

	var eyr, myf Bits

	if eyf.Int() > 0 {
		eyr = eyf.Slice(7, 0)
		myf = NewFromUint(27, myd.Uint()<<se.Uint())
	} else {
		eyr = NewFromUint(8, 0)
		myf = NewFromUint(27, myd.Uint()<<(eyd.Slice(4, 0).Uint()-1))
	}

	var myr Bits
	if (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 0 && stck.Slice(0, 0).Uint() == 0 && myf.Slice(2, 2).Uint() == 1) || (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 0 && s1.Uint() == s2.Uint() && stck.Slice(0, 0).Uint() == 1) || (myf.Slice(1, 1).Uint() == 1 && myf.Slice(0, 0).Uint() == 1) {
		myr = NewFromUint(25, myf.Slice(26, 2).Uint()+1)
	} else {
		myr = myf.Slice(26, 2)
	}

	eyri := NewFromUint(8, eyr.Uint()+1)

	var my, ey Bits
	if myr.Slice(24, 24).Uint() == 1 {
		my = NewFromUint(23, 0)
		ey = eyri
	} else if myr.Slice(23, 0).Uint() == 0 {
		my = NewFromUint(23, 0)
		ey = NewFromUint(8, 0)
	} else {
		my = myr.Slice(22, 0)
		ey = eyr
	}

	var sy Bits
	if ey.Uint() == 0 && my.Uint() == 0 {
		if s1.Uint() > 0 && s2.Uint() > 0 {
			sy = NewFromUint(1, 1)
		} else {
			sy = NewFromUint(1, 0)
		}
	} else {
		sy = ss
	}

	var nzm1 Bits
	if m1.Slice(22, 0).Uint() > 0 {
		nzm1 = NewFromUint(1, 1)
	} else {
		nzm1 = NewFromUint(1, 0)
	}

	var nzm2 Bits
	if m2.Slice(22, 0).Uint() > 0 {
		nzm2 = NewFromUint(1, 1)
	} else {
		nzm2 = NewFromUint(1, 0)
	}

	var y Bits
	if e1.Uint() == 255 && e2.Uint() != 255 {
		y = s1.Append(NewFromUint(8, 255)).Append(nzm1).Append(m1.Slice(21, 0))
	} else if e2.Uint() == 255 && e1.Uint() != 255 {
		y = s2.Append(NewFromUint(8, 255)).Append(nzm2).Append(m2.Slice(21, 0))
	} else if e2.Uint() == 255 && nzm2.Uint() > 0 {
		y = s2.Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(m2.Slice(21, 0))
	} else if e1.Uint() == 255 && nzm1.Uint() > 0 {
		y = s1.Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(m1.Slice(21, 0))
	} else if e1.Uint() == 255 && e2.Uint() == 255 && s1.Uint() == s2.Uint() {
		y = s1.Append(NewFromUint(8, 255)).Append(NewFromUint(23, 0))
	} else if e1.Uint() == 255 && e2.Uint() == 255 {
		y = NewFromUint(1, 1).Append(NewFromUint(8, 255)).Append(NewFromUint(1, 1)).Append(NewFromUint(22, 0))
	} else {
		y = sy.Append(ey).Append(my)
	}

	return y
}
