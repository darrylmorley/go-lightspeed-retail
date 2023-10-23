package lightspeedsdk

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	maxRetries   = 3
	initialDelay = 1 * time.Second
	maxDelay     = 32 * time.Second
)

func (sdk *SDK) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	delay := initialDelay
	for i := 0; i < maxRetries; i++ {
		resp, err = sdk.Client.Do(req)
		if err == nil && resp.StatusCode < http.StatusInternalServerError {
			// Success, or a client error (4xx). Exit the loop.
			return resp, nil
		}

		// If this was the last attempt, return the last error.
		if i == maxRetries-1 {
			break
		}

		time.Sleep(addJitter(delay))
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return nil, err
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func addJitter(duration time.Duration) time.Duration {
	jitter := duration / 2
	return duration + time.Duration(r.Int63n(int64(jitter+1))-int64(jitter)/2)
}
