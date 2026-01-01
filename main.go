package main

import (
	"log"
	"os"
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

	file, err := os.OpenFile("devpool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	defer file.Close()
	cmd.NewLogger(file)

	if err := cmd.InitValkey(); err != nil {
		cmd.Log.Fatal("Failed to initialize Valkey", err)
	}

	githubService, err := services.NewGithubService(config)
	if err != nil {
		cmd.Log.Fatal("Failed to initialize GithubService", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

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

	wg.Wait()
}
