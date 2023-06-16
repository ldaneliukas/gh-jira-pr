package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	jira "github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/cli/go-gh"
	"github.com/go-git/go-git/v5"
	"gopkg.in/alecthomas/kingpin.v2"
	"mvdan.cc/xurls/v2"
)

func main() {
	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-jira-pr failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func cli() error {

	jiraURL := kingpin.Flag("jira-url", "Jira URL").Envar("JIRA_URL").String()
	jiraUser := kingpin.Flag("jira-user", "Jira username").Envar("JIRA_USER").Required().String()
	jiraToken := kingpin.Flag("jira-token", "Jira auth token").Envar("JIRA_TOKEN").Required().String()
	jiraIssue := kingpin.Arg("jira-issue", "Jira issue to base the pul request on").String()
	ref := kingpin.Flag("ref", "Use the current repository HEAD ref as the Jira issue").Bool()

	kingpin.Parse()

	ctx := context.Background()

	if *ref {
		r, err := git.PlainOpen(".")
		if err != nil {
			return fmt.Errorf("could not to load repository: %v", err.Error())
		}
		head, err := r.Head()
		if err != nil {
			return fmt.Errorf("could not retrieve head ref: %v", err.Error())
		}
		branch := strings.Split(head.String(), "refs/heads/")[1]
		*jiraIssue = branch
	}

	*jiraIssue = strings.ToUpper(*jiraIssue)
	if !isJiraIssue(*jiraIssue) {
		return fmt.Errorf("'%s' does not look like a Jira issue", *jiraIssue)
	}

	if *jiraIssue == "" {
		return fmt.Errorf("Jira issue was not provided as an argument and --ref was not specified")
	}

	tp := jira.BasicAuthTransport{
		Username: *jiraUser,
		APIToken: *jiraToken,
	}

	jiraClient, err := jira.NewClient(*jiraURL, tp.Client())
	if err != nil {
		return fmt.Errorf("could not create jira client: %v", err.Error())
	}
	options := &jira.GetQueryOptions{
		Expand: "renderedFields",
	}
	issue, _, err := jiraClient.Issue.Get(ctx, *jiraIssue, options)
	if err != nil {
		return fmt.Errorf("could not retrieve jira issue '%s': %v", *jiraIssue, err.Error())
	}

	// Prepare pull request body
	body := ""
	body += fmt.Sprintf("## [%v - %v](%v/browse/%v)\n\n", *jiraIssue, issue.Fields.Summary, *jiraURL, *jiraIssue)
	body += fmt.Sprintf("### Description\n%v\n\n", htmlToMarkdown(issue.RenderedFields.Description))
	if len(issue.Fields.Subtasks) > 0 {
		body += "### Tasks\n"
		for _, subtask := range issue.Fields.Subtasks {
			body += fmt.Sprintf("- [ ] %v\n", subtask.Fields.Summary)
		}
	}

	// Create pull request
	ghExecArgs := []string{"pr", "create", "--title", issue.Fields.Summary, "--body", body}
	pr, _, err := gh.Exec(ghExecArgs...)
	if err != nil {
		return fmt.Errorf("could not create pull request: %v", err.Error())
	}

	// Get pull request URL
	xurlsStrict := xurls.Strict()
	prURL := xurlsStrict.FindAllString(pr.String(), -1)

	// Update Jira issue with pull reqeust URL
	if len(prURL) > 0 {
		_, _, err = jiraClient.Issue.AddRemoteLink(ctx, *jiraIssue, &jira.RemoteLink{
			Object: &jira.RemoteLinkObject{
				URL:   prURL[0],
				Title: fmt.Sprintf("Pull Request: %v", issue.Fields.Summary),
			},
		})
		if err != nil {
			return fmt.Errorf("could not add pull request link to Jira issue: %v", err.Error())
		}
	}

	fmt.Printf("Pull request: %v\nJira issue: %vbrowse/%v\n", prURL[0], jiraClient.BaseURL, *jiraIssue)

	return nil
}

func isJiraIssue(s string) bool {
	pattern := "^[a-zA-Z]{1,10}-[0-9]{1,5}$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(s)
}

func htmlToMarkdown(s string) string {
	converter := md.NewConverter("", true, nil)
	markdown, _ := converter.ConvertString(s)
	return markdown
}
