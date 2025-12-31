package main

import (
	"log"
	"sync"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/IAmRiteshKoushik/devpool/services"
)

func main() {
	if err := cmd.InitValkey(); err != nil {
		log.Fatalf("Failed to initialize Valkey: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		services.ConsumeIssueStream()
	}()

	go func() {
		defer wg.Done()
		services.ConsumeBountyStream()
	}()

	go func() {
		defer wg.Done()
		services.ConsumeSolutionStream()
	}()

	go func() {
		defer wg.Done()
		services.ConsumeAchievementStream()
	}()

	wg.Wait()
}