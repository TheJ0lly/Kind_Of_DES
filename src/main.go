package main

import (
	"flag"
	"fmt"
)

var RC *RoundComputer = &RoundComputer{Rounds: []int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}}

var Decrypt *bool
var Encrypt *bool

func LeftRotateKey(key *Bitset, round int) *Bitset {
	l, r := key.Split()
	l.ShiftBy(RC.Rounds[round], true)
	r.ShiftBy(RC.Rounds[round], true)

	rotKey := ConcatBitsets(l, r)

	return rotKey
}

func PrecomputeRounds(flk, frk *Bitset) []*Bitset {
	roundKeys := make([]*Bitset, 16)
	firstkey := ConcatBitsets(flk, frk)

	l, r := firstkey.Split()

	// We append the first key, which is the original key + leftshift[Round].
	roundKeys[0] = LeftRotateKey(ConcatBitsets(l, r), RC.Rounds[0])

	// We create all other rounds keys.
	for r := 1; r < 16; r++ {
		roundKeys[r] = LeftRotateKey(roundKeys[r-1], r)
	}

	// We apply the PC2 permutation on each key.
	for r := 0; r < 16; r++ {
		roundKeys[r].Permute(&PC2)
	}

	return roundKeys
}

func TransformKey(left *Bitset, right *Bitset, RoundIndex int) *Bitset {
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

	outBS := CreateBitsetFromUInt32(out)
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

	if len(*key) != 16 {
		logError("key must be exactly 16 characters\n")
		return
	}

	if *Encrypt && len(*text) != 8 {
		logError("encryption: text must be exactly 8 characters\n")
		return
	}

	if *Decrypt && len(*text) != 32 {
		logError("decryption: text must be exactly 32 HEX characters\n")
		return
	}

	var textBS, keyBS, Temp *Bitset

	if *Encrypt {
		// Prepare the text
		textBS = CreateBitsetFromString(*text, false)
		// Adding timestamp (0 FOR NOW)
		textBS.Salt()
		textBS.Permute(&IP)
		// 64 bit halves - we need to split them once each.
		TL, TR := textBS.Split()

		// Left side halves
		TL0, TR0 := TL.Split()

		// Right side halves
		TL1, TR1 := TR.Split()

		// Prepare the key
		keyBS = CreateBitsetFromString(*key, false)
		keyBS.Permute(&PC1)

		// We get the 56 bits sides
		KL, KR := keyBS.Split()

		// The 28 bit halves of the left side
		KL0, KR0 := KL.Split()
		// Compute the key rounds for encryption of the left side
		KeyRounds0 := PrecomputeRounds(KL0, KR0)

		// The 28 bit halves of the right side
		KL1, KR1 := KR.Split()
		// Compute the key rounds for encryption of the right side
		KeyRounds1 := PrecomputeRounds(KL1, KR1)

		for RoundIndex := 0; RoundIndex < 16; RoundIndex++ {
			//============== LEFT SIDE =================
			// The new right side
			Temp = XORBitsets(Feistel(KeyRounds0[RoundIndex], TR0), TL0)
			// The left side is the old right side
			TL0 = TR0
			// The right side is the new computed side.
			TR0 = Temp
			//===========================================

			//============== RIGHT SIDE =================
			// The new right side
			Temp = XORBitsets(Feistel(KeyRounds1[RoundIndex], TR1), TL1)
			// The left side is the old right side
			TL1 = TR1
			// The right side is the new computed side.
			TR1 = Temp
			//===========================================

		}

		// We switch the sides, as after 16 rounds they are inverted.
		FinalLeft := ConcatBitsets(TR0, TL0)
		FinalRight := ConcatBitsets(TR1, TL1)

		Final := ConcatBitsets(FinalLeft, FinalRight)
		Final.Permute(IP.Inverse())

		fmt.Printf("Output: %s\n", Final.ToHexString())
	} else {
		// Prepare the text
		textBS = CreateBitsetFromString(*text, true)
		textBS.Permute(&IP)
		// 64 bit halves - we need to split them once each.
		TL, TR := textBS.Split()

		// Left side halves
		TL0, TR0 := TL.Split()

		// Right side halves
		TL1, TR1 := TR.Split()

		// Prepare the key
		keyBS = CreateBitsetFromString(*key, false)
		keyBS.Permute(&PC1)

		// We get the 56 bits sides
		KL, KR := keyBS.Split()

		// The 28 bit halves of the left side
		KL0, KR0 := KL.Split()
		// Compute the key rounds for encryption of the left side
		KeyRounds0 := PrecomputeRounds(KL0, KR0)

		// The 28 bit halves of the right side
		KL1, KR1 := KR.Split()
		// Compute the key rounds for encryption of the right side
		KeyRounds1 := PrecomputeRounds(KL1, KR1)

		for RoundIndex := 15; RoundIndex >= 0; RoundIndex-- {
			//============== LEFT SIDE =================
			// The new right side
			Temp = XORBitsets(Feistel(KeyRounds0[RoundIndex], TR0), TL0)
			// The left side is the old right side
			TL0 = TR0
			// The right side is the new computed side.
			TR0 = Temp
			//===========================================

			//============== RIGHT SIDE =================
			// The new right side
			Temp = XORBitsets(Feistel(KeyRounds1[RoundIndex], TR1), TL1)
			// The left side is the old right side
			TL1 = TR1
			// The right side is the new computed side.
			TR1 = Temp
			//===========================================

		}

		// We switch the sides, as after 16 rounds they are inverted.
		FinalLeft := ConcatBitsets(TR0, TL0)
		FinalRight := ConcatBitsets(TR1, TL1)

		Final := ConcatBitsets(FinalLeft, FinalRight)
		Final.Permute(IP.Inverse())
		Final.RemoveSalt()

		fmt.Printf("Output: %s\n", Final.ToString())
	}
}
