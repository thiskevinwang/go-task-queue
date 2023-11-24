package heartbeat

import (
	"time"
)

// StartHeartbeat starts the heartbeat functionality
func StartHeartbeat(onHeartbeat func(msg string)) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			onHeartbeat("tick")
		}
	}()
}
