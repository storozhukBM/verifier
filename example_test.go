package verifier

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
	verify.NotNil(p, "person can't be nil")
	verify.PanicOnError()
	verify.That(p.age >= 21, "customer age should be 21 or higher, but yours: %d", p.age)
	verify.That(p.hasLicense, "customer should have license")
	if verify.GetError() != nil {
		return verify.GetError()
	}

	fmt.Println("yes, you can have some alcohol")
	return nil
}
