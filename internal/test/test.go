package test

type Test struct {
	// A user-friendly name of the test
	Name string
	// The reference to the day function to test
	DayFunc  func([]string, ...any) any
	Expected any
	// Extra arguments for example tests with reduced problem space
	Extras []any
	Input  []string
}
