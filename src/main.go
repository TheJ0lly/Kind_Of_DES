package main

import (
	"flag"
	"fmt"
	"slices"
)

var RC *RoundComputer = &RoundComputer{Rounds: make([]int, 16)}

var Decrypt *bool
var Encrypt *bool

func GenerateRounds(key *Bitset) {
	for i := 0; i < key.Len()-7; i++ {
		part := key.Bits[i : i+7]
		val := byte(0)

		for j := 0; j < 7; j++ {
			val <<= 1
			val |= part[j]
		}

		RC.Rounds[val%16]++
	}

	maxRounds := 0

	roundc := make([]int, 16)
	copy(roundc, RC.Rounds)

	// We set the rounds of 1 bit shift
	for maxRounds < 4 {
		index := slices.Index(roundc, slices.Max(roundc))

		RC.Rounds[index] = 1

		roundc = slices.Delete(roundc, index, index+1)

		maxRounds++
	}

	// Everything that's not a 1, we set to 2.
	for i := 0; i < len(RC.Rounds); i++ {
		if RC.Rounds[i] != 1 {
			RC.Rounds[i] = 2
		}
	}
}

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

	LeftChan := make(chan struct{})
	RightChan := make(chan struct{})

	LeftRound := make(chan int)
	RightRound := make(chan int)

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

		// We generate the rounds based on the key
		GenerateRounds(keyBS)

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

		//============== LEFT SIDE =================
		go func(ready chan struct{}, round chan int) {
			for {
				RoundIndex := <-round

				if RoundIndex == 16 {
					// We send the last signal to kill the synchronizer
					ready <- struct{}{}
					break
				}

				// The new right side
				Temp = XORBitsets(Feistel(KeyRounds0[RoundIndex], TR0), TL0)
				// The left side is the old right side
				TL0 = TR0
				// The right side is the new computed side.
				TR0 = Temp

				// Signaling that the left side is done for this round.
				ready <- struct{}{}
			}
		}(LeftChan, LeftRound)
		//===========================================

		//============== RIGHT SIDE =================
		go func(ready chan struct{}, round chan int) {
			for {
				RoundIndex := <-round

				if RoundIndex == 16 {
					// We send the last signal to kill the synchronizer
					ready <- struct{}{}
					break
				}

				// The new left side
				Temp = XORBitsets(Feistel(KeyRounds1[RoundIndex], TL1), TR1)
				// The right side is the old left side
				TR1 = TL1
				// The left side is the new computed side.
				TL1 = Temp

				// Signaling that the right side is done for this round.
				ready <- struct{}{}
			}
		}(RightChan, RightRound)
		//===========================================

		Done := make(chan struct{})

		// Synchronizer
		go func(kill, leftReady, rightReady chan struct{}, leftRound, rightRound chan int) {

			for i := 0; i < 17; i++ {
				// We send the rount through the respective channels
				leftRound <- i
				rightRound <- i

				// Waiting for signals to proces
				<-leftReady
				<-rightReady

				// We switch the rounds after each process step to shuffle everything more.
				KeyRounds0, KeyRounds1 = KeyRounds1, KeyRounds0
			}

			kill <- struct{}{}
		}(Done, LeftChan, RightChan, LeftRound, RightRound)

		// We wait for the synchronizer to die.
		<-Done

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

		// We generate the rounds based on the key
		GenerateRounds(keyBS)

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

		//============== LEFT SIDE =================
		go func(ready chan struct{}, round chan int) {
			for {
				RoundIndex := <-round

				if RoundIndex == -1 {
					// We send the last signal to kill the synchronizer
					ready <- struct{}{}
					break
				}

				// The new right side
				Temp = XORBitsets(Feistel(KeyRounds0[RoundIndex], TR0), TL0)
				// The left side is the old right side
				TL0 = TR0
				// The right side is the new computed side.
				TR0 = Temp

				// Signaling that the left side is done for this round.
				ready <- struct{}{}
			}
		}(LeftChan, LeftRound)
		//===========================================

		//============== RIGHT SIDE =================
		go func(ready chan struct{}, round chan int) {
			for {
				RoundIndex := <-round

				if RoundIndex == -1 {
					// We send the last signal to kill the synchronizer
					ready <- struct{}{}
					break
				}

				// The new left side
				Temp = XORBitsets(Feistel(KeyRounds1[RoundIndex], TL1), TR1)
				// The right side is the old left side
				TR1 = TL1
				// The left side is the new computed side.
				TL1 = Temp

				// Signaling that the right side is done for this round.
				ready <- struct{}{}
			}
		}(RightChan, RightRound)
		//===========================================

		Done := make(chan struct{})

		// Synchronizer
		go func(kill, leftReady, rightReady chan struct{}, leftRound, rightRound chan int) {

			for i := 15; i >= -1; i-- {
				// We switch the rounds before making any operation to match to reverse order of the rounds
				KeyRounds0, KeyRounds1 = KeyRounds1, KeyRounds0

				// We send the round through the respective channels
				leftRound <- i
				rightRound <- i

				// Waiting for signals to proces
				<-leftReady
				<-rightReady

			}

			kill <- struct{}{}
		}(Done, LeftChan, RightChan, LeftRound, RightRound)

		// We wait for the synchronizer to die.
		<-Done

		// We switch the sides, as after 16 rounds they are inverted.
		FinalLeft := ConcatBitsets(TR0, TL0)
		FinalRight := ConcatBitsets(TR1, TL1)

		Final := ConcatBitsets(FinalLeft, FinalRight)
		Final.Permute(IP.Inverse())
		Final.RemoveSalt()

		fmt.Printf("Output: %s\n", Final.ToString())
	}
}
