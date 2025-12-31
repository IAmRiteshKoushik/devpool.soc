package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"

	"github.com/IAmRiteshKoushik/devpool/cmd"
	"github.com/google/go-github/v80/github"
	"golang.org/x/oauth2"
)

type GithubService struct {
	client *github.Client
	ctx    context.Context
}

func NewGithubService(config *cmd.AppConfig) (*GithubService, error) {
	token := config.GithubToken
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GithubService{
			client: client,
			ctx:    ctx,
		},
		nil
}

// Post a comment on an issue or PR
func (s *GithubService) PostComment(issueURL, body string) (*github.IssueComment, error) {
	owner, repo, number, err := parseIssueURL(issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue URL: %w", err)
	}

	comment := &github.IssueComment{
		Body: &body,
	}

	createdComment, _, err := s.client.Issues.CreateComment(s.ctx, owner, repo, number, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	log.Printf("Successfully posted comment: %s\n", createdComment.GetHTMLURL())
	return createdComment, nil
}

func (s *GithubService) AssignIssue(issueURL string, assignee string) (*github.Issue, error) {
	owner, repo, number, err := parseIssueURL(issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue URL: %w", err)
	}

	assignees := []string{assignee}
	issue, _, err := s.client.Issues.AddAssignees(s.ctx, owner, repo, number, assignees)
	if err != nil {
		return nil, fmt.Errorf("failed to assign users: %w", err)
	}

	log.Printf("Successfully assigned %v to issue %s#%d", assignee, repo, number)
	return issue, nil
}

func (s *GithubService) UnassignIssue(issueURL string, assignee string) (*github.Issue, error) {
	owner, repo, number, err := parseIssueURL(issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue URL: %w", err)
	}

	assignees := []string{assignee}
	issue, _, err := s.client.Issues.RemoveAssignees(s.ctx, owner, repo, number, assignees)
	if err != nil {
		return nil, fmt.Errorf("failed to unassign users: %w", err)
	}

	log.Printf("Successfully unassigned %v from issue %s#%d", assignee, repo, number)
	return issue, nil
}

func parseIssueURL(issueURL string) (owner, repo string, number int, err error) {
	u, err := url.Parse(issueURL)
	if err != nil {
		return "", "", 0, err
	}

	// Regex for matching - /owner/repo/(issues|pull)/number
	re := regexp.MustCompile(`^/([^/]+)/([^/]+)/(?:issues|pull)/(\d+)$`)
	matches := re.FindStringSubmatch(u.Path)

	if len(matches) != 4 {
		return "", "", 0, fmt.Errorf("invalid GitHub issue or pull request URL format")
	}

	owner = matches[1]
	repo = matches[2]
	number, err = strconv.Atoi(matches[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid issue number: %s", matches[3])
	}

	return owner, repo, number, nil
}
