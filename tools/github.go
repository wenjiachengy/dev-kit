package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v60/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/dev-kit/util"
)

var githubClient = sync.OnceValue[*github.Client](func() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is required")
	}

	client := github.NewClient(nil).WithAuthToken(token)
	return client
})

// RegisterGitHubTool registers the GitHub tool with the MCP server
func RegisterGitHubTool(s *server.MCPServer) {
	listReposTool := mcp.NewTool("github_list_repos",
		mcp.WithDescription("List GitHub repositories for a user or organization"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("GitHub username or organization name")),
		mcp.WithString("type", mcp.DefaultString("all"), mcp.Description("Type of repositories to list (all/owner/public/private/member)")),
	)

	repoDetailsTool := mcp.NewTool("github_get_repo",
		mcp.WithDescription("Get GitHub repository details"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
	)

	prListTool := mcp.NewTool("github_list_prs",
		mcp.WithDescription("List pull requests"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("state", mcp.DefaultString("open"), mcp.Description("PR state (open/closed/all)")),
	)

	prDetailsTool := mcp.NewTool("github_get_pr_details",
		mcp.WithDescription("Get pull request details"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Pull request number")),
	)

	prCommentTool := mcp.NewTool("github_create_pr_comment",
		mcp.WithDescription("Create a comment on a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Pull request number")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("Comment text")),
	)

	fileContentTool := mcp.NewTool("github_get_file_content",
		mcp.WithDescription("Get file content from a GitHub repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file in the repository")),
		mcp.WithString("ref", mcp.Description("Branch name, tag, or commit SHA")),
	)

	createPRTool := mcp.NewTool("github_create_pr",
		mcp.WithDescription("Create a new pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Pull request title")),
		mcp.WithString("head", mcp.Required(), mcp.Description("Name of the branch where your changes are implemented")),
		mcp.WithString("base", mcp.Required(), mcp.Description("Name of the branch you want your changes pulled into")),
		mcp.WithString("body", mcp.Description("Pull request description")),
	)

	prActionTool := mcp.NewTool("github_pr_action",
		mcp.WithDescription("Approve or close a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Pull request number")),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to take (approve/close)")),
	)

	issueListTool := mcp.NewTool("github_list_issues",
		mcp.WithDescription("List GitHub issues for a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("state", mcp.DefaultString("open"), mcp.Description("Issue state (open/closed/all)")),
		mcp.WithBoolean("include_body", mcp.DefaultBool(false), mcp.Description("Include issue description in the output")),
	)

	issueDetailsTool := mcp.NewTool("github_get_issue",
		mcp.WithDescription("Get GitHub issue details"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Issue number")),
	)

	issueCommentTool := mcp.NewTool("github_comment_issue",
		mcp.WithDescription("Comment on a GitHub issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Issue number")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("Comment text")),
	)

	issueActionTool := mcp.NewTool("github_issue_action",
		mcp.WithDescription("Close or reopen a GitHub issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("Repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name")),
		mcp.WithString("number", mcp.Required(), mcp.Description("Issue number")),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to take (close/reopen)")),
	)

	s.AddTool(listReposTool, util.ErrorGuard(listReposHandler))
	s.AddTool(repoDetailsTool, util.ErrorGuard(getRepoHandler))
	s.AddTool(prListTool, util.ErrorGuard(listPullRequestsHandler))
	s.AddTool(prDetailsTool, util.ErrorGuard(getPullRequestHandler))
	s.AddTool(prCommentTool, util.ErrorGuard(commentOnPullRequestHandler))
	s.AddTool(fileContentTool, util.ErrorGuard(getGitHubFileContentHandler))
	s.AddTool(createPRTool, util.ErrorGuard(createPullRequestHandler))
	s.AddTool(prActionTool, util.ErrorGuard(prActionHandler))
	s.AddTool(issueListTool, util.ErrorGuard(listIssuesHandler))
	s.AddTool(issueDetailsTool, util.ErrorGuard(getIssueHandler))
	s.AddTool(issueCommentTool, util.ErrorGuard(commentOnIssueHandler))
	s.AddTool(issueActionTool, util.ErrorGuard(issueActionHandler))
}

func listReposHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	ownerVal, ok := arguments["owner"]
	if !ok || ownerVal == nil {
		return nil, fmt.Errorf("missing required argument: owner")
	}
	owner, ok := ownerVal.(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner must be a non-empty string")
	}

	repoTypeVal, ok := arguments["type"]
	if !ok || repoTypeVal == nil {
		return nil, fmt.Errorf("missing required argument: type")
	}
	repoType, ok := repoTypeVal.(string)
	if !ok {
		return nil, fmt.Errorf("type must be a string")
	}

	opt := &github.RepositoryListOptions{
		Type: repoType,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var repos []*github.Repository
	var err error

	repos, _, err = githubClient().Repositories.List(context.Background(), owner, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}

	var result strings.Builder
	for _, repo := range repos {
		result.WriteString(fmt.Sprintf("Name: %s\n", repo.GetFullName()))
		result.WriteString(fmt.Sprintf("Description: %s\n", repo.GetDescription()))
		result.WriteString(fmt.Sprintf("URL: %s\n", repo.GetHTMLURL()))
		result.WriteString(fmt.Sprintf("Language: %s\n", repo.GetLanguage()))
		result.WriteString(fmt.Sprintf("Stars: %d\n", repo.GetStargazersCount()))
		result.WriteString(fmt.Sprintf("Forks: %d\n", repo.GetForksCount()))
		result.WriteString(fmt.Sprintf("Created: %s\n", repo.GetCreatedAt().Format("2006-01-02 15:04:05")))
		result.WriteString(fmt.Sprintf("Last Updated: %s\n\n", repo.GetUpdatedAt().Format("2006-01-02 15:04:05")))
	}

	return mcp.NewToolResultText(result.String()), nil
}

func getRepoHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)

	repository, _, err := githubClient().Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Repository Details:\n"))
	result.WriteString(fmt.Sprintf("Full Name: %s\n", repository.GetFullName()))
	result.WriteString(fmt.Sprintf("Description: %s\n", repository.GetDescription()))
	result.WriteString(fmt.Sprintf("URL: %s\n", repository.GetHTMLURL()))
	result.WriteString(fmt.Sprintf("Clone URL: %s\n", repository.GetCloneURL()))
	result.WriteString(fmt.Sprintf("Default Branch: %s\n", repository.GetDefaultBranch()))
	result.WriteString(fmt.Sprintf("Language: %s\n", repository.GetLanguage()))
	result.WriteString(fmt.Sprintf("Stars: %d\n", repository.GetStargazersCount()))
	result.WriteString(fmt.Sprintf("Forks: %d\n", repository.GetForksCount()))
	result.WriteString(fmt.Sprintf("Open Issues: %d\n", repository.GetOpenIssuesCount()))
	result.WriteString(fmt.Sprintf("Created: %s\n", repository.GetCreatedAt().Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("Last Updated: %s\n", repository.GetUpdatedAt().Format("2006-01-02 15:04:05")))

	return mcp.NewToolResultText(result.String()), nil
}

func listPullRequestsHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	state := arguments["state"].(string)

	opt := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	prs, _, err := githubClient().PullRequests.List(context.Background(), owner, repo, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %v", err)
	}

	var result strings.Builder
	for _, pr := range prs {
		result.WriteString(fmt.Sprintf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle()))
		result.WriteString(fmt.Sprintf("State: %s\n", pr.GetState()))
		result.WriteString(fmt.Sprintf("Author: %s\n", pr.GetUser().GetLogin()))
		result.WriteString(fmt.Sprintf("URL: %s\n", pr.GetHTMLURL()))
		result.WriteString(fmt.Sprintf("Created: %s\n", pr.GetCreatedAt().Format("2006-01-02 15:04:05")))
		if !pr.GetMergedAt().IsZero() {
			result.WriteString(fmt.Sprintf("Merged: %s\n", pr.GetMergedAt().Format("2006-01-02 15:04:05")))
		}
		if !pr.GetClosedAt().IsZero() {
			result.WriteString(fmt.Sprintf("Closed: %s\n", pr.GetClosedAt().Format("2006-01-02 15:04:05")))
		}
		result.WriteString(fmt.Sprintf("Base: %s\n", pr.GetBase().GetRef()))
		result.WriteString(fmt.Sprintf("Head: %s\n", pr.GetHead().GetRef()))
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

func getPullRequestHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)

	prNumber := 0
	fmt.Sscanf(number, "%d", &prNumber)

	pr, _, err := githubClient().PullRequests.Get(context.Background(), owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle()))
	result.WriteString(fmt.Sprintf("State: %s\n", pr.GetState()))
	result.WriteString(fmt.Sprintf("Author: %s\n", pr.GetUser().GetLogin()))
	result.WriteString(fmt.Sprintf("URL: %s\n", pr.GetHTMLURL()))
	result.WriteString(fmt.Sprintf("Created: %s\n", pr.GetCreatedAt().Format("2006-01-02 15:04:05")))
	if !pr.GetMergedAt().IsZero() {
		result.WriteString(fmt.Sprintf("Merged: %s\n", pr.GetMergedAt().Format("2006-01-02 15:04:05")))
	}
	if !pr.GetClosedAt().IsZero() {
		result.WriteString(fmt.Sprintf("Closed: %s\n", pr.GetClosedAt().Format("2006-01-02 15:04:05")))
	}
	result.WriteString(fmt.Sprintf("Base: %s\n", pr.GetBase().GetRef()))
	result.WriteString(fmt.Sprintf("Head: %s\n", pr.GetHead().GetRef()))
	result.WriteString(fmt.Sprintf("\nDescription:\n%s\n", pr.GetBody()))

	// Get PR comments
	comments, _, err := githubClient().Issues.ListComments(context.Background(), owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request comments: %v", err)
	}

	if len(comments) > 0 {
		result.WriteString("\nComments:\n")
		for _, comment := range comments {
			result.WriteString(fmt.Sprintf("\nFrom @%s at %s:\n%s\n",
				comment.GetUser().GetLogin(),
				comment.GetCreatedAt().Format("2006-01-02 15:04:05"),
				comment.GetBody()))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

func commentOnPullRequestHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)
	comment := arguments["comment"].(string)

	prNumber := 0
	fmt.Sscanf(number, "%d", &prNumber)

	issueComment := &github.IssueComment{
		Body: github.String(comment),
	}

	_, _, err := githubClient().Issues.CreateComment(context.Background(), owner, repo, prNumber, issueComment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %v", err)
	}

	return mcp.NewToolResultText("Comment created successfully"), nil
}

func getGitHubFileContentHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	path := arguments["path"].(string)
	ref := ""
	if refArg, ok := arguments["ref"]; ok {
		ref = refArg.(string)
	}

	opts := &github.RepositoryContentGetOptions{}
	if ref != "" {
		opts.Ref = ref
	}

	content, _, _, err := githubClient().Repositories.GetContents(context.Background(), owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %v", err)
	}

	decodedContent, err := content.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode content: %v", err)
	}

	return mcp.NewToolResultText(decodedContent), nil
}

func createPullRequestHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	title := arguments["title"].(string)
	head := arguments["head"].(string)
	base := arguments["base"].(string)
	body := ""
	if bodyArg, ok := arguments["body"]; ok {
		body = bodyArg.(string)
	}

	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}

	pr, _, err := githubClient().PullRequests.Create(context.Background(), owner, repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Pull request created successfully: %s", pr.GetHTMLURL())), nil
}

func prActionHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)
	action := arguments["action"].(string)

	prNumber := 0
	fmt.Sscanf(number, "%d", &prNumber)

	switch action {
	case "approve":
		// Create approval review
		review := &github.PullRequestReviewRequest{
			Event: github.String("APPROVE"),
		}
		_, _, err := githubClient().PullRequests.CreateReview(context.Background(), owner, repo, prNumber, review)
		if err != nil {
			return nil, fmt.Errorf("failed to approve pull request: %v", err)
		}
		return mcp.NewToolResultText("Pull request approved successfully"), nil

	case "close":
		// Close PR
		pr := &github.PullRequest{
			State: github.String("closed"),
		}
		_, _, err := githubClient().PullRequests.Edit(context.Background(), owner, repo, prNumber, pr)
		if err != nil {
			return nil, fmt.Errorf("failed to close pull request: %v", err)
		}
		return mcp.NewToolResultText("Pull request closed successfully"), nil

	default:
		return nil, fmt.Errorf("invalid action: %s. Must be either 'approve' or 'close'", action)
	}
}

func listIssuesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	state := arguments["state"].(string)
	includeBody := arguments["include_body"].(bool)

	opt := &github.IssueListByRepoOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	issues, _, err := githubClient().Issues.ListByRepo(context.Background(), owner, repo, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %v", err)
	}

	var result strings.Builder
	for _, issue := range issues {
		// Skip pull requests (they're also returned by the Issues API)
		if issue.IsPullRequest() {
			continue
		}

		result.WriteString(fmt.Sprintf("Issue #%d: %s\n", issue.GetNumber(), issue.GetTitle()))
		result.WriteString(fmt.Sprintf("State: %s\n", issue.GetState()))
		result.WriteString(fmt.Sprintf("Author: %s\n", issue.GetUser().GetLogin()))
		result.WriteString(fmt.Sprintf("URL: %s\n", issue.GetHTMLURL()))
		result.WriteString(fmt.Sprintf("Created: %s\n", issue.GetCreatedAt().Format("2006-01-02 15:04:05")))
		if !issue.GetClosedAt().IsZero() {
			result.WriteString(fmt.Sprintf("Closed: %s\n", issue.GetClosedAt().Format("2006-01-02 15:04:05")))
		}
		if len(issue.Labels) > 0 {
			labels := make([]string, 0, len(issue.Labels))
			for _, label := range issue.Labels {
				labels = append(labels, label.GetName())
			}
			result.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(labels, ", ")))
		}

		if includeBody && issue.GetBody() != "" {
			result.WriteString(fmt.Sprintf("Description:\n%s\n", issue.GetBody()))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

func getIssueHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)

	issueNumber := 0
	fmt.Sscanf(number, "%d", &issueNumber)

	issue, _, err := githubClient().Issues.Get(context.Background(), owner, repo, issueNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Issue #%d: %s\n", issue.GetNumber(), issue.GetTitle()))
	result.WriteString(fmt.Sprintf("State: %s\n", issue.GetState()))
	result.WriteString(fmt.Sprintf("Author: %s\n", issue.GetUser().GetLogin()))
	result.WriteString(fmt.Sprintf("URL: %s\n", issue.GetHTMLURL()))
	result.WriteString(fmt.Sprintf("Created: %s\n", issue.GetCreatedAt().Format("2006-01-02 15:04:05")))
	if !issue.GetClosedAt().IsZero() {
		result.WriteString(fmt.Sprintf("Closed: %s\n", issue.GetClosedAt().Format("2006-01-02 15:04:05")))
	}
	if len(issue.Labels) > 0 {
		labels := make([]string, 0, len(issue.Labels))
		for _, label := range issue.Labels {
			labels = append(labels, label.GetName())
		}
		result.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(labels, ", ")))
	}
	result.WriteString(fmt.Sprintf("\nDescription:\n%s\n", issue.GetBody()))

	// Get issue comments
	comments, _, err := githubClient().Issues.ListComments(context.Background(), owner, repo, issueNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue comments: %v", err)
	}

	if len(comments) > 0 {
		result.WriteString("\nComments:\n")
		for _, comment := range comments {
			result.WriteString(fmt.Sprintf("\nFrom @%s at %s:\n%s\n",
				comment.GetUser().GetLogin(),
				comment.GetCreatedAt().Format("2006-01-02 15:04:05"),
				comment.GetBody()))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

func commentOnIssueHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)
	comment := arguments["comment"].(string)

	issueNumber := 0
	fmt.Sscanf(number, "%d", &issueNumber)

	issueComment := &github.IssueComment{
		Body: github.String(comment),
	}

	_, _, err := githubClient().Issues.CreateComment(context.Background(), owner, repo, issueNumber, issueComment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %v", err)
	}

	return mcp.NewToolResultText("Comment created successfully"), nil
}

func issueActionHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	owner := arguments["owner"].(string)
	repo := arguments["repo"].(string)
	number := arguments["number"].(string)
	action := arguments["action"].(string)

	issueNumber := 0
	fmt.Sscanf(number, "%d", &issueNumber)

	var state string
	switch action {
	case "close":
		state = "closed"
	case "reopen":
		state = "open"
	default:
		return nil, fmt.Errorf("invalid action: %s. Must be either 'close' or 'reopen'", action)
	}

	issue := &github.IssueRequest{
		State: &state,
	}

	_, _, err := githubClient().Issues.Edit(context.Background(), owner, repo, issueNumber, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to %s issue: %v", action, err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Issue %s successfully", action+"d")), nil
}
