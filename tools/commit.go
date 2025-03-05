package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/dev-kit/util"
)

// CommitData contains the information for the commit process
type CommitData struct {
	Thought            string `json:"thought"`
	CriticalQuestions  string `json:"critical_questions"`
	NextStep           string `json:"next_step"`
	Analysis           string `json:"analysis,omitempty"`
	ThoughtNumber      int    `json:"thoughtNumber"`
	TotalThoughts      int    `json:"totalThoughts"`
	NextThoughtNeeded  bool   `json:"nextThoughtNeeded"`
	IsRevision         *bool  `json:"isRevision,omitempty"`
	RevisesThought     *int   `json:"revisesThought,omitempty"`
	BranchFromThought  *int   `json:"branchFromThought,omitempty"`
	BranchID           string `json:"branchId,omitempty"`
	NeedsMoreThoughts  *bool  `json:"needsMoreThoughts,omitempty"`
}

// CommitServer manages the chain of thought process for the commit workflow
type CommitServer struct {
	thoughtHistory []CommitData
	branches       map[string][]CommitData
	// Store the final commit information
	changes        []GitChange
	commitType     string
	scope          string
	jiraIssueKey   string
	commitSubject  string
	commitBody     string
	breakingChange string
}

// GitChange represents a change to a file
type GitChange struct {
	Path     string
	Status   string
	Diff     string
	Filename string
}

// NewCommitServer creates a new CommitServer
func NewCommitServer() *CommitServer {
	return &CommitServer{
		thoughtHistory: make([]CommitData, 0),
		branches:       make(map[string][]CommitData),
		changes:        make([]GitChange, 0),
	}
}

// validateCommitData validates the commit data
func (s *CommitServer) validateCommitData(input map[string]interface{}) (*CommitData, error) {
	thought, ok := input["thought"].(string)
	if !ok || thought == "" {
		return nil, fmt.Errorf("invalid thought: must be a string")
	}

	criticalQuestions, ok := input["critical_questions"].(string)
	if !ok || criticalQuestions == "" {
		return nil, fmt.Errorf("invalid critical_questions: must be a string")
	}

	nextStep, ok := input["next_step"].(string)
	if !ok || nextStep == "" {
		return nil, fmt.Errorf("invalid next_step: must be a string")
	}

	thoughtNumber, ok := input["thoughtNumber"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid thoughtNumber: must be a number")
	}

	totalThoughts, ok := input["totalThoughts"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid totalThoughts: must be a number")
	}

	nextThoughtNeeded, ok := input["nextThoughtNeeded"].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid nextThoughtNeeded: must be a boolean")
	}

	data := &CommitData{
		Thought:           thought,
		CriticalQuestions: criticalQuestions,
		NextStep:          nextStep,
		ThoughtNumber:     int(thoughtNumber),
		TotalThoughts:     int(totalThoughts),
		NextThoughtNeeded: nextThoughtNeeded,
	}

	if analysis, ok := input["analysis"].(string); ok {
		data.Analysis = analysis
	}

	if isRevision, ok := input["isRevision"].(bool); ok {
		data.IsRevision = &isRevision
	}

	if revisesThought, ok := input["revisesThought"].(float64); ok {
		temp := int(revisesThought)
		data.RevisesThought = &temp
	}

	if branchFromThought, ok := input["branchFromThought"].(float64); ok {
		temp := int(branchFromThought)
		data.BranchFromThought = &temp
	}

	if branchID, ok := input["branchId"].(string); ok {
		data.BranchID = branchID
	}

	if needsMoreThoughts, ok := input["needsMoreThoughts"].(bool); ok {
		data.NeedsMoreThoughts = &needsMoreThoughts
	}

	return data, nil
}

// processThought processes the commit thought and updates the server state
func (s *CommitServer) processThought(input map[string]interface{}) (*mcp.CallToolResult, error) {
	data, err := s.validateCommitData(input)
	if err != nil {
		return nil, err
	}

	if data.BranchID != "" {
		// This is a branch thought
		if _, ok := s.branches[data.BranchID]; !ok {
			s.branches[data.BranchID] = make([]CommitData, 0)
		}
		s.branches[data.BranchID] = append(s.branches[data.BranchID], *data)
	} else {
		// This is a main thought
		s.thoughtHistory = append(s.thoughtHistory, *data)
	}

	// If this is the first thought, fetch git changes
	if data.ThoughtNumber == 1 {
		err := s.fetchGitChanges()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch git changes: %v", err)
		}
	}

	// Extract relevant information from the thought process
	s.extractCommitInfo(data)

	result := ""
	if data.NextThoughtNeeded {
		result = fmt.Sprintf("Thought %d processed. Proceed with your next thought.", data.ThoughtNumber)
	} else {
		// Generate final commit message
		commitMsg := s.generateCommitMessage()
		result = fmt.Sprintf("Commit process completed.\n\nFinal conventional commit message:\n\n%s", commitMsg)
	}

	return mcp.NewToolResultText(result), nil
}

// fetchGitChanges gets the current git changes
func (s *CommitServer) fetchGitChanges() error {
	// Run git diff --cached to get the staged changes
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git diff: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}

		status := parts[0]
		path := parts[1]
		filename := path[strings.LastIndex(path, "/")+1:]

		// Get the diff for this file
		diffCmd := exec.Command("git", "diff", "--cached", path)
		var diffOut bytes.Buffer
		diffCmd.Stdout = &diffOut
		if err := diffCmd.Run(); err != nil {
			return fmt.Errorf("failed to get diff for %s: %v", path, err)
		}

		change := GitChange{
			Path:     path,
			Status:   convertStatus(status),
			Diff:     diffOut.String(),
			Filename: filename,
		}
		s.changes = append(s.changes, change)
	}

	return nil
}

