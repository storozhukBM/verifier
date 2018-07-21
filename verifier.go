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

func SetUnhandledVerificationsWriter(w io.Writer) {
	newWriter := unsafe.Pointer(&w)
	atomic.StorePointer(&verificationsWriter, newWriter)
}

type Verification struct {
	creationStack []uintptr
	err           error
	checked       bool
}

func New() *Verification {
	v := &Verification{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, printWarningOnUncheckedVerification)
	return v
}

func Offensive() *Verification {
	v := &Verification{
		creationStack: captureCreationStack(),
		checked:       false,
	}
	runtime.SetFinalizer(v, failProcessOnUncheckedVerification)
	return v
}

func Silent() *Verification {
	v := &Verification{
		checked: false,
	}
	return v
}

func (v *Verification) String() string {
	if v.err == nil {
		return "Verification success"
	}
	return "Verification failure: " + v.err.Error()
}

func (v *Verification) GetError() error {
	v.checked = true
	return v.err
}

func (v *Verification) PanicOnError() {
	v.checked = true
	if v.err != nil {
		panic("Verification failure: " + v.err.Error())
	}
}

func (v *Verification) Predicate(predicate func() bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if predicate() {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

func (v *Verification) That(positiveCondition bool, message string, args ...interface{}) {
	v.checked = false
	if v.err != nil {
		return
	}
	if positiveCondition {
		return
	}
	v.err = fmt.Errorf(message, args...)
}

func (v *Verification) Nil(ref interface{}, message string, args ...interface{}) {
	v.That(ref == nil, message, args...)
}

func (v *Verification) NotNil(ref interface{}, message string, args ...interface{}) {
	v.That(ref != nil, message, args...)
}

func (v *Verification) printCreationStack(writer io.Writer) {
	frames := runtime.CallersFrames(v.creationStack)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(writer, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func failProcessOnUncheckedVerification(v *Verification) {
	printWarningOnUncheckedVerification(v)
	os.Exit(1)
}

func printWarningOnUncheckedVerification(v *Verification) {
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
