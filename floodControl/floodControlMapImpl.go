package floodControl

import (
	"context"
	"sync"
	"time"
)

type MapImpl struct {
	mu         sync.Mutex
	requestMap map[int64][]time.Time
	N          time.Duration
	K          int
}

func NewFloodControlMapImpl(N time.Duration, K int) *MapImpl {
	return &MapImpl{
		requestMap: make(map[int64][]time.Time),
		N:          N,
		K:          K,
	}
}

func (fc *MapImpl) Check(ctx context.Context, userID int64) (bool, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	requests, ok := fc.requestMap[userID]
	currentTime := time.Now()

	if !ok {
		fc.requestMap[userID] = []time.Time{currentTime}
		return true, nil
	}

	var validRequests []time.Time
	for _, t := range requests {
		if t.After(currentTime.Add(-fc.N * time.Second)) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= fc.K {
		return false, nil
	}

	validRequests = append(validRequests, currentTime)
	fc.requestMap[userID] = validRequests

	return true, nil
}
