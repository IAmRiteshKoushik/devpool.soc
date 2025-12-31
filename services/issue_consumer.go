package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/IAmRiteshKoushik/devpool/models"
	"github.com/redis/go-redis/v9"
)

func ConsumeIssueStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.IssueClaim
	groupName := "issue-group"
	consumerName := "issue-consumer-1"

	err := cmd.Valkey.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Printf("Error creating consumer group: %v", err)
		}
	}

	for {
		streams, err := cmd.Valkey.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()

		if err != nil {
			log.Printf("Error reading from stream %s: %v", streamName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var issueAction models.IssueAction
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					log.Printf("Could not find data in message: %v", message.ID)
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &issueAction)
				if err != nil {
					log.Printf("Error unmarshalling issue action: %v", err)
					continue
				}

				log.Printf("Received issue action: %+v", issueAction)

				// Example of using the GithubService to post a comment
				if issueAction.Url != "" {
					commentBody := fmt.Sprintf("Hello @%s! Thanks for your interest in this issue. This is an automated message from DevPool.", issueAction.ParticipantUsername)
					_, err := githubService.PostComment(issueAction.Url, commentBody)
					if err != nil {
						log.Printf("Failed to post comment on %s: %v", issueAction.Url, err)
					}
				}

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}
