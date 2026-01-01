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

func ConsumeIssueStream(githubService *GithubService) {
	ctx := context.Background()
	streamName := cmd.IssueClaim
	groupName := "issue-group"
	consumerName := "issue-consumer-1"

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
				var issueAction models.IssueAction
				jsonData, ok := message.Values["data"].(string)
				if !ok {
					cmd.Log.Warn(fmt.Sprintf("Could not find data in message: %v", message.ID))
					continue
				}

				err := json.Unmarshal([]byte(jsonData), &issueAction)
				if err != nil {
					cmd.Log.Error("Error unmarshalling issue action", err)
					continue
				}

				cmd.Log.Info(fmt.Sprintf("Received issue action: %+v", issueAction))

				if issueAction.Claim {
					deadline := time.Now().Add(8 * 24 * time.Hour).Format("2006-01-02")
					commentBody := fmt.Sprintf(cmd.IssueClaimed, issueAction.ParticipantUsername, deadline)
					_, err := githubService.PostComment(issueAction.Url, commentBody)
					if err != nil {
						cmd.Log.Error(fmt.Sprintf("Failed to post comment on %s", issueAction.Url), err)
					}
					_, err = githubService.AssignIssue(issueAction.Url, issueAction.ParticipantUsername)
					if err != nil {
						cmd.Log.Error(fmt.Sprintf("Failed to assign issue %s to %s", issueAction.Url, issueAction.ParticipantUsername), err)
					}
				} else {
					commentBody := fmt.Sprintf(cmd.IssueUnclaimed, issueAction.ParticipantUsername)
					_, err := githubService.PostComment(issueAction.Url, commentBody)
					if err != nil {
						cmd.Log.Error(fmt.Sprintf("Failed to post comment on %s", issueAction.Url), err)
					}
					_, err = githubService.UnassignIssue(issueAction.Url, issueAction.ParticipantUsername)
					if err != nil {
						cmd.Log.Error(fmt.Sprintf("Failed to unassign issue %s to %s", issueAction.Url, issueAction.ParticipantUsername), err)
					}
				}

				cmd.Valkey.XAck(ctx, streamName, groupName, message.ID)
			}
		}
	}
}