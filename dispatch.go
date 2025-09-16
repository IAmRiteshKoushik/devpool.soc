package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Function that determines the elapse of time 7 days from current time based on
// UTC.
func findEnd() string {
	nowInUTC := time.Now().UTC()
	future := nowInUTC.AddDate(0, 0, 7)
	result := future.Format("Monday, January 2, 2006")

	return result
}

func ManageExtension(username, url string) {
	comment := fmt.Sprintf(ExtensionGranted, username)

	info, err := ParseGitHubURL(url)
	if err != nil {
		Log.Error("Failed to parse github url", err)
		return
	}
	repoUrl := fmt.Sprintf("https://github.com/%s/%s", info.RepoOwner, info.RepoName)
	token, err := FetchInstallationToken(repoUrl)
	if err != redis.Nil {
		Log.Error("Could not fetch installation token", err)
		return
	}
	if err == redis.Nil {
		Log.Warn("Cache miss. Installation token not found. Trying to generate")
		token, err = NewInstallationToken(url)
		if err != nil {
			Log.Error("Failed to obtain new installation-token from GitHub", err)
			return
		}
	}
	if token == "" {
		Log.Error("No token returned by GitHub after request",
			fmt.Errorf("Missing token despite successful GitHub API call"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = postComment(ctx, token, info.RepoOwner, info.RepoName, info.Number,
		comment, info.Type)
	if err != nil {
		Log.Error("Failed to post comment on github", err)
		return
	}
	Log.Info(fmt.Sprintf("Successfully posted comment on github %s", repoUrl))
}

func ManageIssueClaim(username string, claim bool, url string) {
	timeout := findEnd()
	var comment string
	if claim {
		comment = fmt.Sprintf(IssueClaimed, username, timeout)
	} else {
		comment = fmt.Sprintf(IssueUnclaimed, username)
	}

	info, err := ParseGitHubURL(url)
	if err != nil {
		Log.Error("Failed to parse github issue url", err)
		return
	}
	repoUrl := fmt.Sprintf("https://github.com/%s/%s", info.RepoOwner, info.RepoName)
	token, err := FetchInstallationToken(repoUrl)
	if err != redis.Nil {
		Log.Error("Could not fetch installation token", err)
		return
	}
	if err == redis.Nil {
		Log.Warn("Cache miss. Installation token not found. Trying to generate")
		token, err = NewInstallationToken(url)
		if err != nil {
			Log.Error("Failed to obtain new installation-token from GitHub", err)
			return
		}
	}
	if token == "" {
		Log.Error("No token returned by GitHub after request",
			fmt.Errorf("Missing token despite successful GitHub API call"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = postComment(ctx, token, info.RepoOwner, info.RepoName, info.Number,
		comment, info.Type)
	if err != nil {
		Log.Error("Failed to post comment on github issue", err)
		return
	}
	Log.Info(fmt.Sprintf("Successfully posted comment on github issue %s", repoUrl))
}

func ManageBounty(username string, amt int, action, url string) {
	var comment string
	if action == "bounty" {
		comment = fmt.Sprintf(BountyDelivered, username)
	} else {
		comment = fmt.Sprintf(PenaltyDelivered, username)
	}

	info, err := ParseGitHubURL(url)
	if err != nil {
		Log.Error("Failed to parse github url", err)
		return
	}
	repoUrl := fmt.Sprintf("https://github.com/%s/%s", info.RepoOwner, info.RepoName)
	token, err := FetchInstallationToken(repoUrl)
	if err != redis.Nil {
		Log.Error("Could not fetch installation token", err)
		return
	}
	if err == redis.Nil {
		Log.Warn("Cache miss. Installation token not found. Trying to generate")
		token, err = NewInstallationToken(url)
		if err != nil {
			Log.Error("Failed to obtain new installation-token from GitHub", err)
			return
		}
	}
	if token == "" {
		Log.Error("No token returned by GitHub after request",
			fmt.Errorf("Missing token despite successful GitHub API call"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = postComment(ctx, token, info.RepoOwner, info.RepoName, info.Number,
		comment, info.Type)
	if err != nil {
		Log.Error("Failed to post comment on github", err)
		return
	}
	Log.Info(fmt.Sprintf("Successfully posted comment on github: %s", repoUrl))
}

func ManageSolution(username string, mergeStatus bool, url string) {
	var comment string
	if mergeStatus {
		comment = fmt.Sprintf(PRMerged, username)
	} else {
		comment = fmt.Sprintf(PROpened, username)
	}

	info, err := ParseGitHubURL(url)
	if err != nil {
		Log.Error("Failed to parse github pull-request url", err)
		return
	}
	repoUrl := fmt.Sprintf("https://github.com/%s/%s", info.RepoOwner, info.RepoName)
	token, err := FetchInstallationToken(repoUrl)
	if err != redis.Nil {
		Log.Error("Could not fetch installation token", err)
		return
	}
	if err == redis.Nil {
		Log.Warn("Cache miss. Installation token not found. Trying to generate")
		token, err = NewInstallationToken(url)
		if err != nil {
			Log.Error("Failed to obtain new installation-token from GitHub", err)
			return
		}
	}
	if token == "" {
		Log.Error("No token returned by GitHub after request",
			fmt.Errorf("Missing token despite successful GitHub API call"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = postComment(ctx, token, info.RepoOwner, info.RepoName, info.Number,
		comment, info.Type)
	if err != nil {
		Log.Error("Failed to post comment on github pull-request", err)
		return
	}
	Log.Info(fmt.Sprintf("Successfully posted comment on github pull-request %s", repoUrl))
}

func ManageAchivement(username string, aType string, url string) {
	var comment string
	switch aType {
	case "doc":
		comment = fmt.Sprintf(DocSubmissions, username)
	case "test":
		comment = fmt.Sprintf(Tester, username)
	case "help":
		comment = fmt.Sprintf(Helper, username)
	case "impact":
		comment = fmt.Sprintf(HighImpact, username)
	case "bug":
		comment = fmt.Sprintf(BugReport, username)
	default:
		Log.Warn("invalid parameter for achivement type")
		return
	}

	info, err := ParseGitHubURL(url)
	if err != nil {
		Log.Error("Failed to parse github url", err)
		return
	}
	repoUrl := fmt.Sprintf("https://github.com/%s/%s", info.RepoOwner, info.RepoName)
	token, err := FetchInstallationToken(repoUrl)
	if err != redis.Nil {
		Log.Error("Could not fetch installation token", err)
		return
	}
	if err == redis.Nil {
		Log.Warn("Cache miss. Installation token not found. Trying to generate")
		token, err = NewInstallationToken(url)
		if err != nil {
			Log.Error("Failed to obtain new installation-token from GitHub", err)
			return
		}
	}
	if token == "" {
		Log.Error("No token returned by GitHub after request",
			fmt.Errorf("Missing token despite successful GitHub API call"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = postComment(ctx, token, info.RepoOwner, info.RepoName, info.Number,
		comment, info.Type)
	if err != nil {
		Log.Error("Failed to post comment on github", err)
		return
	}
	Log.Info(fmt.Sprintf("Successfully posted comment on github %s", repoUrl))
}
