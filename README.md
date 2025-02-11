# Dev Kit - Model Context Protocol (MCP) Server

[Tutorial](https://www.youtube.com/watch?v=XnDFtYKU6xU)

## Community

For community support, discussions, and updates, please visit our forum at [community.aiocean.io](https://community.aiocean.io/).


## Prerequisites

- Go 1.23.2 or higher
- Various API keys and tokens for the services you want to use

## Installation

### Installing via Go

1. Install the server:
```bash
go install github.com/nguyenvanduocit/dev-kit@latest
```

2. **Manual setup required** - Create a `.env` file with your configuration:
```env
ENABLE_TOOLS=
ATLASSIAN_HOST=
ATLASSIAN_EMAIL=
GITLAB_HOST=
GITLAB_TOKEN=
ATLASSIAN_TOKEN=
PROXY_URL=
OPENAI_API_KEY=
DEEPSEEK_API_KEY=
```

3. Config your claude's config:

```json{claude_desktop_config.json}
{
  "mcpServers": {
    "dev_kit": {
      "command": "dev-kit",
      "args": ["-env", "/path/to/.env"],
    }
  }
}
```

## Enable Tools

There are a hidden variable `ENABLE_TOOLS` in the environment variable. It is a comma separated list of tools group to enable. If not set, all tools will be enabled. Leave it empty to enable all tools.


Here is the list of tools group:

- `confluence`: Confluence tools
- `jira`: Jira tools
- `gitlab`: GitLab tools
- `script`: Script tools

## Available Tools

### confluence_search

Search Confluence

Arguments:

- `query` (String) (Required): Atlassian Confluence Query Language (CQL)

### confluence_get_page

Get Confluence page content

Arguments:

- `page_id` (String) (Required): Confluence page ID

### confluence_create_page

Create a new Confluence page

Arguments:

- `space_key` (String) (Required): The key of the space where the page will be created
- `title` (String) (Required): Title of the page
- `content` (String) (Required): Content of the page in storage format (XHTML)
- `parent_id` (String): ID of the parent page (optional)

### confluence_update_page

Update an existing Confluence page

Arguments:

- `page_id` (String) (Required): ID of the page to update
- `title` (String): New title of the page (optional)
- `content` (String): New content of the page in storage format (XHTML)
- `version_number` (String): Version number for optimistic locking (optional)

### gitlab_list_projects

List GitLab projects

Arguments:

- `group_id` (String) (Required): gitlab group ID
- `search` (String): Multiple terms can be provided, separated by an escaped space, either + or %20, and will be ANDed together. Example: one+two will match substrings one and two (in any order).

### gitlab_get_project

Get GitLab project details

Arguments:

- `project_path` (String) (Required): Project/repo path

### gitlab_list_mrs

List merge requests

Arguments:

- `project_path` (String) (Required): Project/repo path
- `state` (String) (Default: all): MR state (opened/closed/merged)

### gitlab_get_mr_details

Get merge request details

Arguments:

- `project_path` (String) (Required): Project/repo path
- `mr_iid` (String) (Required): Merge request IID

### gitlab_create_MR_note

Create a note on a merge request

Arguments:

- `project_path` (String) (Required): Project/repo path
- `mr_iid` (String) (Required): Merge request IID
- `comment` (String) (Required): Comment text

### gitlab_get_file_content

Get file content from a GitLab repository

Arguments:

- `project_path` (String) (Required): Project/repo path
- `file_path` (String) (Required): Path to the file in the repository
- `ref` (String) (Required): Branch name, tag, or commit SHA

### gitlab_list_pipelines

List pipelines for a GitLab project

Arguments:

- `project_path` (String) (Required): Project/repo path
- `status` (String) (Default: all): Pipeline status (running/pending/success/failed/canceled/skipped/all)

### gitlab_list_commits

List commits in a GitLab project within a date range

Arguments:

- `project_path` (String) (Required): Project/repo path
- `since` (String) (Required): Start date (YYYY-MM-DD)
- `until` (String): End date (YYYY-MM-DD). If not provided, defaults to current date
- `ref` (String) (Required): Branch name, tag, or commit SHA

### gitlab_get_commit_details

Get details of a commit

Arguments:

- `project_path` (String) (Required): Project/repo path
- `commit_sha` (String) (Required): Commit SHA

### gitlab_list_user_events

List GitLab user events within a date range

Arguments:

- `username` (String) (Required): GitLab username
- `since` (String) (Required): Start date (YYYY-MM-DD)
- `until` (String): End date (YYYY-MM-DD). If not provided, defaults to current date

### gitlab_list_group_users

List all users in a GitLab group

Arguments:

- `group_id` (String) (Required): GitLab group ID

### gitlab_create_mr

Create a new merge request

Arguments:

- `project_path` (String) (Required): Project/repo path
- `source_branch` (String) (Required): Source branch name
- `target_branch` (String) (Required): Target branch name
- `title` (String) (Required): Merge request title
- `description` (String): Merge request description

### jira_get_issue

Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions

Arguments:

- `issue_key` (String) (Required): The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)

### jira_search_issue

Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues

Arguments:

- `jql` (String) (Required): JQL query string (e.g., 'project = KP AND status = \"In Progress\"')

### jira_list_sprints

List all active and future sprints for a specific Jira board, including sprint IDs, names, states, and dates

Arguments:

- `board_id` (String) (Required): Numeric ID of the Jira board (can be found in board URL)

### jira_create_issue

Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL

Arguments:

- `project_key` (String) (Required): Project identifier where the issue will be created (e.g., KP, PROJ)
- `summary` (String) (Required): Brief title or headline of the issue
- `description` (String) (Required): Detailed explanation of the issue
- `issue_type` (String) (Required): Type of issue to create (common types: Bug, Task, Story, Epic)

### jira_update_issue

Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed

Arguments:

- `issue_key` (String) (Required): The unique identifier of the issue to update (e.g., KP-2)
- `summary` (String): New title for the issue (optional)
- `description` (String): New description for the issue (optional)

### jira_list_statuses

Retrieve all available issue status IDs and their names for a specific Jira project

Arguments:

- `project_key` (String) (Required): Project identifier (e.g., KP, PROJ)

### jira_transition_issue

Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue

Arguments:

- `issue_key` (String) (Required): The issue to transition (e.g., KP-123)
- `transition_id` (String) (Required): Transition ID from available transitions list
- `comment` (String): Optional comment to add with transition

### execute_comand_line_script

Safely execute command line scripts on the user's system with security restrictions. Features sandboxed execution, timeout protection, and output capture. Supports cross-platform scripting with automatic environment detection.

Arguments:

- `content` (String) (Required): 
- `interpreter` (String) (Default: /bin/sh): Path to interpreter binary (e.g. /bin/sh, /bin/bash, /usr/bin/python, cmd.exe). Validated against allowed list for security
- `working_dir` (String): Execution directory path (default: user home). Validated to prevent unauthorized access to system locations

