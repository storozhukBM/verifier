package verifier_test

import (
	"github.com/storozhukBM/verifier"
	"fmt"
)

type Person struct {
	name       string
	age        int32
	hasLicense bool
}

func Example() {
	person := &Person{
		name:       "John Smith",
		age:        42,
		hasLicense: false,
	}
	err := sellAlcohol(person)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
	// customer should have license
}

func sellAlcohol(p *Person) error {
	verify := verifier.New()
	verify.That(p != nil, "person can't be nil")
	verify.PanicOnError() // let's imagine that we don't want to tolerate such calls in our system
	verify.That(p.age >= 21, "customer age should be 21 or higher, but yours: %d", p.age)
	verify.That(p.hasLicense, "customer should have license")
	if verify.GetError() != nil {
		return verify.GetError()
	}

	fmt.Print("yes, you can have some alcohol")
	return nil
}
