package funct

import (
	"log"
	"time"
)

type FailableNullaryFunc func() error
type FloatSequence func(uint) float32
type ErrorHandler func(error)

// DefaultBackoffSequence waits 1 second for first 4 seconds, and then increases wait time by 1.5 seconds per attempt
var DefaultBackoffSequence FloatSequence = func(a uint) float32 {
	if a < 4 {
		return 1.0
	}
	return 1.5*float32(a) - 4.0
}

// DefaultErrorHandler logs the error
var DefaultErrorHandler ErrorHandler = func(err error) {
	log.Println(err)
}

type RetryConf struct {
	Retries uint
	// WaitMap should map the current attempt number to number of seconds to wait before next attempt.
	// The first time WaitMap is called will be during attempt 2.
	WaitMap FloatSequence
	// ErrorHandler will be called as a go routine each time.
	ErrHandler           ErrorHandler
	ConcurrentErrHandler bool
}

// Retry keeps calling `f()` until it returns a nil error or number of retries is exceeded.
// If after each call to `f()`, there is an error, then its passed to an error handler.
func (f FailableNullaryFunc) Retry(conf *RetryConf) error {
	retries := conf.Retries
	errHandler := conf.ErrHandler
	if errHandler == nil {
		errHandler = DefaultErrorHandler
	}
	wm := conf.WaitMap
	if wm == nil {
		wm = DefaultBackoffSequence
	}
	concerr := conf.ConcurrentErrHandler
	// Initiliaze attemp counter and error
	err := f.Call(errHandler, concerr)
	var attempt uint = 2
	// Loop until success or attempts run out
	for err != nil {
		// Make sure we can still make an attempt
		if attempt >= conf.Retries {
			// If out of attempts, return last attempts error
			return err
		}
		// Wait based on func w
		waitTime := time.Duration(wm(retries)) * time.Second
		time.Sleep(waitTime)
		// Call function
		err = f.Call(errHandler, true)
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
