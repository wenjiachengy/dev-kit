# Dev Kit - Model Context Protocol (MCP) Server

Dev Kit is a powerful Model Context Protocol (MCP) server designed to enhance AI model capabilities by providing seamless integration with various development tools and services. It serves as a bridge between AI models (like Claude) and essential development services such as Jira, Confluence, and GitLab.

## Overview

This server implements the Model Context Protocol, allowing AI models to interact with development tools in a structured and secure way. It's particularly useful for automating development workflows, managing documentation, and handling version control operations through natural language interactions.

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
# Required for Atlassian services (Jira & Confluence)
ATLASSIAN_HOST=        # Your Atlassian instance URL (e.g., https://your-domain.atlassian.net)
ATLASSIAN_EMAIL=       # Your Atlassian account email
ATLASSIAN_TOKEN=       # Your Atlassian API token

# Required for GitLab services
GITLAB_HOST=           # Your GitLab instance URL
GITLAB_TOKEN=          # Your GitLab personal access token

# Optional configurations
ENABLE_TOOLS=          # Comma-separated list of tool groups to enable (empty = all enabled)
PROXY_URL=            # Optional: HTTP/HTTPS proxy URL if needed
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


## Available Tools

### Group: confluence

#### confluence_search

Search Confluence

#### confluence_get_page

Get Confluence page content

#### confluence_create_page

Create a new Confluence page

#### confluence_update_page

Update an existing Confluence page

### Group: gitlab

#### gitlab_list_projects

List GitLab projects

#### gitlab_get_project

Get GitLab project details

#### gitlab_list_mrs

List merge requests

#### gitlab_get_mr_details

Get merge request details

#### gitlab_create_MR_note

Create a note on a merge request

#### gitlab_get_file_content

Get file content from a GitLab repository

#### gitlab_list_pipelines

List pipelines for a GitLab project

#### gitlab_list_commits

List commits in a GitLab project within a date range

#### gitlab_get_commit_details

Get details of a commit

#### gitlab_list_user_events

List GitLab user events within a date range

#### gitlab_list_group_users

List all users in a GitLab group

#### gitlab_create_mr

Create a new merge request

### Group: jira

#### jira_get_issue

Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions

#### jira_search_issue

Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues

#### jira_list_sprints

List all active and future sprints for a specific Jira board, including sprint IDs, names, states, and dates

#### jira_create_issue

Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL

#### jira_update_issue

Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed

#### jira_list_statuses

Retrieve all available issue status IDs and their names for a specific Jira project

#### jira_transition_issue

Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue

### Group: script

#### execute_comand_line_script

Safely execute command line scripts on the user's system with security restrictions. Features sandboxed execution, timeout protection, and output capture. Supports cross-platform scripting with automatic environment detection.

