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

func ConsumeBountyStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.Bounty
	groupName := "bounty-group"
	consumerName := "bounty-consumer-1"

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
				var bountyAction models.BountyAction
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					cmd.Log.Warn(fmt.Sprintf("Could not find data in message: %v", message.ID))
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &bountyAction)
				if err != nil {
					cmd.Log.Error("Error unmarshalling bounty action", err)
					continue
				}

				cmd.Log.Info(fmt.Sprintf("Received bounty action: %+v", bountyAction))

				var commentBody string
				switch bountyAction.Action {
				case "BOUNTY":
					commentBody = fmt.Sprintf(cmd.BountyDelivered, bountyAction.ParticipantUsername)
				case "PENALTY":
					commentBody = fmt.Sprintf(cmd.PenaltyDelivered, bountyAction.ParticipantUsername)
				default:
					cmd.Log.Warn(fmt.Sprintf("Unknown bounty action: %s", bountyAction.Action))
					cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
					continue
				}

				withRetry(func() error {
					_, err := githubService.PostComment(bountyAction.Url, commentBody)
					return err
				})

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}
