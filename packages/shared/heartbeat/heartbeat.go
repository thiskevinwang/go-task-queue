package heartbeat

import (
	"time"
)

// StartHeartbeat starts the heartbeat functionality
func StartHeartbeat(every int64, onHeartbeat func(msg string)) {
	ticker := time.NewTicker(time.Duration(every) * time.Second)
	go func() {
		for range ticker.C {
			onHeartbeat("tick")
		}
	}()
}
