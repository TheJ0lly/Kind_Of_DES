package main

import (
	"flag"
	"fmt"
)

var RoundIndex = 0
var TEXT_LEFT *Bitset
var TEXT_RIGHT *Bitset
var KEY_LEFT *Bitset
var KEY_RIGHT *Bitset
var RC *RoundComputer = &RoundComputer{Rounds: []int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}}

var Decrypt *bool
var Encrypt *bool

func TransformKey(left *Bitset, right *Bitset) *Bitset {
	left.ShiftBy(RC.Rounds[RoundIndex], true)
	right.ShiftBy(RC.Rounds[RoundIndex], true)

	ck := ConcatBitsets(left, right)
	ck.Permute(&PC2)

	return ck
}

func Feistel(key *Bitset, right *Bitset) *Bitset {
	// Do the expansion permutation
	rc := CopyBitset(right)
	rc.Permute(&Expansion)

	xor := XORBitsets(key, rc)

	out := uint32(0)
	for i := 0; i < 8; i++ {
		out <<= 4
		v := xor.ApplySBox(i)
		out |= v
	}

	outBS := CreateBitsetFromUInt32(out, false)
	outBS.Permute(&PPermutation)

	return outBS
}

func main() {
	key := flag.String("k", "", "the key used for encryption/decryption")
	text := flag.String("t", "", "the text to be encrypted/decrypted")
	Decrypt = flag.Bool("d", false, "starts the decryption process")
	Encrypt = flag.Bool("e", false, "starts the encryption process")
	dbg := flag.Bool("l", false, "toggles the logs")
	help := flag.Bool("h", false, "help menu")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	CanLog = *dbg

	if *Decrypt && *Encrypt {
		logError("both the encryption and decryption flags have been toggled - must have only 1\n")
		return
	}

	if !*Decrypt && !*Encrypt {
		logError("error: neither encryption nor decryption flag has been toggled - must have 1\n")
		return
	}

	if len(*key) != 8 {
		logError("key must be exactly 8 characters\n")
		return
	}

	if len(*text) != 8 {
		logError("text must be exactly 8 characters\n")
		return
	}

	keyBS := CreateBitsetFromString(*key, false, false)

	// Drop the parity bits - 56 bit key
	keyBS.Permute(&PC1)

	KEY_LEFT, KEY_RIGHT = keyBS.Split()

	textBS := CreateBitsetFromString(*text, false, false)

	textBS.Permute(&IP)

	L0, R0 := textBS.Split()
	var Temp *Bitset

	if *Encrypt {
		RoundIndex = 0
		for ; RoundIndex < 16; RoundIndex++ {
			// We concatenate the key
			tk := TransformKey(KEY_LEFT, KEY_RIGHT)

			// The new right side
			Temp = XORBitsets(Feistel(tk, R0), L0)

			// The left side is the old right side
			L0 = R0

			// The right side is the new computed side.
			R0 = Temp
		}

		// We switch the sides, as after 16 rounds they are inverted.
		Final := ConcatBitsets(R0, L0)
		Final.Permute(IP.Inverse())

		fmt.Printf("Output: %s\n", Final.ToHexString())
	}
}
