package main

import (
	"log"
	"sync"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/IAmRiteshKoushik/devpool/services"
)

func main() {
	config, err := cmd.NewAppConfig()
	if err != nil {
		log.Fatalf("Failed to load app config: %v", err)
	}
	cmd.App = config

	if err := cmd.InitValkey(); err != nil {
		log.Fatalf("Failed to initialize Valkey: %v", err)
	}

	githubService, err := services.NewGithubService(config)
	if err != nil {
		log.Fatalf("Failed to initialize GithubService: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		services.ConsumeIssueStream(githubService)
	}()

	go func() {
		defer wg.Done()
		services.ConsumeBountyStream(githubService)
	}()

	go func() {
		defer wg.Done()
		services.ConsumeSolutionStream(githubService)
	}()

	go func() {
		defer wg.Done()
		services.ConsumeAchievementStream(githubService)
	}()

	wg.Wait()
}
