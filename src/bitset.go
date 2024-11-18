package main

import (
	"fmt"
	"strconv"
	"time"
)

type Bitset struct {
	Bits []byte
}

// CreateBitsetFromString creates a Bitset from a string.
func CreateBitsetFromString(str string, hex bool) *Bitset {
	// We allocate space for the bits of the string + the timestamp bits(64 bits),
	// and we set the size to the length of the string * 8, to get all bits per character.
	var bs *Bitset

	bsI := 0

	if hex {
		bs = &Bitset{Bits: make([]byte, len(str)*4)}
		for i := 0; i < len(str); i += 2 {
			b := GetByteFromHex(str[i], str[i+1])
			// We initialize with 128, because this is the value of the 8 bit alone.
			bit := byte(128)

			for i := 7; i >= 0; i-- {
				val := byte(b) & bit >> i
				bs.Bits[bsI] = val

				bit /= 2
				bsI++
			}
		}
	} else {
		bs = &Bitset{Bits: make([]byte, len(str)*8)}
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
	}

	return bs
}

func CreateBitsetFromUInt32(num uint32) *Bitset {
	// We allocate space for the bits of the string + the timestamp bits(64 bits),
	// and we set the size to the length of the string * 8, to get all bits per character.

	bs := &Bitset{Bits: make([]byte, 32)}
	bsI := 0

	for i := 31; i >= 0; i-- {
		bs.Bits[bsI] = byte((num >> i) & 1)
		bsI++
	}

	return bs
}

func CreateBitsetFromInt64(num int64) *Bitset {
	bs := &Bitset{Bits: make([]byte, 64)}
	bsI := 0

	for i := 63; i >= 0; i-- {
		bs.Bits[bsI] = byte((num >> i) & 1)
		bsI++
	}

	return bs
}

func CopyBitset(bs *Bitset) *Bitset {
	ret := &Bitset{Bits: make([]byte, len(bs.Bits))}

	copy(ret.Bits, bs.Bits)

	return ret
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

	left := &Bitset{Bits: make([]byte, middle)}
	right := &Bitset{Bits: make([]byte, middle)}

	copy(left.Bits, bs.Bits[:middle])
	copy(right.Bits, bs.Bits[middle:])

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

	uxTDup := (uxT << 32) | uxT

	uxTBS := CreateBitsetFromInt64(uxTDup)

	bitIndex := 0

	// We default it to 128, because salting, will get us to a 128 bit bitset
	nbs := make([]byte, 128)

	for i := 0; i < len(bs.Bits); i++ {
		nbs[bitIndex] = uxTBS.Bits[i]
		nbs[bitIndex+1] = bs.Bits[i]
		bitIndex += 2
	}

	bs.Bits = nbs
	logInfo("inserted timestamp: current length %d\n", bs.Len())
}

func (bs *Bitset) RemoveSalt() {
	logInfo("removing salt\n")

	nbs := make([]byte, 64)
	bitIndex := 0

	for i := 1; i < len(bs.Bits); i += 2 {
		nbs[bitIndex] = bs.Bits[i]
		bitIndex++
	}

	bs.Bits = nbs
}

// Permute applies a permutation to the bitset.
func (bs *Bitset) Permute(p *Permutation) {
	nbs := make([]byte, len(p.Data))

	for i := 0; i < len(p.Data); i++ {
		nbs[i] = bs.Bits[p.Data[i]-1]
	}

	bs.Bits = nbs
	logInfo("new length after permute %d\n", bs.Len())
}

// ShiftBy shifts the Bitset by `n` bits.
func (bs *Bitset) ShiftBy(n int, direction bool) {
	// If `direction` is true, we shift left.
	if direction {
		val := []byte{bs.Bits[0]}
		if n == 2 {
			val = append(val, bs.Bits[1])
		}

		bs.Bits = bs.Bits[n:]
		bs.Bits = append(bs.Bits, val...)
	} else {
		// Otherwise we shift right
		val := []byte{bs.Bits[bs.Len()-1]}

		if n == 2 {
			val = append(val, bs.Bits[bs.Len()-2])
		}

		bs.Bits = bs.Bits[:bs.Len()-n]
		bs.Bits = append(val, bs.Bits...)
	}

}

func ConcatBitsets(l *Bitset, r *Bitset) *Bitset {

	left := make([]byte, l.Len())
	right := make([]byte, l.Len())

	copy(left, l.Bits)
	copy(right, r.Bits)

	ret := &Bitset{Bits: append(left, right...)}
	logInfo("concatenated bitsets: new length %d\n", ret.Len())
	return ret
}

func XORBitsets(l *Bitset, r *Bitset) *Bitset {
	ret := &Bitset{Bits: make([]byte, l.Len())}

	for i := 0; i < l.Len(); i++ {
		ret.Bits[i] = l.Bits[i] ^ r.Bits[i]
	}

	return ret
}

func Get6BitVal(bits []byte) byte {
	sum := byte(0)

	for i := 0; i < len(bits); i++ {
		sum = (sum << 1) | bits[i]
	}

	return sum
}

func (bs *Bitset) ApplySBox(which int) uint32 {
	index := 6 * which

	// We get the 6 bits before applying the SBox
	val := Get6BitVal(bs.Bits[index : index+6])

	// We get the first bit, then we get the 6th bit
	row := ((val & 32) >> 4) | (val & 1)

	// We nullify the first and last bit, and we get the middle 4
	col := (val & 30) >> 1

	return uint32(SBoxes[which][row][col])
}

func GetByteAsHex(bits []byte) string {
	char := byte(0)

	HexMap := map[byte]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
	}

	fstr := ""
	for i := 0; i < len(bits); i++ {
		char = (char << 1) | bits[i]
	}

	fstr += HexMap[(char&240)>>4]
	fstr += HexMap[char&15]

	return fstr
}

// GetByteFromHex will return a byte formed out of 2 hex characters
func GetByteFromHex(first byte, second byte) byte {
	var final byte
	ASCIIToHexDifference := byte(55)

	if first >= '0' && first <= '9' {
		num, _ := strconv.Atoi(string(first))

		final = byte(num) << 4
	} else {
		final = (first - ASCIIToHexDifference) << 4
	}

	if second >= '0' && second <= '9' {
		num, _ := strconv.Atoi(string(second))

		final = final | byte(num)
	} else {
		final = final | (second - ASCIIToHexDifference)
	}

	return final

}

func (bs *Bitset) ToHexString() string {
	fstr := ""

	for i := 0; i < bs.Len(); i += 8 {
		fstr += GetByteAsHex(bs.Bits[i : i+8])
	}

	return fstr
}

func (bs *Bitset) ToString() string {
	fstr := ""

	for i := 0; i < bs.Len(); i += 8 {
		sum := byte(0)

		for j := i; j < i+8; j++ {
			sum = (sum << 1) | bs.Bits[j]
		}

		fstr += string(sum)
	}

	return fstr
}
