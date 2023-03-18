# gh-jira-pr

✨ A GitHub (`gh`) [CLI](https://cli.github.com/) extension to create GitHub pull requests based on Jira issues.

## Setup

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
2. Install this extension:

    ```shell
    gh extension install ldaneliukas/gh-jira-pr
    ```

3. Login to GitHub `gh auth login`
4. Create [Jira Token](https://id.atlassian.com/manage-profile/security/api-tokens)

## Usage

```shell
gh jira-pr <command> [flags]
```

### Commands

Commands  | Description
--------- | -------------
create    | create a pull request from supplied Jira ticket

#### Create

Create a pull request from the supplied Jira ticket

```shell
USAGE:
 gh jira-pr create [flags]


ARGUMENTS:
 No Arguments


FLAGS:
 --jira-url   (env JIRA_URL)   <string>  Jira server URL
 --jira-user  (env JIRA_USER)  <string>  Jira username 
 --jira-token (env JIRA_TOKEN) <string>  Jira auth token


INHERITED FLAGS
 --help  Show help for command


EXAMPLES:
 $ gh jira-pr create IT-1234
 $ gh jira-pr create IT-1234 --jira-url https://company.atlassian.net
```