// convertStatus converts git status codes to human-readable status
func convertStatus(status string) string {
	switch status[0] {
	case 'A':
		return "Added"
	case 'M':
		return "Modified"
	case 'D':
		return "Deleted"
	case 'R':
		return "Renamed"
	case 'C':
		return "Copied"
	default:
		return status
	}
}

// extractCommitInfo extracts commit information from the thought process
func (s *CommitServer) extractCommitInfo(data *CommitData) {
	// Extract JIRA issue key using regex pattern
	jiraPattern := regexp.MustCompile(`[A-Z]+-\d+`)
	if jiraMatches := jiraPattern.FindStringSubmatch(data.Thought); len(jiraMatches) > 0 && s.jiraIssueKey == "" {
		s.jiraIssueKey = jiraMatches[0]
	}

	// Look for commit type and scope in the thought
	typePattern := regexp.MustCompile(`(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(([^)]+)\))?`)
	if typeMatches := typePattern.FindStringSubmatch(data.Thought); len(typeMatches) > 0 && s.commitType == "" {
		s.commitType = typeMatches[1]
		if len(typeMatches) > 3 {
			s.scope = typeMatches[3]
		}
	}

	// Extract commit subject from thought process
	if strings.Contains(data.Thought, "commit subject:") || strings.Contains(data.Thought, "commit message:") {
		lines := strings.Split(data.Thought, "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "commit subject:") || strings.Contains(strings.ToLower(line), "commit message:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					s.commitSubject = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	// Extract commit body from thought process
	if strings.Contains(data.Thought, "commit body:") {
		lines := strings.Split(data.Thought, "commit body:")
		if len(lines) > 1 {
			s.commitBody = strings.TrimSpace(lines[1])
		}
	}

	// Extract breaking change from thought process
	if strings.Contains(data.Thought, "BREAKING CHANGE:") {
		lines := strings.Split(data.Thought, "BREAKING CHANGE:")
		if len(lines) > 1 {
			s.breakingChange = strings.TrimSpace(lines[1])
		}
	}
}

// generateCommitMessage creates a conventional commit message
func (s *CommitServer) generateCommitMessage() string {
	// Default to "chore" if no commit type was extracted
	if s.commitType == "" {
		s.commitType = "chore"
	}

	// Build the header
	header := s.commitType
	if s.scope != "" {
		header += "(" + s.scope + ")"
	}
	header += ": "

	// Add breaking change indicator if there is a breaking change
	if s.breakingChange != "" {
		header += "!"
	}

	// Add the subject
	if s.commitSubject != "" {
		header += s.commitSubject
	} else {
		// Default subject if none was extracted
		header += "update code"
	}

	// Start building the full message
	message := header

	// Add Jira issue reference
	if s.jiraIssueKey != "" {
		message += "\n\nRef: " + s.jiraIssueKey
	}

	// Add the body if provided
	if s.commitBody != "" {
		message += "\n\n" + s.commitBody
	}

	// Add breaking changes section if there is a breaking change
	if s.breakingChange != "" {
		message += "\n\nBREAKING CHANGE: " + s.breakingChange
	}

	return message
}

// RegisterCommitTool registers the commit tool to the MCP server
func RegisterCommitTool(s *server.MCPServer) {
	commitServer := NewCommitServer()

	commitTool := mcp.NewTool("commit",
		mcp.WithDescription(`MCP Guide LLM for code commit process. This tool helps review changes, search for JIRA issues, and create conventional commit messages using a chain of thought process.

The tool follows a structured approach:
1. Reviews git changes to understand what has been modified
2. Searches for JIRA issue keys in the changes or branch name
3. Determines the appropriate conventional commit type (feat, fix, docs, etc.)
4. Creates a properly formatted commit message

Features:
- Uses chain of thought reasoning to understand code changes
- Automatically detects JIRA issue references
- Follows conventional commit format
- Allows for branching thought processes to explore alternatives
- Provides critical questions at each step to guide the reasoning process`),
		mcp.WithString("thought", mcp.Required(), mcp.Description("Your current thinking step")),
		mcp.WithString("critical_questions", mcp.Required(), mcp.Description("The critical questions of the current step, helps guide critical thinking")),
		mcp.WithString("next_step", mcp.Required(), mcp.Description("The next step to take in the process")),
		mcp.WithString("analysis", mcp.Description("The analysis of the current step")),
		mcp.WithBoolean("nextThoughtNeeded", mcp.Required(), mcp.Description("Whether another thought step is needed")),
		mcp.WithNumber("thoughtNumber", mcp.Required(), mcp.Description("Current thought number")),
		mcp.WithNumber("totalThoughts", mcp.Required(), mcp.Description("Estimated total thoughts needed")),
		mcp.WithBoolean("isRevision", mcp.Description("Whether this revises previous thinking")),
		mcp.WithNumber("revisesThought", mcp.Description("Which thought is being reconsidered")),
		mcp.WithNumber("branchFromThought", mcp.Description("Branching point thought number")),
		mcp.WithString("branchId", mcp.Description("Branch identifier")),
		mcp.WithBoolean("needsMoreThoughts", mcp.Description("If more thoughts are needed")),
	)

	s.AddTool(commitTool, util.ErrorGuard(commitServer.processThought))
}
