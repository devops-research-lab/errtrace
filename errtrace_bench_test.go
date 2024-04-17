package errtrace

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func g1() error {
	return Errorf("g1: %w", g2())
}

func g2() error {
	return Errorf("g2: %w", g3())
}

func g3() error {
	return Errorf("g3: %w", g4())
}

func g4() error {
	return Errorf("g4: %w", g5())
}

func g5() error {
	return Errorf("g5: %w", g6())
}

func g6() error {
	return Errorf("g6 err")
}

func stacksBatch(err error) []string {
	pcs := make([]uintptr, 0, 10)
	for {
		if et, ok := err.(*errTrace); ok {
			pcs = append(pcs, et.pc)
			err = et.err
			continue
		}

		if e := errors.Unwrap(err); e != nil {
			err = e
			continue
		}

		break
	}

	funcs := make([]string, 0, 10)
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		funcs = append(funcs, f.Function)
		if !more {
			break
		}
	}

	return funcs
}

func stacksUnwrap(err error) []string {
	funcs := make([]string, 0, 10)
	for {
		if frame, underlying, ok := UnwrapFrame(err); ok {
			err = underlying
			funcs = append(funcs, frame.Func)
			continue
		}

		if e := errors.Unwrap(err); e != nil {
			err = e
			continue
		}

		break
	}

	return funcs
}

func BenchmarkUnwrapBatch(b *testing.B) {
	var got []string
	err := g1()
	for i := 0; i < b.N; i++ {
		got = stacksBatch(err)
	}
	if b.N == 1 {
		fmt.Println(got)
	}
}

func BenchmarkUnwrapFunc(b *testing.B) {
	var got []string
	err := g1()
	for i := 0; i < b.N; i++ {
		got = stacksUnwrap(err)
	}
	if b.N == 1 {
		fmt.Println(got)
	}
}
