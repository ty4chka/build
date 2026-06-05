package build

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GitHubBuilder triggers GitHub Actions workflows
type GitHubBuilder struct {
	token  string
	owner  string
	repo   string
	client *http.Client
}

// NewGitHubBuilder creates a new builder
func NewGitHubBuilder(token, owner, repo string) *GitHubBuilder {
	return &GitHubBuilder{
		token:  token,
		owner:  owner,
		repo:   repo,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// BuildResult holds the build results
type BuildResult struct {
	RunID    int64
	Status   string
	Logs     string
	Artifact string
}

// TriggerBuild starts a GitHub Actions workflow
func (b *GitHubBuilder) TriggerBuild(ctx context.Context, goos, goarch string) (*BuildResult, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/build.yml/dispatches",
		b.owner, b.repo)

	payload := map[string]interface{}{
		"ref": "main",
		"inputs": map[string]string{
			"goos":   goos,
			"goarch": goarch,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %d, %s", resp.StatusCode, string(body))
	}

	return &BuildResult{Status: "triggered"}, nil
}

// GetRunStatus checks the status of a workflow run
func (b *GitHubBuilder) GetRunStatus(ctx context.Context, runID int64) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d",
		b.owner, b.repo, runID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Conclusion string `json:"conclusion"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status == "completed" {
		return result.Conclusion, nil
	}
	return result.Status, nil
}
