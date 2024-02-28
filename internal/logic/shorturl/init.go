package shorturl

import (
	"context"
	"time"
)

func init() {
	go func() {
		time.Sleep(60 * time.Second)
		ctx := context.Background()
		go TimerUpdateCount(ctx)
		go TimerCheckStatus(ctx)
	}()
}
