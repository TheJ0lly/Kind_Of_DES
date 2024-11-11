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

var Running = true

func WatchReady(left chan bool, right chan bool, done chan bool) {
	leftReady := false
	rightReady := false
	for RoundIndex < 16 {
		select {
		case <-left:
			leftReady = true
			logInfo("left side READY")
		case <-right:
			rightReady = true
			logInfo("right side READY")
		}

		if leftReady && rightReady {
			logInfo("proceeding to next round")
			left <- true
			right <- true
			leftReady = false
			rightReady = false

			// We advance the round index
			RoundIndex++
		}
	}

	// We wait for the final signals
	for !leftReady && !rightReady {
		select {
		case <-left:
			leftReady = true
		case <-right:
			rightReady = true
		}
	}

	done <- true
}

// We get the 32 bits text halves
func EncryptPlainTextSide(textSide **Bitset, KeyHalf **Bitset) {
	var lt, rt *Bitset

	// for RoundIndex < 16 {
	lt, rt = (*textSide).Split()
	lk, rk := (*KeyHalf).Split()

	// crt is a copy of the right text, that will be used for the feistel function.
	crt := Copy(rt)

	// We expand the "right" text side.
	crt.Permute(&Expansion)

	// We shift the sides by RC.Rounds[RoundIndex] bits
	lk.ShiftBy(RC.Rounds[RoundIndex], true)
	rk.ShiftBy(RC.Rounds[RoundIndex], true)

	// We concatenate the key sides
	ck := ConcatBitsets(lk, rk)

	// We permute with PC2
	ck.Permute(&PC2)

	// We get the 48 bit XOR'd "right" side
	sBS := XORBitsets(ck, crt)

	// We apply and then concatenate the SBoxes output
	sBSOut := int32(0)
	for i := 0; i < 8; i++ {
		sBSOut <<= 4
		sBSOut |= sBS.ApplySBox(i)
	}

	// We create the bitset from the 32 bit SBoxes output
	crt = CreateBitsetFromInt32(sBSOut, false)

	// We apply the P permutation
	crt.Permute(&PPermutation)

	// We XOR the left and right side and switch text sides
	rt, lt = XORBitsets(lt, crt), rt

	// We concatenate the text sides and then the new text is the text halves combined.
	(*textSide) = ConcatBitsets(lt, rt)

	// We combine the key sides and then assign it to the KeyHalf.
	(*KeyHalf) = ConcatBitsets(lk, rk)
}

func DecryptTextSide(textSide **Bitset, KeyHalf **Bitset) {
	// We get the 32 bits text halves
	var lt, rt *Bitset

	// for RoundIndex < 16 {
	lt, rt = (*textSide).Split()
	lk, rk := (*KeyHalf).Split()

	// crt is a copy of the right text, that will be used for the feistel function.
	crt := Copy(rt)

	// We expand the "right" text side.
	crt.Permute(&Expansion)

	// We skip the first round of rotations in decrypt
	if RoundIndex != 15 {
		// We shift the sides by RC.Rounds[RoundIndex] bits
		lk.ShiftBy(RC.Rounds[RoundIndex], false)
		rk.ShiftBy(RC.Rounds[RoundIndex], false)
	}

	// We concatenate the key sides
	ck := ConcatBitsets(lk, rk)
	// We permute with PC2
	ck.Permute(&PC2)

	// We get the 48 bit XOR'd "right" side
	sBS := XORBitsets(ck, crt)

	// We apply and then concatenate the SBoxes output
	sBSOut := int32(0)
	for i := 0; i < 8; i++ {
		sBSOut <<= 4
		sBSOut |= sBS.ApplySBox(i)
	}

	// We create the bitset from the 32 bit SBoxes output
	crt = CreateBitsetFromInt32(sBSOut, false)

	// We apply the P permutation
	crt.Permute(&PPermutation)

	// We XOR the left and right side and switch text sides
	rt, lt = XORBitsets(lt, crt), rt

	// We concatenate the text sides and then the new text is the text halves combined.
	(*textSide) = ConcatBitsets(lt, rt)

	// We combine the key sides and then assign it to the KeyHalf.
	(*KeyHalf) = ConcatBitsets(lk, rk)
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
		logError("text must be exactly 8 characters\n")
		return
	}

	keyBS := CreateBitsetFromString(*key, false, false)
	// RC = CreateRoundComputer(keyBS)

	// keyBS.RemoveParityBits()
	keyBS.Permute(&PC1)

	KEY_LEFT, KEY_RIGHT = keyBS.Split()

	var textBS *Bitset

	if *Decrypt {
		// If we decrypt the text, the input will be in hex format
		textBS = CreateBitsetFromString(*text, true, true)

		// If we decrypt we start with the INV_IP
		textBS.Permute(&IP)
		textBS.Print()
	} else {
		// If we encrypt the text, the input will be in plaintext format
		textBS = CreateBitsetFromString(*text, true, false)

		// If we encrypt we add the salt, and permute with IP
		textBS.Salt()
		textBS.Permute(&IP)
	}

	TEXT_LEFT, TEXT_RIGHT = textBS.Split()

	if *Encrypt {
		RoundIndex = 0
		for ; RoundIndex < 16; RoundIndex++ {
			// Left 64 bits of the text and the left 56 bits of the key
			EncryptPlainTextSide(&TEXT_LEFT, &KEY_LEFT)
			// Right 64 bits of the text and the right 56 bits of the key
			EncryptPlainTextSide(&TEXT_RIGHT, &KEY_RIGHT)
			KEY_LEFT, KEY_RIGHT = KEY_RIGHT, KEY_LEFT
		}
	} else {
		RoundIndex = 15
		for ; RoundIndex >= 0; RoundIndex-- {
			// Left 64 bits of the text and the left 56 bits of the key
			DecryptTextSide(&TEXT_LEFT, &KEY_LEFT)

			// Right 64 bits of the text and the right 56 bits of the key
			DecryptTextSide(&TEXT_RIGHT, &KEY_RIGHT)
			KEY_LEFT, KEY_RIGHT = KEY_RIGHT, KEY_LEFT

			fmt.Printf("Round %d: ", RoundIndex)
			ConcatBitsets(KEY_LEFT, KEY_RIGHT).Print()
		}
	}

	Final := ConcatBitsets(TEXT_LEFT, TEXT_RIGHT)

	// Final.Print()
	Final.Permute(IP.Inverse())

	if *Decrypt {
		Final.RemoveSalt()
	}

	if *Encrypt {
		KeyOutput := ConcatBitsets(KEY_LEFT, KEY_RIGHT)
		KeyOutput.Permute(PC1.Inverse())
		fmt.Printf("Output: %s| Length %d | Key %s\n", Final.ToHexString(), Final.Len(), KeyOutput.ToString())
	} else {
		KeyOutput := ConcatBitsets(KEY_LEFT, KEY_RIGHT)
		KeyOutput.Permute(PC1.Inverse())
		fmt.Printf("Output: %s| Length %d | Key %s\n", Final.ToString(), Final.Len(), KeyOutput.ToString())
	}
}
