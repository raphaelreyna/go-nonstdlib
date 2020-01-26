package funct

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	timeStamps := []time.Time{}
	var testFunc FailableNullaryFunc = func() error {
		timeStamps = append(timeStamps, time.Now())
		return errors.New("fail")
	}
	conf := &RetryConf{
		Retries:              7,
		ConcurrentErrHandler: false,
	}
	err := testFunc.Retry(conf)
	if err == nil {
		t.Fatalf("found nil, expected error")
	}
	for i := 1; i < len(timeStamps); i++ {
		trueDelta := time.Duration(DefaultBackoffSequence(uint(i+1)) * 1000000000)
		trueDelta = trueDelta.Truncate(500 * time.Millisecond)
		testDelta := timeStamps[i].Sub(timeStamps[i-1])
		testDelta = testDelta.Truncate(500 * time.Millisecond)
		if testDelta != trueDelta {
			t.Fatalf("testDelta: %v\ttestDelta: %v", testDelta, testDelta)
		}
	}
}
