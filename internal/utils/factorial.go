package utils

func FactorialNumber(n uint, ch chan uint) {
	if n == 0 {
		ch <- 1
		return
	}
	fact := n
	for i := n - 1; i > 0; i-- {
		fact *= i
	}
	ch <- fact
}
