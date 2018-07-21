package verifier

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var verificationsWriter = unsafe.Pointer(&os.Stdout)

// New creates verification instance (recommended).
// It tracks verification state.
// If you forget to check internal error, using `GetError` or `PanicOnError` methods,
// it will write error message to UnhandledVerificationsWriter (default: os.Stdout).
// This mechanism will help you track down possible unhandled verifications.
// If you don't wan't to track anything, create verifier with Silent() function.
func New() *Verify {
	v := &Verify{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, printWarningOnUncheckedVerification)
	return v
}

// Creates verification instance without any tracking features.
// It's silent about unhandled verifications.
func Silent() *Verify {
	v := &Verify{
		checked: false,
	}
	return v
}

// Offensive creates verification instance (not-recommended).
// It tracks verification state and stops application process when founds unchecked verification.
// If you forget to check internal error, using `GetError` or `PanicOnError` methods,
// it will write error message to UnhandledVerificationsWriter (default: os.Stdout) and WILL STOP YOUR PROCESS.
// Created for people who adopt offensive programming(https://en.wikipedia.org/wiki/Offensive_programming).
// This mechanism will help you track down possible unhandled verifications.
// USE IT WISELY.
func Offensive() *Verify {
	v := &Verify{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, failProcessOnUncheckedVerification)
	return v
}

// Verify represents verification instance.
// All checks can be performed on it using `That` or `Predicate` functions.
// After one failed check all others won't count and predicates won't be evaluated.
// Use Verify.GetError function to check if there where any during verification process.
type Verify struct {
	creationStack []uintptr
	err           error
	checked       bool
}

// That verifies condition passed as first argument.
// If `positiveCondition == true`, verification will proceed for other checks.
// If `positiveCondition == false`, internal state will be filled with error,
// using message argument as format in fmt.Errorf(message, args...).
// After the first failed verification all others won't count and predicates won't be evaluated.
func (v *Verify) That(positiveCondition bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if positiveCondition {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

// That evaluates predicate passed as first argument.
// If `predicate() == true`, verification will proceed for other checks.
// If `predicate() == false`, internal state will be filled with error,
// using message argument as format in fmt.Errorf(message, args...).
// After the first failed verification all others won't count and predicates won't be evaluated.
func (v *Verify) Predicate(predicate func() bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if predicate() {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

// GetError extracts error from internal state to check if there where any during verification process.
func (v *Verify) GetError() error {
	v.checked = true
	return v.err
}

// PanicOnError panics if there is an error in internal state.
// Created for people who adopt offensive programming(https://en.wikipedia.org/wiki/Offensive_programming).
func (v *Verify) PanicOnError() {
	v.checked = true
	if v.err != nil {
		panic("verification failure: " + v.err.Error())
	}
}

// String represents verification and it's status as string type.
func (v *Verify) String() string {
	if v.err == nil {
		return "verification success"
	}
	return "verification failure: " + v.err.Error()
}

func (v *Verify) printCreationStack(writer io.Writer) {
	frames := runtime.CallersFrames(v.creationStack)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(writer, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func failProcessOnUncheckedVerification(v *Verify) {
	if v.checked {
		return
	}
	printWarningOnUncheckedVerification(v)
	os.Exit(1)
}

func printWarningOnUncheckedVerification(v *Verify) {
	if v.checked {
		return
	}
	writer := *(*io.Writer)(atomic.LoadPointer(&verificationsWriter))
	fmt.Fprintf(writer, "[ERROR] found verification with unhandled error: %s\n", v.err.Error())
	fmt.Fprint(writer, "verification was created here:\n")
	v.printCreationStack(writer)
}

func captureCreationStack() []uintptr {
	var rawStack [32]uintptr
	numberOfFrames := runtime.Callers(3, rawStack[:])
	return rawStack[:numberOfFrames]
}

// SetUnhandledVerificationsWriter gives you ability to override UnhandledVerificationsWriter (default: os.Stdout).
func SetUnhandledVerificationsWriter(w io.Writer) {
	newWriter := unsafe.Pointer(&w)
	atomic.StorePointer(&verificationsWriter, newWriter)
}
