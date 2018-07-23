package verifier_test

import (
	"bytes"
	"fmt"
	"github.com/storozhukBM/verifier"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestVerifier_positive_conditions(t *testing.T) {
	verify := verifier.Offensive()
	verify.That(rand.Float32() >= 0.0, "random should be positive")
	verify.That(rand.Float32() < 1.0, "random should less then 1.0")
	verify.That(true, "some other check with format %s; %d", "testCheck", 35)
	if verify.GetError() != nil {
		t.Error("verifier should be empty")
	}
	verify.That(rand.Float32() < 0.0, "expect error %s", "here")
	verify.That(rand.Float32() == 0.0, "should not have any deference")
	if verify.GetError() == nil {
		t.Fatal("verifier should be filled")
	}
	if verify.GetError().Error() != "expect error here" {
		t.Errorf("unexpected error message: %s", verify.GetError())
	}
	if fmt.Sprintf("%s", verify) != "verification failure: expect error here" {
		t.Errorf("unexpected verifier string representation: %s", verify)
	}
}

func TestVerifier_positive_not_evaluate_after_failure(t *testing.T) {
	counter := 0
	verify := verifier.New()

	verify.Predicate(func() bool {
		counter++
		return true
	}, "should be ok")
	verify.Predicate(func() bool {
		counter++
		return true
	}, "still OK")
	verify.Predicate(func() bool {
		counter++
		return false
	}, "should break here")
	verify.Predicate(func() bool {
		counter++
		return true
	}, "won't evaluate")

	if verify.GetError() == nil {
		t.Fatal("verifier should be filled")
	}
	if verify.GetError().Error() != "should break here" {
		t.Errorf("unexpected error message: %s", verify.GetError())
	}
	if counter != 3 {
		t.Errorf("unexpected evaluations happened")
	}
}

func TestVerifier_positive_panic_on_error(t *testing.T) {
	verify := verifier.New()
	verify.That(len("") != 0, "empty string is not nil")
	defer func() {
		panicObj := recover()
		if panicObj == nil {
			t.Fatal("verifier should have panic")
		}
		if panicObj != "verification failure: empty string is not nil" {
			t.Errorf("unexpected error message: %s", panicObj)
		}

	}()

	verify.PanicOnError()
}

func TestVerifier_negative_unhandled_error(t *testing.T) {
	localBuffer := &safeBuffer{}
	verifier.SetUnhandledVerificationsWriter(localBuffer)
	defer verifier.SetUnhandledVerificationsWriter(os.Stdout)

	verify := verifier.New()
	verify.That(len("") != 0, "empty string is not nil")
	runtime.GC()
	time.Sleep(10 * time.Millisecond)

	resultBuffer := localBuffer.String()
	if len(resultBuffer) == 0 {
		t.Fatal("unhandled error not found")
	}
	if !strings.HasPrefix(resultBuffer, "[ERROR] found unhandled verification: verification failure: empty string is not nil") {
		t.Fatalf("unexpected verifier buffer: %s", resultBuffer)
	}
}

func TestVerifier_negative_unhandled_success(t *testing.T) {
	localBuffer := &safeBuffer{}
	verifier.SetUnhandledVerificationsWriter(localBuffer)
	defer verifier.SetUnhandledVerificationsWriter(os.Stdout)

	verify := verifier.New()
	verify.That(true, "empty string is not nil")
	runtime.GC()
	time.Sleep(10 * time.Millisecond)

	resultBuffer := localBuffer.String()
	if len(resultBuffer) == 0 {
		t.Fatal("unhandled error not found")
	}
	if !strings.HasPrefix(resultBuffer, "[ERROR] found unhandled verification: verification success") {
		t.Fatalf("unexpected verifier buffer: %s", resultBuffer)
	}
}

func TestVerifier_negative_silent(t *testing.T) {
	localBuffer := &safeBuffer{}
	verifier.SetUnhandledVerificationsWriter(localBuffer)
	defer verifier.SetUnhandledVerificationsWriter(os.Stdout)

	verify := verifier.Verify{}
	verify.That(len("") != 0, "empty string is not nil")
	runtime.GC()
	time.Sleep(10 * time.Millisecond)

	resultBuffer := localBuffer.String()
	if len(resultBuffer) != 0 {
		t.Fatalf("unhandled printed something: %s", resultBuffer)
	}
}

func TestVerifier_negative_nil_error(t *testing.T) {
	var verify *verifier.Verify
	verify.That(len("") != 0, "empty string is not nil")
	if verify.GetError() == nil {
		t.Fatal("verifier should be filled")
	}
	if verify.GetError().Error() != "verifier instance is nil" {
		t.Errorf("unexpected error message: %s", verify.GetError())
	}
	if fmt.Sprintf("%s", verify) != "nil" {
		t.Errorf("unexpected verifier string representation: %s", verify)
	}
}

func TestVerifier_negative_nil_panic(t *testing.T) {
	var verify *verifier.Verify
	counter := 0
	verify.Predicate(func() bool {
		counter++
		return false
	}, "should not evaluate")
	defer func() {
		panicObj := recover()
		if panicObj == nil {
			t.Fatal("verifier should have panic")
		}
		if panicObj != "verifier instance is nil" {
			t.Errorf("unexpected error message: %s", panicObj)
		}

	}()

	verify.PanicOnError()
}

type safeBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (s *safeBuffer) Write(p []byte) (n int, err error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.b.Write(p)
}

func (s *safeBuffer) String() string {
	s.m.Lock()
	defer s.m.Unlock()
	return s.b.String()
}
