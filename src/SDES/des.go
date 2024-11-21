package sdes

import "slices"

var RC []int = make([]int, 16)

func GenerateRounds(key *Bitset) {
	for i := 0; i < key.Len()-7; i++ {
		part := key.Bits[i : i+7]
		val := byte(0)

		for j := 0; j < 7; j++ {
			val <<= 1
			val |= part[j]
		}

		RC[val%16]++
	}

	maxIndexes := make(map[int]int, 16)

	// We set the map to vector values and indexes.
	for i, v := range RC {
		maxIndexes[v] = i
	}

	keys := make([]int, 0, len(maxIndexes))
	for k := range maxIndexes {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for i := len(keys) - 1; i >= len(keys)-4; i-- {
		RC[maxIndexes[keys[i]]] = -1
	}

	// Everything that's not a 1, we set to 2.
	for i := 0; i < len(RC); i++ {
		if RC[i] != -1 {
			RC[i] = 2
		} else if RC[i] == -1 {
			RC[i] = 1
		}
	}
}

func LeftRotateKey(key *Bitset, round int) *Bitset {
	l, r := key.Split()
	l.ShiftBy(RC[round], true)
	r.ShiftBy(RC[round], true)

	rotKey := ConcatBitsets(l, r)

	return rotKey
}

func PrecomputeRounds(flk, frk *Bitset) []*Bitset {
	roundKeys := make([]*Bitset, 16)
	firstkey := ConcatBitsets(flk, frk)

	l, r := firstkey.Split()

	// We append the first key, which is the original key + leftshift[Round].
	roundKeys[0] = LeftRotateKey(ConcatBitsets(l, r), RC[0])

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
