package goft

import (
	"log"
	"time"
)

type FailableNullaryFunc func() error
type FloatSequence func(uint) float32
type ErrorHandler func(error)

var DefaultBackoffSequence FloatSequence = func(a uint) float32 {
	if a < 4 {
		return 1.0
	}
	return 1.5*float32(a) - 4.0
}

var DefaultErrorHandler ErrorHandler = func(err error) {
	log.Println(err)
}

// Retry keeps calling `f()` until it returns a nil error or number of attempts exceeds `attemmpts`.
// If after each call to `f()`, there is an error, then its passed to errorHandler.
// The `errorHandler` function will be called as a go routine each time.
// The parameter `waitMap` should map the current attempt number to number of seconds to wait before next attempt.
// The first time `waitMap` is called will be during attempt 2.
func (f FailableNullaryFunc) Retry(attempts uint, waitMap FloatSequence, errorHandler ErrorHandler) error {
	// Initiliaze attemp counter and error
	err := f.Call(errorHandler, true)
	var attempt uint = 2
	// Loop until success or attempts run out
	for err != nil {
		// Make sure we can still make an attempt
		if attempt >= attempts {
			// If out of attempts, return last attempts error
			return err
		}
		// Wait based on func w
		waitTime := time.Duration(waitMap(attempt)) * time.Second
		time.Sleep(waitTime)
		// Call function
		err = f.Call(errorHandler, true)
		// Increase attempt count
		attempt = attempt + 1
	}
	return nil
}

// Calls the function `f()` and passes the error, if there is one, to `errorHandler()`; also returns the error.
// The `errorHandler` function may be called as a go routine if `concurrent` is true.
func (f FailableNullaryFunc) Call(errorHandler ErrorHandler, concurrent bool) error {
	if err := f(); err != nil {
		if concurrent {
			go errorHandler(err)
		} else {
			errorHandler(err)
		}
		return err
	}
	return nil
}
