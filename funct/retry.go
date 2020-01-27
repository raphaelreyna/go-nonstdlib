package funct

import (
	"log"
	"time"
)

// DefaultBackoffSequence waits 1 second for first 4 seconds, and then increases wait time by 1.5 seconds per attempt
func DefaultBackoffSequence(a uint) float32 {
	if a < 4 {
		return 1.0
	}
	return 1.5*float32(a) - 3.5
}

// DefaultErrorHandler logs the error
func DefaultErrorHandler(err error) {
	log.Println(err)
}

// RetryConf holds the configuration for a function call retry
type RetryConf struct {
	Retries uint
	// WaitMap should map the current attempt number to number of seconds to wait before next attempt.
	// The first time WaitMap is called will be during attempt 2.
	WaitMap func(uint) float32
	// ErrorHandler will be called as a go routine each time.
	ErrHandler           func(error)
	ConcurrentErrHandler bool
}

// Retry keeps calling `f()` until it returns a nil error or number of retries is exceeded.
// If after each call to `f()`, there is an error, then its passed to an error handler.
func Retry(conf *RetryConf, f func() error) error {
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
	// Initiliaze attempt counter and error
	err := call(concerr, errHandler, f)
	var attempt uint = 2
	// Loop until success or attempts run out
	for err != nil {
		// Make sure we can still make an attempt
		if attempt > retries {
			// If out of attempts, return last attempts error
			return err
		}
		// Wait based on func w
		waitTime := time.Duration(wm(attempt) * 1000000000)
		time.Sleep(waitTime)
		// Call function
		err = call(concerr, errHandler, f)
		// Increase attempt count
		attempt = attempt + 1
	}
	return nil
}

func call(concurrent bool, errorHandler func(error), f func() error) error {
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
