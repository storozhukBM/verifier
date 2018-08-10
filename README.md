# Verifier "Defend against the impossible, because the impossible will happen." [![Build Status](https://travis-ci.org/storozhukBM/verifier.svg?branch=master)](https://travis-ci.org/storozhukBM/verifier)  [![Go Report Card](https://goreportcard.com/badge/github.com/storozhukBM/verifier)](https://goreportcard.com/report/github.com/storozhukBM/verifier) [![Coverage Status](https://coveralls.io/repos/github/storozhukBM/verifier/badge.svg?branch=master)](https://coveralls.io/github/storozhukBM/verifier?branch=master&k=1) [![GoDoc](https://godoc.org/github.com/storozhukBM/verifier?status.svg)](http://godoc.org/github.com/storozhukBM/verifier)

Package `verifier` provides simple [defensive programing](https://en.wikipedia.org/wiki/Defensive_programming) primitives.

Some software has higher than usual requirements for availability, safety or security.
Often in such projects people practice pragmatic paranoia with specific set of rules.
For example: each public function on any level of your application should check all arguments passed to it. Obviously checking for nil pointers, but also states and conditions, that it will rely on.

When you use such approaches your code can become a mess. And sometimes it's hard to distinguish between verification code and underlying business logic.

This small library is built on error handling pattern described by Rob Pike in Go blog called 
[Errors are values](https://blog.golang.org/errors-are-values).

It helps you quickly transform code from this
```go
if transfer == nil {
	return nil, errors.New("transfer can't be nil")
}
if person == nil {
	return nil, errors.New("person can't be nil")
}
if transfer.Destination == "" {
	return nil, errors.New("transfer destination can't be empty")
}
if transfer.Amount <= 0 {
	return nil, errors.New("transfer amount should be greater than zero")
}
if person.Name == "" {
	return nil, errors.New("name can't be empty")
}
if person.Age < 21 {
	return nil, fmt.Errorf("age should be 21 or higher, but yours: %d", person.Age)
}
if !person.HasLicense {
	return nil, errors.New("customer should have license")
}
```
to this
```go
verify := verifier.New()
verify.That(transfer != nil, "transfer can't be nil")
verify.That(person != nil, "person can't be nil")
if verify.GetError() != nil {
	return nil, verify.GetError()
}
verify.That(transfer.Destination != "", "transfer destination can't be empty")
verify.That(transfer.Amount > 0, "transfer amount should be greater than zero")
verify.That(person.Name != "", "name can't be empty")
verify.That(person.Age >= 21, "age should be 21 or higher, but yours: %d", person.Age)
verify.That(person.HasLicense, "customer should have license")
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

If you don't want/need such tracking use zero verifier `verifier.Verify{}`.
You can also redirect output from this package using `verifier.SetUnhandledVerificationsWriter(io.Writer)` method.

---
##### There is other libraries that can be useful when you employ defensive programming style
* [vala](https://github.com/kat-co/vala)
* [govalidate](https://github.com/tonyhb/govalidate)

---
##### License
Copyright 2018 Bohdan Storozhuk

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
