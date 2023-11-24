package utils

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func BoolToInt(x bool) int {
	if x {
		return 1
	}
	return 0
}

func SumOfBools(x ...bool) int {
	sum := 0
	for _, v := range x {
		sum += BoolToInt(v)
	}
	return sum
}
