package main

import (
	sdes "TheJ0lly/SDES/SDES"
	"flag"
	"fmt"
)

var Decrypt *bool
var Encrypt *bool

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

	sdes.CanLog = *dbg

	if *Decrypt && *Encrypt {
		sdes.LogError("both the encryption and decryption flags have been toggled - must have only 1\n")
		return
	}

	if !*Decrypt && !*Encrypt {
		sdes.LogError("error: neither encryption nor decryption flag has been toggled - must have 1\n")
		return
	}

	if len(*key) != 16 {
		sdes.LogError("key must be exactly 16 characters\n")
		return
	}

	if *Encrypt && len(*text) != 8 {
		sdes.LogError("encryption: text must be exactly 8 characters\n")
		return
	}

	if *Decrypt && len(*text) != 32 {
		sdes.LogError("decryption: text must be exactly 32 HEX characters\n")
		return
	}

	var textBS, keyBS, Temp *sdes.Bitset

	LeftChan := make(chan struct{})
	RightChan := make(chan struct{})

	LeftRound := make(chan int)
	RightRound := make(chan int)

	if *Encrypt {
		// Prepare the text
		textBS = sdes.CreateBitsetFromString(*text, false)
		// Adding timestamp
		textBS.Salt()
		textBS.Permute(&sdes.IP)
		// 64 bit halves - we need to split them once each.
		TL, TR := textBS.Split()

		// Left side halves
		TL0, TR0 := TL.Split()

		// Right side halves
		TL1, TR1 := TR.Split()

		// Prepare the key
		keyBS = sdes.CreateBitsetFromString(*key, false)
		keyBS.Permute(&sdes.PC1)

		// We generate the rounds based on the key
		sdes.GenerateRounds(keyBS)

		// We get the 56 bits sides
		KL, KR := keyBS.Split()

		// The 28 bit halves of the left side
		KL0, KR0 := KL.Split()
		// Compute the key rounds for encryption of the left side
		KeyRounds0 := sdes.PrecomputeRounds(KL0, KR0)

		// The 28 bit halves of the right side
		KL1, KR1 := KR.Split()
		// Compute the key rounds for encryption of the right side
		KeyRounds1 := sdes.PrecomputeRounds(KL1, KR1)

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
				Temp = sdes.XORBitsets(sdes.Feistel(KeyRounds0[RoundIndex], TR0), TL0)
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
				Temp = sdes.XORBitsets(sdes.Feistel(KeyRounds1[RoundIndex], TL1), TR1)
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
		FinalLeft := sdes.ConcatBitsets(TR0, TL0)
		FinalRight := sdes.ConcatBitsets(TR1, TL1)

		Final := sdes.ConcatBitsets(FinalLeft, FinalRight)
		Final.Permute(sdes.IP.Inverse())

		fmt.Printf("Output: %s\n", Final.ToHexString())
	} else {
		// Prepare the text
		textBS = sdes.CreateBitsetFromString(*text, true)
		textBS.Permute(&sdes.IP)
		// 64 bit halves - we need to split them once each.
		TL, TR := textBS.Split()

		// Left side halves
		TL0, TR0 := TL.Split()

		// Right side halves
		TL1, TR1 := TR.Split()

		// Prepare the key
		keyBS = sdes.CreateBitsetFromString(*key, false)
		keyBS.Permute(&sdes.PC1)

		// We generate the rounds based on the key
		sdes.GenerateRounds(keyBS)

		// We get the 56 bits sides
		KL, KR := keyBS.Split()

		// The 28 bit halves of the left side
		KL0, KR0 := KL.Split()
		// Compute the key rounds for encryption of the left side
		KeyRounds0 := sdes.PrecomputeRounds(KL0, KR0)

		// The 28 bit halves of the right side
		KL1, KR1 := KR.Split()
		// Compute the key rounds for encryption of the right side
		KeyRounds1 := sdes.PrecomputeRounds(KL1, KR1)

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
				Temp = sdes.XORBitsets(sdes.Feistel(KeyRounds0[RoundIndex], TR0), TL0)
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
				Temp = sdes.XORBitsets(sdes.Feistel(KeyRounds1[RoundIndex], TL1), TR1)
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
		FinalLeft := sdes.ConcatBitsets(TR0, TL0)
		FinalRight := sdes.ConcatBitsets(TR1, TL1)

		Final := sdes.ConcatBitsets(FinalLeft, FinalRight)
		Final.Permute(sdes.IP.Inverse())
		Final.RemoveSalt()

		fmt.Printf("Output: %s\n", Final.ToString())
	}
}
