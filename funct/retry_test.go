package funct

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	timeStamps := []time.Time{}
	var testFunc = func() error {
		timeStamps = append(timeStamps, time.Now())
		return errors.New("fail")
	}
	conf := &RetryConf{
		Retries:              3,
		ConcurrentErrHandler: false,
		ErrHandler:           func(_ error) {},
		WaitMap:              func(_ uint) float32 { return 0.1 },
	}
	err := Retry(conf, testFunc)
	if err == nil {
		t.Fatalf("found nil, expected error")
	}
	// Check retry count
	if len(timeStamps) != int(conf.Retries) {
		t.Errorf("expected %d retries, found %d instead", conf.Retries, len(timeStamps))
	}
	for i := 1; i < len(timeStamps); i++ {
		trueDelta := time.Duration(conf.WaitMap(uint(i+1)) * 1000000000)
		trueDelta = trueDelta.Truncate(10 * time.Millisecond)
		testDelta := timeStamps[i].Sub(timeStamps[i-1])
		testDelta = testDelta.Truncate(10 * time.Millisecond)
		if testDelta != trueDelta {
			t.Fatalf("testDelta: %v\ttestDelta: %v", testDelta, testDelta)
		}
	}
}
