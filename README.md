# gh-jira-pr

✨ A GitHub (`gh`) [CLI](https://cli.github.com/) extension to create GitHub pull requests based on Jira issues. The pull request will then be added to the Jira task as a web link.

Field mapping:

- Jira issue summary     -> Pull request title
- Jira issue description -> Pull request body
- Jira issue subtasks    -> Pull request body tasks

Note, when `--web` is used, the pull request will not be linked to the Jira issue.

## Setup

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
2. Install this extension:

    ```shell
    gh extension install ldaneliukas/gh-jira-pr
    ```

3. Login to GitHub `gh auth login`
4. Create [Jira Token](https://id.atlassian.com/manage-profile/security/api-tokens)

## Usage

Create a pull request from the supplied Jira ticket

```shell
USAGE:
 gh jira-pr <issue> [flags]


ARGUMENTS:
 issue                         <string>  Jira Issue


FLAGS:
 --web                         <string>  Open the web browser to create a pull request
 --ref                         <string>  Use the current repository HEAD ref as the Jira issue
 --jira-url   (env JIRA_URL)   <string>  Jira server URL
 --jira-user  (env JIRA_USER)  <string>  Jira username 
 --jira-token (env JIRA_TOKEN) <string>  Jira auth token


INHERITED FLAGS
 --help  Show help for command


EXAMPLES:
 $ gh jira-pr IT-1234
 $ gh jira-pr IT-1234 --jira-url https://company.atlassian.net
 $ gh jira-pr --ref
```
