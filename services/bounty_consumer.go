package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/IAmRiteshKoushik/devpool/models"
	"github.com/redis/go-redis/v9"
)

func ConsumeBountyStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.Bounty
	groupName := "bounty-group"
	consumerName := "bounty-consumer-1"

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
				var bountyAction models.BountyAction
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					log.Printf("Could not find data in message: %v", message.ID)
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &bountyAction)
				if err != nil {
					log.Printf("Error unmarshalling bounty action: %v", err)
					continue
				}

				log.Printf("Received bounty action: %+v", bountyAction)

				// TODO: Take appropriate action with the bountyAction

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}