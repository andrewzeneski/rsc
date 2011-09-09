// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gf256

import (
	"bytes"
	"fmt"
	"testing"
)

var f = NewField(0x11d) // x^8 + x^4 + x^3 + x^2 + 1

func TestBasic(t *testing.T) {
	if f.Exp(0) != 1 || f.Exp(1) != 2 || f.Exp(255) != 1 {
		panic("bad Exp")
	}
}

func TestEncode(t *testing.T) {
	data := []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11}
	check := []byte{0xa5, 0x24, 0xd4, 0xc1, 0xed, 0x36, 0xc7, 0x87, 0x2c, 0x55}
	out := f.ECBytes(data, len(check))
	if !bytes.Equal(out, check) {
		t.Errorf("have %x want %x", out, check)
	}
}

func TestLinear(t *testing.T) {
	d1 := []byte{0x00, 0x00}
	c1 := []byte{0x00, 0x00}
	if out := f.ECBytes(d1, len(c1)); !bytes.Equal(out, c1) {
		t.Errorf("ECBytes(%x, %d) = %x, want 0", d1, len(c1), out)
	}
	d2 := []byte{0x00, 0x01}
	c2 := f.ECBytes(d2, 2)
	d3 := []byte{0x00, 0x02}
	c3 := f.ECBytes(d3, 2)
	cx := make([]byte, 2)
	for i := range cx {
		cx[i] = c2[i] ^ c3[i]
	}
	d4 := []byte{0x00, 0x03}
	c4 := f.ECBytes(d4, 2)
	if !bytes.Equal(cx, c4) {
		t.Errorf("ECBytes(%x, 2) = %x\nECBytes(%x, 2) = %x\nxor = %x\nECBytes(%x, 2) = %x",
			d2, c2, d3, c3, cx, d4, c4)
	}
}

func TestGaussJordan(t *testing.T) {

	m := make([][]byte, 16)
	for i := range m {
		m[i] = make([]byte, 4)
		m[i][i/8] = 1 << uint(i%8)
		copy(m[i][2:], f.ECBytes(m[i][:2], 2))
	}
	fmt.Printf("---\n")
	for _, row := range m {
		fmt.Printf("%x\n", row)
	}
	b := []uint{0, 1, 2, 3, 12, 13, 14, 15, 20, 21, 22, 23, 24, 25, 26, 27}
	for i := 0; i < 16; i++ {
		bi := b[i]
		if m[i][bi/8]&(1<<(7-bi%8)) == 0 {
			for j := i + 1; ; j++ {
				if j >= len(m) {
					t.Errorf("lost track for %d", bi)
					break
				}
				if m[j][bi/8]&(1<<(7-bi%8)) != 0 {
					m[i], m[j] = m[j], m[i]
					break
				}
			}
		}
		for j := i + 1; j < len(m); j++ {
			if m[j][bi/8]&(1<<(7-bi%8)) != 0 {
				for k := range m[j] {
					m[j][k] ^= m[i][k]
				}
			}
		}
	}
	fmt.Printf("---\n")
	for _, row := range m {
		fmt.Printf("%x\n", row)
	}
	for i := 15; i >= 0; i-- {
		bi := b[i]
		for j := i - 1; j >= 0; j-- {
			if m[j][bi/8]&(1<<(7-bi%8)) != 0 {
				for k := range m[j] {
					m[j][k] ^= m[i][k]
				}
			}
		}
	}
	fmt.Printf("---\n")
	for _, row := range m {
		fmt.Printf("%x", row)
		if out := f.ECBytes(row[:2], 2); !bytes.Equal(out, row[2:]) {
			fmt.Printf(" - want %x", out)
		}
		fmt.Printf("\n")
	}
}
