package main

import (
	"flag"
)

func main() {
	key := flag.String("k", "", "the key used for encryption/decryption")
	text := flag.String("t", "", "the text to be encrypted/decrypted")
	decrypt := flag.Bool("d", false, "starts the decryption process")
	encrypt := flag.Bool("e", false, "starts the encryption process")
	help := flag.Bool("h", false, "help menu")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *decrypt && *encrypt {
		logError("both the encryption and decryption flags have been toggled - must have only 1\n")
		return
	}

	if !*decrypt && !*encrypt {
		logError("error: neither encryption nor decryption flag has been toggled - must have 1\n")
		return
	}

	if len(*key) != 16 {
		logError("key must be exactly 16 characters\n")
		return
	}

	if len(*text) != 8 {
		logError("text must be exactly 8 characters\n")
		return
	}

	keyBS := CreateBitset(*key, false)
	// rc := CreateRoundComputer(keyBS)
	keyBS.RemoveParityBits()

	textBS := CreateBitset(*text, true)
	textBS.Salt()
	textBS.Permute(&IP)
}
