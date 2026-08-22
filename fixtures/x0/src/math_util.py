import math


def add(a: int, b: int) -> int:
    """Return the sum of two integers (fixture symbol for P0-X indexing)."""
    return a + b


def hypotenuse(a: float, b: float) -> float:
    """Return the Euclidean length of a right triangle's hypotenuse."""
    return math.sqrt(a * a + b * b)
