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

func ConsumeAchievementsStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.AutomaticEvents
	groupName := "achievements-group"
	consumerName := "achievements-consumer-1"

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
				var achievement models.Achievement
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					cmd.Log.Warn(fmt.Sprintf("Could not find data in message: %v", message.ID))
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &achievement)
				if err != nil {
					cmd.Log.Error("Error unmarshalling achievement", err)
					continue
				}

				cmd.Log.Info(fmt.Sprintf("Received achievement: %+v", achievement))

				var commentBody string
				switch achievement.Type {
				case "IMPACT":
					commentBody = fmt.Sprintf(cmd.HighImpact, achievement.ParticipantUsername)
				case "DOC":
					commentBody = fmt.Sprintf(cmd.DocSubmissions, achievement.ParticipantUsername)
				case "BUG":
					commentBody = fmt.Sprintf(cmd.BugReport, achievement.ParticipantUsername)
				case "TEST":
					commentBody = fmt.Sprintf(cmd.Tester, achievement.ParticipantUsername)
				case "HELP":
					commentBody = fmt.Sprintf(cmd.Helper, achievement.ParticipantUsername)
				default:
					cmd.Log.Warn(fmt.Sprintf("Unknown achievement type: %s", achievement.Type))
					cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
					continue
				}

				_, err = githubService.PostComment(achievement.Url, commentBody)
				if err != nil {
					cmd.Log.Error(fmt.Sprintf("Failed to post comment on %s", achievement.Url), err)
				}

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}
