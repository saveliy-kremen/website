package utils

import "testing"

func TestFactorialNumber(t *testing.T) {
	testCases := []struct {
		testName string
		testData uint
		testRes  uint
	}{
		{"factorial_0", 0, 1},
		{"factorial_10", 10, 3628800},
		{"factorial_20", 20, 2432902008176640000},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			ch := make(chan uint)
			go FactorialNumber(testCase.testData, ch)
			fact := <-ch
			if fact != testCase.testRes {
				t.Errorf("Expected factorial of %d must be %d", testCase.testData, testCase.testRes)
			}
		})
	}
}
