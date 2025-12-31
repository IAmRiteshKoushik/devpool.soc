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

func ConsumeSolutionStream() {
	ctx := context.Background()
	streamName := cmd.SolutionMerge
	groupName := "solution-group"
	consumerName := "solution-consumer-1"

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
				var solution models.Solution
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					log.Printf("Could not find data in message: %v", message.ID)
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &solution)
				if err != nil {
					log.Printf("Error unmarshalling solution: %v", err)
					continue
				}

				log.Printf("Received solution: %+v", solution)

				// TODO: Take appropriate action with the solution

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}
