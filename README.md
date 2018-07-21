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
	return nil, errors.New("person can't be nil")
}
if len(person.name) == "" {
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
verify.That(person.name != "", "name can't be empty")
verify.That(person.age >= 21, "age should be 21 or higher, but yours: %d", p.age)
verify.That(person.hasLicense, "customer should have license")
if verify.GetError() != nil {
	return nil, verify.GetError()
}
```

It also can help you to track down unchecked verifiers when you forget to do it 
```go
verify := verifier.New()
verify.That(user.HashPermission(READ, ACCOUNTS), "user has no permission to read accounts")
```

In this example we have forgot to check verifier using verify.GetError() and
after some time you'll see in your log such error:
```
[ERROR] found unhandled verification: verification failure: user has no permission to read accounts
  verification was created here:
  github.com/storozhukBM/verifier_test.TestVerifier
    /Users/bogdanstorozhuk/verification-lib/src/github.com/storozhukBM/verifier/verifier_test.go:91
  testing.tRunner
    /usr/local/Cellar/go/1.10.2/libexec/src/testing/testing.go:777
  runtime.goexit
    /usr/local/Cellar/go/1.10.2/libexec/src/runtime/asm_amd64.s:2361
```

If you don't want/need such tracking use `verifier.Silent()` function to create your verifiers.
You can also redirect output from this package using `verifier.SetUnhandledVerificationsWriter(io.Writer)` method.
---

##### License
Copyright 2018 Bohdan Storozhuk

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
