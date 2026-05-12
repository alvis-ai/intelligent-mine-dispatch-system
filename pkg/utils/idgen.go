package utils

import (
	"sync"
	"time"
)

var (
	machineID uint8 = 1
	seq       uint32
	mu        sync.Mutex
	lastTS    int64
)

func SetMachineID(id uint8) {
	machineID = id
}

func NextID() uint64 {
	mu.Lock()
	defer mu.Unlock()

	ts := time.Now().UnixMilli()
	if ts == lastTS {
		seq++
	} else {
		seq = 0
		lastTS = ts
	}

	return uint64(ts)<<22 | uint64(machineID)<<12 | uint64(seq)
}
