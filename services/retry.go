package services

import (
	"time"

	"github.com/IAmRiteshKoushik/devpool/cmd"
)

const (
	retryDelay = 5 * time.Second
)

func withRetry(fn func() error) {
	for {
		err := fn()
		if err == nil {
			return
		}
		cmd.Log.Error("Retrying after error", err)
		time.Sleep(retryDelay)
	}
}
