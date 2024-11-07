IP = [
    58, 50, 42, 34, 26, 18, 10, 2,
    60, 52, 44, 36, 28, 20, 12, 4,
    62, 54, 46, 38, 30, 22, 14, 6,
    64, 56, 48, 40, 32, 24, 16, 8,
    57, 49, 41, 33, 25, 17, 9, 1,
    59, 51, 43, 35, 27, 19, 11, 3,
    61, 53, 45, 37, 29, 21, 13, 5,
    63, 55, 47, 39, 31, 23, 15, 7
    ]

class Permutation:
    def __init__(self, *args) -> None:
        """
        A permutation.
        """
        self.data = []

        for x in args:
            self.data.extend(x)

    # I don't know why it works in strings. Its called a "Forward Reference Type Hint".
    def compute_inverse(self) -> "Permutation":
        """
        It computes the inverse of the permutation and returns the new permutation.
        """
        inverse = [None] * len(self.data)

        for i in range(len(self.data)):
            inverse[self.data[i] - 1] = i+1


        return Permutation(inverse)

