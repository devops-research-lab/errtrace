package errtrace_test

import (
	"fmt"

	"braces.dev/errtrace"
)

func g1() error {
	return errtrace.Errorf("g2: %w", g2())
}

func g2() error {
	return errtrace.Errorf("g3: %w", g3())
}

func g3() error {
	return errtrace.New("failed")
}

func Example_unwrap() {
	err := g1()
	for {
		frame, inner, ok := errtrace.UnwrapFrame(err)
		if !ok {
			break
		}
		err = inner

		fmt.Println(frame.Func)
	}

	// Output:
	//failed
	//
	//braces.dev/errtrace_test.f3
	//	/path/to/errtrace/example_trace_test.go:3
	//braces.dev/errtrace_test.f2
	//	/path/to/errtrace/example_trace_test.go:2
	//braces.dev/errtrace_test.f1
	//	/path/to/errtrace/example_trace_test.go:1
}
