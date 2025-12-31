package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IssueAction struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"url"`
	Claim               bool   `json:"claimed"`
	Extend              bool   `json:"extend"`
}

type BountyAction struct {
	ParticipantUsername string `json:"github_username"`
	Amount              int    `json:"amount"`
	Url                 string `json:"url"`
	Action              string `json:"action"`
}

type Achievement struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"url"`
	Type                string `json:"type"`
}

type Solution struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"pull_request_url"`
	Merged              bool   `json:"merged"`
}

func ReadIssueStream(bm *BotMux) {
	lastEntry, err := ReadLastEntry("issue")
	if err != nil {
		Log.Error("Could not setup read-issue streeam", err)
		return
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{"issue-stream", lastEntry.StreamId},
			Count:   1,
			Block:   0,
		}
		streams, err := Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			Log.Error("failed to read from issue-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Retry logic
				for {
					entry, err := CreateEntry("issue", message.ID)
					if err != nil {
						Log.Error("failed to add event from issue-stream to local-db. Retrying in 2 seconds...", err)
						time.Sleep(2 * time.Second)
						continue
					}
					lastEntry = entry
					break
				}

				for _, val := range message.Values {

					var result IssueAction
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skipping the event as it's a malformed JSON
						Log.Error("Failed to unmarshal JSON at issue-stream", err)
						continue
					}
					mux.IssueClaim <- result
					Log.Info("Successfully sent issue-action-result to mux.")
				}
			}
		}
		// End of processing, reading next stream element
	}
}

func ReadBountyStream(bs *BotMux) {
	lastEntry, err := ReadLastEntry("bounty")
	if err != nil {
		Log.Error("Could not setup read-bounty stream", err)
		return
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{"bounty-stream", lastEntry.StreamId},
			Count:   1,
			Block:   0,
		}
		streams, err := Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			Log.Error("failed to read from bounty-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Retry logic
				for {
					entry, err := CreateEntry("bounty", message.ID)
					if err != nil {
						Log.Error("failed to add event from bounty-stream to local-db. Retrying in 2 seconds...", err)
						time.Sleep(2 * time.Second)
						continue
					}
					lastEntry = entry
					break
				}

				for _, val := range message.Values {

					var result BountyAction
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skipping the event as it's a malformed JSON
						Log.Error("Failed to unmarshal JSON at bounty-stream", err)
						continue
					}
					mux.Bounty <- result
					Log.Info("Successfully sent bounty-result to mux.")
				}
			}
		}
		// End of processing, reading next stream element
	}
}

func ReadSolutionStream(bs *BotMux) {
	lastEntry, err := ReadLastEntry("solution")
	if err != nil {
		Log.Error("Could not setup read-solution stream", err)
		return
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{"solution-merged-stream", lastEntry.StreamId},
			Count:   1,
			Block:   0,
		}
		streams, err := Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			Log.Error("failed to read from solution-merged-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Retry logic
				for {
					entry, err := CreateEntry("solution", message.ID)
					if err != nil {
						Log.Error("failed to add event from solution-merged-stream to local-db. Retrying in 2 seconds...", err)
						time.Sleep(2 * time.Second)
						continue
					}
					lastEntry = entry
					break
				}

				for _, val := range message.Values {

					var result Solution
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skipping the event as it's a malformed JSON
						Log.Error("Failed to unmarshal JSON at solution-merged-stream", err)
						continue
					}
					mux.Solution <- result
					Log.Info("Successfully sent solution-result to mux.")
				}
			}
		}
		// End of processing, reading next stream element
	}
}

func ReadAchivementStream(bs *BotMux) {
	lastEntry, err := ReadLastEntry("achivement")
	if err != nil {
		Log.Error("Could not setup read-achivement streeam", err)
		return
	}

	for {
		args := &redis.XReadArgs{
			Streams: []string{"automatic-events-stream", lastEntry.StreamId},
			Count:   1,
			Block:   0,
		}
		streams, err := Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(10 * time.Second)
				continue
			}
			Log.Error("failed to read from automatic-events-stream. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {

				// Retry logic
				for {
					entry, err := CreateEntry("achivement", message.ID)
					if err != nil {
						Log.Error("failed to add event from automatic-events-stream to local-db. Retrying in 2 seconds...", err)
						time.Sleep(2 * time.Second)
						continue
					}
					lastEntry = entry
					break
				}

				for _, val := range message.Values {

					var result Achievement
					err := json.Unmarshal([]byte(val.(string)), &result)
					if err != nil {
						// Skipping the event as it's a malformed JSON
						Log.Error("Failed to unmarshal JSON at automatic-events-stream", err)
						continue
					}
					mux.Achievement <- result
					Log.Info("Successfully sent achivement-result to mux.")
				}
			}
		}
		// End of processing, reading next stream element
	}
}
