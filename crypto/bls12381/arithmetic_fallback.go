


































package bls12381

import (
	"math/bits"
)

func add(z, x, y *fe) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], carry = bits.Add64(x[3], y[3], carry)
	z[4], carry = bits.Add64(x[4], y[4], carry)
	z[5], _ = bits.Add64(x[5], y[5], carry)

	
	
	if !(z[5] < 1873798617647539866 || (z[5] == 1873798617647539866 && (z[4] < 5412103778470702295 || (z[4] == 5412103778470702295 && (z[3] < 7239337960414712511 || (z[3] == 7239337960414712511 && (z[2] < 7435674573564081700 || (z[2] == 7435674573564081700 && (z[1] < 2210141511517208575 || (z[1] == 2210141511517208575 && (z[0] < 13402431016077863595))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 13402431016077863595, 0)
		z[1], b = bits.Sub64(z[1], 2210141511517208575, b)
		z[2], b = bits.Sub64(z[2], 7435674573564081700, b)
		z[3], b = bits.Sub64(z[3], 7239337960414712511, b)
		z[4], b = bits.Sub64(z[4], 5412103778470702295, b)
		z[5], _ = bits.Sub64(z[5], 1873798617647539866, b)
	}
}

func addAssign(x, y *fe) {
	var carry uint64

	x[0], carry = bits.Add64(x[0], y[0], 0)
	x[1], carry = bits.Add64(x[1], y[1], carry)
	x[2], carry = bits.Add64(x[2], y[2], carry)
	x[3], carry = bits.Add64(x[3], y[3], carry)
	x[4], carry = bits.Add64(x[4], y[4], carry)
	x[5], _ = bits.Add64(x[5], y[5], carry)

	
	
	if !(x[5] < 1873798617647539866 || (x[5] == 1873798617647539866 && (x[4] < 5412103778470702295 || (x[4] == 5412103778470702295 && (x[3] < 7239337960414712511 || (x[3] == 7239337960414712511 && (x[2] < 7435674573564081700 || (x[2] == 7435674573564081700 && (x[1] < 2210141511517208575 || (x[1] == 2210141511517208575 && (x[0] < 13402431016077863595))))))))))) {
		var b uint64
		x[0], b = bits.Sub64(x[0], 13402431016077863595, 0)
		x[1], b = bits.Sub64(x[1], 2210141511517208575, b)
		x[2], b = bits.Sub64(x[2], 7435674573564081700, b)
		x[3], b = bits.Sub64(x[3], 7239337960414712511, b)
		x[4], b = bits.Sub64(x[4], 5412103778470702295, b)
		x[5], _ = bits.Sub64(x[5], 1873798617647539866, b)
	}
}

func ladd(z, x, y *fe) {
	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], carry = bits.Add64(x[3], y[3], carry)
	z[4], carry = bits.Add64(x[4], y[4], carry)
	z[5], _ = bits.Add64(x[5], y[5], carry)
}

func laddAssign(x, y *fe) {
	var carry uint64
	x[0], carry = bits.Add64(x[0], y[0], 0)
	x[1], carry = bits.Add64(x[1], y[1], carry)
	x[2], carry = bits.Add64(x[2], y[2], carry)
	x[3], carry = bits.Add64(x[3], y[3], carry)
	x[4], carry = bits.Add64(x[4], y[4], carry)
	x[5], _ = bits.Add64(x[5], y[5], carry)
}

func double(z, x *fe) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], carry = bits.Add64(x[2], x[2], carry)
	z[3], carry = bits.Add64(x[3], x[3], carry)
	z[4], carry = bits.Add64(x[4], x[4], carry)
	z[5], _ = bits.Add64(x[5], x[5], carry)

	
	
	if !(z[5] < 1873798617647539866 || (z[5] == 1873798617647539866 && (z[4] < 5412103778470702295 || (z[4] == 5412103778470702295 && (z[3] < 7239337960414712511 || (z[3] == 7239337960414712511 && (z[2] < 7435674573564081700 || (z[2] == 7435674573564081700 && (z[1] < 2210141511517208575 || (z[1] == 2210141511517208575 && (z[0] < 13402431016077863595))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 13402431016077863595, 0)
		z[1], b = bits.Sub64(z[1], 2210141511517208575, b)
		z[2], b = bits.Sub64(z[2], 7435674573564081700, b)
		z[3], b = bits.Sub64(z[3], 7239337960414712511, b)
		z[4], b = bits.Sub64(z[4], 5412103778470702295, b)
		z[5], _ = bits.Sub64(z[5], 1873798617647539866, b)
	}
}

func doubleAssign(z *fe) {
	var carry uint64

	z[0], carry = bits.Add64(z[0], z[0], 0)
	z[1], carry = bits.Add64(z[1], z[1], carry)
	z[2], carry = bits.Add64(z[2], z[2], carry)
	z[3], carry = bits.Add64(z[3], z[3], carry)
	z[4], carry = bits.Add64(z[4], z[4], carry)
	z[5], _ = bits.Add64(z[5], z[5], carry)

	
	
	if !(z[5] < 1873798617647539866 || (z[5] == 1873798617647539866 && (z[4] < 5412103778470702295 || (z[4] == 5412103778470702295 && (z[3] < 7239337960414712511 || (z[3] == 7239337960414712511 && (z[2] < 7435674573564081700 || (z[2] == 7435674573564081700 && (z[1] < 2210141511517208575 || (z[1] == 2210141511517208575 && (z[0] < 13402431016077863595))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 13402431016077863595, 0)
		z[1], b = bits.Sub64(z[1], 2210141511517208575, b)
		z[2], b = bits.Sub64(z[2], 7435674573564081700, b)
		z[3], b = bits.Sub64(z[3], 7239337960414712511, b)
		z[4], b = bits.Sub64(z[4], 5412103778470702295, b)
		z[5], _ = bits.Sub64(z[5], 1873798617647539866, b)
	}
}

func ldouble(z, x *fe) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], carry = bits.Add64(x[2], x[2], carry)
	z[3], carry = bits.Add64(x[3], x[3], carry)
	z[4], carry = bits.Add64(x[4], x[4], carry)
	z[5], _ = bits.Add64(x[5], x[5], carry)
}

func sub(z, x, y *fe) {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	z[1], b = bits.Sub64(x[1], y[1], b)
	z[2], b = bits.Sub64(x[2], y[2], b)
	z[3], b = bits.Sub64(x[3], y[3], b)
	z[4], b = bits.Sub64(x[4], y[4], b)
	z[5], b = bits.Sub64(x[5], y[5], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 13402431016077863595, 0)
		z[1], c = bits.Add64(z[1], 2210141511517208575, c)
		z[2], c = bits.Add64(z[2], 7435674573564081700, c)
		z[3], c = bits.Add64(z[3], 7239337960414712511, c)
		z[4], c = bits.Add64(z[4], 5412103778470702295, c)
		z[5], _ = bits.Add64(z[5], 1873798617647539866, c)
	}
}

func subAssign(z, x *fe) {
	var b uint64
	z[0], b = bits.Sub64(z[0], x[0], 0)
	z[1], b = bits.Sub64(z[1], x[1], b)
	z[2], b = bits.Sub64(z[2], x[2], b)
	z[3], b = bits.Sub64(z[3], x[3], b)
	z[4], b = bits.Sub64(z[4], x[4], b)
	z[5], b = bits.Sub64(z[5], x[5], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 13402431016077863595, 0)
		z[1], c = bits.Add64(z[1], 2210141511517208575, c)
		z[2], c = bits.Add64(z[2], 7435674573564081700, c)
		z[3], c = bits.Add64(z[3], 7239337960414712511, c)
		z[4], c = bits.Add64(z[4], 5412103778470702295, c)
		z[5], _ = bits.Add64(z[5], 1873798617647539866, c)
	}
}

func lsubAssign(z, x *fe) {
	var b uint64
	z[0], b = bits.Sub64(z[0], x[0], 0)
	z[1], b = bits.Sub64(z[1], x[1], b)
	z[2], b = bits.Sub64(z[2], x[2], b)
	z[3], b = bits.Sub64(z[3], x[3], b)
	z[4], b = bits.Sub64(z[4], x[4], b)
	z[5], b = bits.Sub64(z[5], x[5], b)
}

func neg(z *fe, x *fe) {
	if x.isZero() {
		z.zero()
		return
	}
	var borrow uint64
	z[0], borrow = bits.Sub64(13402431016077863595, x[0], 0)
	z[1], borrow = bits.Sub64(2210141511517208575, x[1], borrow)
	z[2], borrow = bits.Sub64(7435674573564081700, x[2], borrow)
	z[3], borrow = bits.Sub64(7239337960414712511, x[3], borrow)
	z[4], borrow = bits.Sub64(5412103778470702295, x[4], borrow)
	z[5], _ = bits.Sub64(1873798617647539866, x[5], borrow)
}

func mul(z, x, y *fe) {
	var t [6]uint64
	var c [3]uint64
	{
		
		v := x[0]
		c[1], c[0] = bits.Mul64(v, y[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd1(v, y[1], c[1])
		c[2], t[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd1(v, y[2], c[1])
		c[2], t[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd1(v, y[3], c[1])
		c[2], t[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd1(v, y[4], c[1])
		c[2], t[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd1(v, y[5], c[1])
		t[5], t[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}
	{
		
		v := x[1]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		t[5], t[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}
	{
		
		v := x[2]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		t[5], t[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}
	{
		
		v := x[3]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		t[5], t[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}
	{
		
		v := x[4]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		t[5], t[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}
	{
		
		v := x[5]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9940570264628428797
		c[2] = madd0(m, 13402431016077863595, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], z[0] = madd2(m, 2210141511517208575, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], z[1] = madd2(m, 7435674573564081700, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], z[2] = madd2(m, 7239337960414712511, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], z[3] = madd2(m, 5412103778470702295, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		z[5], z[4] = madd3(m, 1873798617647539866, c[0], c[2], c[1])
	}

	
	
	if !(z[5] < 1873798617647539866 || (z[5] == 1873798617647539866 && (z[4] < 5412103778470702295 || (z[4] == 5412103778470702295 && (z[3] < 7239337960414712511 || (z[3] == 7239337960414712511 && (z[2] < 7435674573564081700 || (z[2] == 7435674573564081700 && (z[1] < 2210141511517208575 || (z[1] == 2210141511517208575 && (z[0] < 13402431016077863595))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 13402431016077863595, 0)
		z[1], b = bits.Sub64(z[1], 2210141511517208575, b)
		z[2], b = bits.Sub64(z[2], 7435674573564081700, b)
		z[3], b = bits.Sub64(z[3], 7239337960414712511, b)
		z[4], b = bits.Sub64(z[4], 5412103778470702295, b)
		z[5], _ = bits.Sub64(z[5], 1873798617647539866, b)
	}
}

func square(z, x *fe) {

	var p [6]uint64

	var u, v uint64
	{
		
		u, p[0] = bits.Mul64(x[0], x[0])
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		var t uint64
		t, u, v = madd1sb(x[0], x[1], u)
		C, p[0] = madd2(m, 2210141511517208575, v, C)
		t, u, v = madd1s(x[0], x[2], t, u)
		C, p[1] = madd2(m, 7435674573564081700, v, C)
		t, u, v = madd1s(x[0], x[3], t, u)
		C, p[2] = madd2(m, 7239337960414712511, v, C)
		t, u, v = madd1s(x[0], x[4], t, u)
		C, p[3] = madd2(m, 5412103778470702295, v, C)
		_, u, v = madd1s(x[0], x[5], t, u)
		p[5], p[4] = madd3(m, 1873798617647539866, v, C, u)
	}
	{
		
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		u, v = madd1(x[1], x[1], p[1])
		C, p[0] = madd2(m, 2210141511517208575, v, C)
		var t uint64
		t, u, v = madd2sb(x[1], x[2], p[2], u)
		C, p[1] = madd2(m, 7435674573564081700, v, C)
		t, u, v = madd2s(x[1], x[3], p[3], t, u)
		C, p[2] = madd2(m, 7239337960414712511, v, C)
		t, u, v = madd2s(x[1], x[4], p[4], t, u)
		C, p[3] = madd2(m, 5412103778470702295, v, C)
		_, u, v = madd2s(x[1], x[5], p[5], t, u)
		p[5], p[4] = madd3(m, 1873798617647539866, v, C, u)
	}
	{
		
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		C, p[0] = madd2(m, 2210141511517208575, p[1], C)
		u, v = madd1(x[2], x[2], p[2])
		C, p[1] = madd2(m, 7435674573564081700, v, C)
		var t uint64
		t, u, v = madd2sb(x[2], x[3], p[3], u)
		C, p[2] = madd2(m, 7239337960414712511, v, C)
		t, u, v = madd2s(x[2], x[4], p[4], t, u)
		C, p[3] = madd2(m, 5412103778470702295, v, C)
		_, u, v = madd2s(x[2], x[5], p[5], t, u)
		p[5], p[4] = madd3(m, 1873798617647539866, v, C, u)
	}
	{
		
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		C, p[0] = madd2(m, 2210141511517208575, p[1], C)
		C, p[1] = madd2(m, 7435674573564081700, p[2], C)
		u, v = madd1(x[3], x[3], p[3])
		C, p[2] = madd2(m, 7239337960414712511, v, C)
		var t uint64
		t, u, v = madd2sb(x[3], x[4], p[4], u)
		C, p[3] = madd2(m, 5412103778470702295, v, C)
		_, u, v = madd2s(x[3], x[5], p[5], t, u)
		p[5], p[4] = madd3(m, 1873798617647539866, v, C, u)
	}
	{
		
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		C, p[0] = madd2(m, 2210141511517208575, p[1], C)
		C, p[1] = madd2(m, 7435674573564081700, p[2], C)
		C, p[2] = madd2(m, 7239337960414712511, p[3], C)
		u, v = madd1(x[4], x[4], p[4])
		C, p[3] = madd2(m, 5412103778470702295, v, C)
		_, u, v = madd2sb(x[4], x[5], p[5], u)
		p[5], p[4] = madd3(m, 1873798617647539866, v, C, u)
	}
	{
		
		m := p[0] * 9940570264628428797
		C := madd0(m, 13402431016077863595, p[0])
		C, z[0] = madd2(m, 2210141511517208575, p[1], C)
		C, z[1] = madd2(m, 7435674573564081700, p[2], C)
		C, z[2] = madd2(m, 7239337960414712511, p[3], C)
		C, z[3] = madd2(m, 5412103778470702295, p[4], C)
		u, v = madd1(x[5], x[5], p[5])
		z[5], z[4] = madd3(m, 1873798617647539866, v, C, u)
	}

	
	
	if !(z[5] < 1873798617647539866 || (z[5] == 1873798617647539866 && (z[4] < 5412103778470702295 || (z[4] == 5412103778470702295 && (z[3] < 7239337960414712511 || (z[3] == 7239337960414712511 && (z[2] < 7435674573564081700 || (z[2] == 7435674573564081700 && (z[1] < 2210141511517208575 || (z[1] == 2210141511517208575 && (z[0] < 13402431016077863595))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 13402431016077863595, 0)
		z[1], b = bits.Sub64(z[1], 2210141511517208575, b)
		z[2], b = bits.Sub64(z[2], 7435674573564081700, b)
		z[3], b = bits.Sub64(z[3], 7239337960414712511, b)
		z[4], b = bits.Sub64(z[4], 5412103778470702295, b)
		z[5], _ = bits.Sub64(z[5], 1873798617647539866, b)
	}
}


















func madd(a, b, t, u, v uint64) (uint64, uint64, uint64) {
	var carry uint64
	hi, lo := bits.Mul64(a, b)
	v, carry = bits.Add64(lo, v, 0)
	u, carry = bits.Add64(hi, u, carry)
	t, _ = bits.Add64(t, 0, carry)
	return t, u, v
}


func madd0(a, b, c uint64) (hi uint64) {
	var carry, lo uint64
	hi, lo = bits.Mul64(a, b)
	_, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}


func madd1(a, b, c uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}


func madd2(a, b, c, d uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}


func madd2s(a, b, c, d, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}

func madd1s(a, b, d, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}

func madd2sb(a, b, c, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

func madd1sb(a, b, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

func madd3(a, b, c, d, e uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, e, carry)
	return
}
