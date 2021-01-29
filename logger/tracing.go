package logger

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"
)

// getCallStack
func getCallStack() (stack []string) {
	const (
		maxStackDepth = 10000
		skipRuntime   = 2 // skip runtime and startup
		skipCallers   = 2
		thisPackage   = "logger"
	)

	callers := make([]uintptr, maxStackDepth) // min 1

	stackDepth := runtime.Callers(skipCallers, callers)
	frames := runtime.CallersFrames(callers)

	for i := 0; i < skipRuntime; i++ {
		frames.Next()
	}

	for i := 0; i < stackDepth-skipRuntime-skipCallers; i++ {
		frame, next := frames.Next()
		if !strings.Contains(frame.File, thisPackage) {
			stack = append(stack, fmt.Sprintf("[%s %d]", frame.File, frame.Line))
		}
		if !next {
			break
		}
	}

	reverseStack(stack)

	return stack
}

// in place
func reverseStack(stack []string) {
	middle := int(math.Floor(float64(len(stack)) / 2.0))

	lastIndex := len(stack) - 1
	for i := 0; i < middle; i++ {
		stack[i], stack[lastIndex] = stack[lastIndex], stack[i]
		lastIndex--
	}
}

func getTrace(format string, v ...interface{}) (t string) {
	stack := getCallStack()

	t = time.Now().Format(dateTimeFormat) + ":\n"
	for _, s := range stack {
		t += s + "\n"
	}
	t += fmt.Sprintf(format, v...) + "\n"

	return t
}
