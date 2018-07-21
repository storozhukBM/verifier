# Verifier "Defend against the impossible, because the impossible will happen." [![Build Status](https://travis-ci.org/storozhukBM/verifier.svg?branch=master)](https://travis-ci.org/storozhukBM/verifier)  [![Go Report Card](https://goreportcard.com/badge/github.com/storozhukBM/verifier)](https://goreportcard.com/report/github.com/storozhukBM/verifier) [![Coverage Status](https://coveralls.io/repos/github/storozhukBM/verifier/badge.svg?branch=master)](https://coveralls.io/github/storozhukBM/verifier?branch=master) [![GoDoc](https://godoc.org/github.com/storozhukBM/verifier?status.svg)](http://godoc.org/github.com/storozhukBM/verifier)

Package `verifier` provides simple [defensive programing](https://en.wikipedia.org/wiki/Defensive_programming) primitives.

Some software have higher than usual requirements for availability, safety or security.
Very often in such projects people practice pragmatic paranoia with specific set of rules.
For example: each public function on any level of your application should check all arguments passed to it.
It obviously includes checking for nil pointer but also all sub-fields, states, conditions 
that your code will use, expect or rely on.

When you use such approaches your code can quickly become a mess and it becomes hard to distinguish between 
verification code and underlying business logic.

This small library is built on error handling pattern described by Rob Pike in Go blog called 
[Errors are values](https://blog.golang.org/errors-are-values).

It helps you quickly transform code from this
```go
    if person == nil {
    	panic("person can't be nil") // this state is impossible, if happens we should fail fast.
    }
    if person.name == "" {
        return nil, errors.New("name can't be empty")
    }
    if person.age < 21 {
    	return nil, fmt.Errorf("age should be 21 or higher, but yours: %d", p.age)
    }
    if !person.hasLicense {
    	return nil, errors.New("customer should have license")
    }
```
to this
```go
	verify := verifier.New()
	verify.That(person != nil, "person can't be nil")
	verify.PanicOnError() // use if you don't want to tolerate such errors
	
	verify.That(person.name != "", "name can't be empty")
	verify.That(person.age >= 21, "age should be 21 or higher, but yours: %d", p.age)
	verify.That(person.hasLicense, "customer should have license")
	if verify.GetError() != nil {
		return verify.GetError()
	}
```


##### License
Copyright 2018 Bohdan Storozhuk

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
