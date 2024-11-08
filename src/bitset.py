class Bitset:
    def __init__(self, value: str):
        """
        Turns a string into a bitset.
        """
        self.bits = []
        self.left: list[int]
        self.right: list[int]
        
        for char in value:
            # We encode the character, then we remove the '0b', then we make sure that the length of the bitstring is 8,
            # then we cast it to a list, then we add the contents inside the self.bits. 
            self.bits.extend(list(bin(char.encode()[0])[2:].zfill(8)))


    def len_bytes(self) -> int:
        """
        Returns the number of bytes.
        """
        return int(len(self.bits) / 8)
    
    def split(self):
        """
        Splits the key into 2 equal halves.
        """
        halfbitsindex = int(len(self.bits) / 2)

        self.left = [bit for bit in self.bits[:halfbitsindex]]
        self.right = [bit for bit in self.bits[halfbitsindex:]]


    def to_string(self) -> str:
        """
        Returns the string made by the bitset.
        """
        s = ""
        i = 0

        while i < len(self.bits):
            s += chr(int(f"0b{"".join(self.bits[i:i+8])}", 2))
            i += 8
        
        return s


class KeyBitset(Bitset):
    def __init__(self, value: str):
        super().__init__(value)


    def remove_parity_bits(self):
        """
        Alters the bits of the bitset, and reduces it from 64 bits to 56 bits.
        """
        i = 56
        while i >= 0:
            self.bits.pop(i)
            i -= 8



class TextBitset(Bitset):
    def __init__(self, value: str):
        """
        The bitset of the text
        """
        super().__init__(value)


    def adjust_len(self):
        """
        If the length of the bytes is not a multiple of 8, we add extra 0 bytes for padding.
        """

        target = 8
        # We divide by 8 to get the bytes
        currsize = int(len(self.bits) / 8)

        while currsize > target:
            target += 8
        
        # The amount of 0 bytes to add
        padding = target - currsize

        for i in range(padding):
            self.bits.extend(list(bin(0)[2:].zfill(8)))
