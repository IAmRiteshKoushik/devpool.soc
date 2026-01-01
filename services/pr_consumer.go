package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/IAmRiteshKoushik/devpool/models"
	"github.com/redis/go-redis/v9"
)

func ConsumeSolutionStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.SolutionMerge
	groupName := "solution-group"
	consumerName := "solution-consumer-1"

	err := cmd.Valkey.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			cmd.Log.Error("Error creating consumer group", err)
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
			cmd.Log.Error(fmt.Sprintf("Error reading from stream %s", streamName), err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var solution models.Solution
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					cmd.Log.Warn(fmt.Sprintf("Could not find data in message: %v", message.ID))
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &solution)
				if err != nil {
					cmd.Log.Error("Error unmarshalling solution", err)
					continue
				}

				cmd.Log.Info(fmt.Sprintf("Received solution: %+v", solution))

				var commentBody string
				if solution.Merged {
					commentBody = fmt.Sprintf(cmd.PRMerged, solution.ParticipantUsername)
				} else {
					commentBody = fmt.Sprintf(cmd.PROpened, solution.ParticipantUsername)
				}

				_, err = githubService.PostComment(solution.Url, commentBody)
				if err != nil {
					cmd.Log.Error(fmt.Sprintf("Failed to post comment on %s", solution.Url), err)
				}

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}
