package verifier

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
)

// New creates verification instance (recommended).
// It tracks verification state.
// If you forget to check internal error, using `GetError` or `PanicOnError` methods,
// it will write error message to UnhandledVerificationsWriter (default: os.Stdout).
// This mechanism will help you track down possible unhandled verifications.
// If you don't wan't to track anything, create zero verifier `Verify{}`.
func New() *Verify {
	v := &Verify{
		creationStack: captureCreationStack(),
	}
	runtime.SetFinalizer(v, printWarningOnUncheckedVerification)
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

// WithError verifies condition passed as first argument.
// If `positiveCondition == true`, verification will proceed for other checks.
// If `positiveCondition == false`, internal state will be filled with error specified as second argument.
// After the first failed verification all others won't count and predicates won't be evaluated.
func (v *Verify) WithError(positiveCondition bool, err error) *Verify {
	vObj := v
	if v == nil {
		vObj = &Verify{}
	}

	vObj.checked = false
	if vObj.err != nil {
		return vObj
	}
	if positiveCondition {
		return vObj
	}
	vObj.err = err
	return vObj
}

// That verifies condition passed as first argument.
// If `positiveCondition == true`, verification will proceed for other checks.
// If `positiveCondition == false`, internal state will be filled with error,
// using message argument as format in fmt.Errorf(message, args...).
// After the first failed verification all others won't count and predicates won't be evaluated.
func (v *Verify) That(positiveCondition bool, message string, args ...interface{}) *Verify {
	vObj := v
	if v == nil {
		vObj = &Verify{}
	}

	vObj.checked = false
	if vObj.err != nil {
		return vObj
	}
	if positiveCondition {
		return vObj
	}
	vObj.err = fmt.Errorf(message, args...)
	return vObj
}

// That evaluates predicate passed as first argument.
// If `predicate() == true`, verification will proceed for other checks.
// If `predicate() == false`, internal state will be filled with error,
// using message argument as format in fmt.Errorf(message, args...).
// After the first failed verification all others won't count and predicates won't be evaluated.
func (v *Verify) Predicate(predicate func() bool, message string, args ...interface{}) *Verify {
	vObj := v
	if v == nil {
		vObj = &Verify{}
	}
	vObj.checked = false
	if vObj.err != nil {
		return vObj
	}
	if predicate() {
		return vObj
	}
	vObj.err = fmt.Errorf(message, args...)
	return vObj
}

// GetError extracts error from internal state to check if there where any during verification process.
func (v *Verify) GetError() error {
	if v == nil {
		return errors.New("verifier instance is nil")
	}
	v.checked = true
	return v.err
}

// PanicOnError panics if there is an error in internal state.
// Created for people who adopt offensive programming(https://en.wikipedia.org/wiki/Offensive_programming).
func (v *Verify) PanicOnError() {
	if v == nil {
		panic("verifier instance is nil")
	}
	v.checked = true
	if v.err != nil {
		panic("verification failure: " + v.err.Error())
	}
}

// String represents verification and it's status as string type.
func (v *Verify) String() string {
	if v == nil {
		return "nil"
	}
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

func captureCreationStack() []uintptr {
	var rawStack [32]uintptr
	numberOfFrames := runtime.Callers(3, rawStack[:])
	return rawStack[:numberOfFrames]
}

type writerWrapper struct {
	value io.Writer
}

var verificationsWriter atomic.Value

// SetUnhandledVerificationsWriter gives you ability to override UnhandledVerificationsWriter (default: os.Stdout).
func SetUnhandledVerificationsWriter(w io.Writer) {
	verificationsWriter.Store(writerWrapper{w})
}

func init() {
	SetUnhandledVerificationsWriter(os.Stdout)
}

func printWarningOnUncheckedVerification(v *Verify) {
	if v.checked {
		return
	}
	rawWriter := verificationsWriter.Load()
	if rawWriter == nil || rawWriter.(writerWrapper).value == nil {
		rawWriter = writerWrapper{os.Stdout}
	}
	writer := rawWriter.(writerWrapper).value
	fmt.Fprintf(writer, "[ERROR] found unhandled verification: %s\n", v)
	fmt.Fprint(writer, "verification was created here:\n")
	v.printCreationStack(writer)
}
