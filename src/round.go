package main

type RoundComputer struct {
	Rounds []int
}

// ComputeByte gets the value of the byte and transforms it such that it maps to a specific index.
func ComputeByte(bits []byte) byte {
	sum := byte(0)

	for i := 0; i < len(bits); i++ {
		sum = (sum << 1) | bits[i]
	}

	return sum % 16
}

// CreateRoundComputer generates the set of rounds where the bits are shifted 1 or 2 bits.
func CreateRoundComputer(k *Bitset) *RoundComputer {
	rounds := make([]int, 16)

	maxRounds := 0

	for i := 0; i < k.Len(); i += 8 {
		index := ComputeByte(k.Bits[i : i+8])
		cur := &rounds[index]

		if *cur == 0 { // If current is 0, we increase it to 1, to signal that in the respective round, we shift by 1.
			(*cur)++
		} else if *cur == 1 && maxRounds < 4 { // Else we have found a 2-bit shift round, and we must have at most 4 2-bit rounds.
			(*cur)++
			maxRounds++
		}
	}

	// We make sure there is no round with 0 bit shift.
	for i := 0; i < 16; i++ {
		if rounds[i] == 0 {
			rounds[i] = 1
		}
	}

	// We invert the rounds to get the 4 x 1 and 12 x 2 rounds
	for i := 0; i < 16; i++ {
		if rounds[i] == 1 {
			rounds[i] = 2
		} else {
			rounds[i] = 1
		}
	}

	return &RoundComputer{Rounds: rounds}
}
