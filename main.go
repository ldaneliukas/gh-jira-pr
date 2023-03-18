package main

import (
	"context"
	"fmt"
	"os"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/cli/go-gh"
	"gopkg.in/alecthomas/kingpin.v2"
	"mvdan.cc/xurls/v2"
)

func main() {
	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-ask failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func cli() error {

	jiraURL := kingpin.Flag("jira-url", "Jira URL").Envar("JIRA_URL").String()
	jiraUser := kingpin.Flag("jira-user", "Jira user").Envar("JIRA_USER").Required().String()
	jiraToken := kingpin.Flag("jira-token", "Jira token").Envar("JIRA_TOKEN").Required().String()
	jiraIssue := kingpin.Arg("jira-issue", "Jira issue").Required().String()
	kingpin.Parse()

	ctx := context.Background()

	tp := jira.BasicAuthTransport{
		Username: *jiraUser,
		APIToken: *jiraToken,
	}

	jiraClient, err := jira.NewClient(*jiraURL, tp.Client())
	if err != nil {
		return fmt.Errorf("could not create jira client: %v", err.Error())
	}
	issue, _, err := jiraClient.Issue.Get(ctx, *jiraIssue, nil)
	if err != nil {
		return fmt.Errorf("could not retrieve jira issue: %v", err.Error())
	}

	// Prepare pull request body
	body := ""
	body += fmt.Sprintf("## [%v: %v](%v/browse/%v)\n", *jiraIssue, issue.Fields.Summary, *jiraURL, *jiraIssue)
	body += fmt.Sprintf("### Description\n%v\n", issue.Fields.Description)
	body += "### Tasks\n"
	if issue.Fields.Subtasks != nil {
		for _, subtask := range issue.Fields.Subtasks {
			body += fmt.Sprintf("- [ ] %v\n", subtask.Fields.Summary)
		}
	}

	// Create pull request
	pr, _, err := gh.Exec("pr", "create", "--title", issue.Fields.Summary, "--body", body)
	if err != nil {
		return fmt.Errorf("could not create pull request: %v", err.Error())
	}

	// Get pull request URL
	xurlsStrict := xurls.Strict()
	prURL := xurlsStrict.FindAllString(pr.String(), -1)[0]

	// Update Jira issue with pull reqeust URL
	if prURL != "" {
		_, _, err = jiraClient.Issue.AddRemoteLink(ctx, *jiraIssue, &jira.RemoteLink{
			Object: &jira.RemoteLinkObject{
				URL:   prURL,
				Title: issue.Fields.Summary,
			},
		})
		if err != nil {
			return fmt.Errorf("could not add pull request link to Jira issue: %v", err.Error())
		}
	}
	return nil
}
