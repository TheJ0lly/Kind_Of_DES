class Bitset:
    def __init__(self, value: str):
        """
        Turns a string into a bitset.
        """
        self.bits = []
        
        for char in value:
            # We encode the character, then we remove the '0b', then we make sure that the length of the bitstring is 8,
            # then we cast it to a list, then we add the contents inside the self.bits. 
            self.bits.extend(list(bin(char.encode()[0])[2:].zfill(8)))
    
    def len_bytes(self) -> int:
        """
        Returns the number of bytes.
        """

        return int(len(self.bits) / 8)

    def adjust_len(self, value: int):
        """
        If the length of the bytes is not a multiple of `value`, we add extra 0 bytes for padding.
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
            

