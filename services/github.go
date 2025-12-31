package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v80/github"
	"golang.org/x/oauth2"
)

// PostComment posts a comment on a GitHub issue or pull request.
// It requires a GitHub Personal Access Token with `public_repo` or `repo` scope
// to be set in the GITHUB_TOKEN environment variable.
func PostComment(issueURL, body string) (*github.IssueComment, error) {
	// Get the GitHub token from an environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	// Parse the URL to get owner, repo, and issue number
	owner, repo, number, err := parseIssueURL(issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue URL: %w", err)
	}

	// Create an authenticated GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create the comment
	comment := &github.IssueComment{
		Body: &body,
	}

	// Post the comment
	createdComment, _, err := client.Issues.CreateComment(ctx, owner, repo, number, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	log.Printf("Successfully posted comment: %s\n", createdComment.GetHTMLURL())
	return createdComment, nil
}

// parseIssueURL extracts the owner, repository, and issue/PR number from a GitHub URL.
// It supports both issue and pull request URLs.
func parseIssueURL(issueURL string) (owner, repo string, number int, err error) {
	u, err := url.Parse(issueURL)
	if err != nil {
		return "", "", 0, err
	}

	// Regex to match /owner/repo/(issues|pull)/number
	re := regexp.MustCompile(`^/([^/]+)/([^/]+)/(?:issues|pull)/(\d+)$`)
	matches := re.FindStringSubmatch(u.Path)

	if len(matches) != 4 {
		return "", "", 0, fmt.Errorf("invalid GitHub issue or pull request URL format")
	}

	owner = matches[1]
	repo = matches[2]
	number, err = strconv.Atoi(matches[3])
	if err != nil {
		// This should not happen with the regex, but handle it just in case
		return "", "", 0, fmt.Errorf("invalid issue number: %s", matches[3])
	}

	return owner, repo, number, nil
}
