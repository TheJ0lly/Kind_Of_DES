package main

import (
	"fmt"
	"slices"
	"time"
)

type Bitset struct {
	Bits []byte
}

// CreateBitset creates a Bitset from a string.
func CreateBitset(str string, double bool) *Bitset {
	// We allocate space for the bits of the string + the timestamp bits(64 bits),
	// and we set the size to the length of the string * 8, to get all bits per character.

	capVal := 1

	if double {
		capVal = 2
	}

	bs := &Bitset{Bits: make([]byte, len(str)*8, (len(str)*8)*capVal)}
	bsI := 0

	for _, char := range str {
		// We initialize with 128, because this is the value of the 8 bit alone.
		bit := byte(128)

		// We get each bit of the current character in the string.
		for i := 7; i >= 0; i-- {
			val := byte(char) & bit >> i
			bs.Bits[bsI] = val

			bit /= 2
			bsI++
		}
	}

	return bs
}

// LenBytes returns the number of full bytes.
func (bs *Bitset) LenBytes() int {
	return len(bs.Bits) / 8
}

// Len returns the number of bits.
func (bs *Bitset) Len() int {
	return len(bs.Bits)
}

// RemoveParityBits will remove the parity bits of the bitset.
//
// Should only be used on the key one time before the process.
func (bs *Bitset) RemoveParityBits() {
	logInfo("removing parity bits: current length %d -> expected length %d\n", bs.Len(), bs.Len()-bs.LenBytes())
	// We do minus 9 because we offset 1 from the length of the bits,
	// and 8 because we need the last parity bit
	i := bs.Len() - 8

	for i >= 0 {
		bs.Bits = append(bs.Bits[:i], bs.Bits[i+1:]...)
		i -= 8
	}

	logInfo("parity bits removed: current length %d\n", bs.Len())
}

// Split creates 2 equal halves of a bitset.
func (bs *Bitset) Split() (*Bitset, *Bitset) {
	middle := bs.Len() / 2
	logInfo("spliting bitset into 2 halves of %d bits each\n", middle)

	left := &Bitset{Bits: bs.Bits[:middle]}
	right := &Bitset{Bits: bs.Bits[middle:]}

	return left, right
}

func (bs *Bitset) Print() {
	for i := 0; i < bs.Len(); i++ {
		fmt.Printf("%d", bs.Bits[i])
	}
	fmt.Println()
}

// Salt adds the timestamp bits into the text bits to give 128 bits.
func (bs *Bitset) Salt() {
	logInfo("inserting timestamp 64 bits into bitset: length %d --- capacity %d\n", bs.Len(), cap(bs.Bits))

	uxT := time.Now().UnixMilli()
	bit := 0

	// We insert the timestamp bits into even places.
	for i := 0; i < cap(bs.Bits); i += 2 {
		val := byte((uxT >> bit) & 0x01)
		bs.Bits = slices.Insert(bs.Bits, i, val)
		bit++
	}

	logInfo("inserted timestamp: current length %d\n", bs.Len())
}

// Permute applies a permutation to the bitset.
func (bs *Bitset) Permute(p *Permutation) {
	nbs := make([]byte, bs.Len())

	for i := 0; i < bs.Len(); i++ {
		nbs[i] = bs.Bits[p.Data[i]-1]
	}

	bs.Bits = nbs
}

// ShiftBy shifts the Bitset by `n` bits.
func (bs *Bitset) ShiftBy(n int) {
	val := []byte{bs.Bits[0]}
	if n == 2 {
		val = append(val, bs.Bits[1])
	}

	bs.Bits = bs.Bits[n:]
	bs.Bits = append(bs.Bits, val...)
}

func ConcatBitsets(l *Bitset, r *Bitset) *Bitset {
	return &Bitset{Bits: append(l.Bits, r.Bits...)}
}
