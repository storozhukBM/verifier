package verifier

import (
	"fmt"
	"runtime"
	"io"
	"os"
	"sync/atomic"
	"unsafe"
)

var verificationsWriter = unsafe.Pointer(&os.Stdout)

func SetUnhandledVerificationsWriter(w io.Writer) {
	newWriter := unsafe.Pointer(&w)
	atomic.StorePointer(&verificationsWriter, newWriter)
}

type verification struct {
	creationStack []uintptr
	err           error
	checked       bool
}

func New() *verification {
	v := &verification{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, printWarningOnUncheckedVerification)
	return v
}

func Offensive() *verification {
	v := &verification{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, failProcessOnUncheckedVerification)
	return v
}

func Silent() *verification {
	v := &verification{
		checked: false,
	}
	return v
}

func (v *verification) String() string {
	if v.err == nil {
		return "verification success"
	}
	return "verification failure: " + v.err.Error()
}

func (v *verification) GetError() error {
	v.checked = true
	return v.err
}

func (v *verification) PanicOnError() {
	v.checked = true
	if v.err != nil {
		panic("verification failure: " + v.err.Error())
	}
}

func (v *verification) Predicate(predicate func() bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if predicate() {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

func (v *verification) That(positiveCondition bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if positiveCondition {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

func (v *verification) Nil(ref interface{}, message string, args ...interface{}) {
	v.That(ref == nil, message, args...)
}

func (v *verification) NotNil(ref interface{}, message string, args ...interface{}) {
	v.That(ref != nil, message, args...)
}

func (v *verification) printCreationStack(writer io.Writer) {
	frames := runtime.CallersFrames(v.creationStack)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(writer, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func failProcessOnUncheckedVerification(v *verification) {
	printWarningOnUncheckedVerification(v)
	os.Exit(1)
}

func printWarningOnUncheckedVerification(v *verification) {
	if v.checked {
		return
	}
	writer := *(*io.Writer)(atomic.LoadPointer(&verificationsWriter))
	fmt.Fprintf(writer, "[ERROR] found verifier with unhandled error: %s\n", v.err.Error())
	fmt.Fprint(writer, "verifier was created here:\n")
	v.printCreationStack(writer)
}

func captureCreationStack() []uintptr {
	var rawStack [32]uintptr
	numberOfFrames := runtime.Callers(3, rawStack[:])
	return rawStack[:numberOfFrames]
}
