package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	//authenticate
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	//main query to github
	var query MainQuery
	//repo data we need
	var allRepos []repository
	repos := readRepos()
	for _, repo := range repos {

		variables := map[string]interface{}{
			"repo":   githubv4.String(repo),
			"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
		}

		for {
			err := client.Query(context.Background(), &query, variables)
			if err != nil {
				fmt.Println(err)
				return
			}
			allRepos = append(allRepos, query.Organization.Repositories.Nodes...)
			if !query.Organization.Repositories.PageInfo.HasNextPage {
				break
			}
			variables["cursor"] = githubv4.String(query.Organization.Repositories.PageInfo.EndCursor)
		}

	}

	if getRAWJSONFlag() {
		printJSON(allRepos)
	} else {
		writeCSV(allRepos)
	}

}

/// slurp all lines from file and use as input to GraphQL query
func readRepos() []string {

	filename := getOrgsFile()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return strings.Split(string(content), "\n")
}

/// write and process result from GraphQL query
/// calculate last activity, add custom comments in .csv etc
func writeCSV(repos []repository) {

	filepath := getCSVFile()
	csvFile, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Comma = ';'

	var header = []string{"Repo URL", "Primary language", "All languages", "Last commit", "PushedAt", "Total pull requests", "Calculated status"}
	writer.Write(header)
	for _, repo := range repos {
		var row []string
		row = append(row, string(repo.Url))
		row = append(row, string(repo.PrimaryLanguage.Name))
		var allLangs string
		for _, lang := range repo.Languages.Nodes {
			allLangs = allLangs + string(lang.Name) + "\n"
		}
		row = append(row, allLangs)

		row = append(row, lastActivity(repo).String())
		row = append(row, repo.PushedAt.String())

		row = append(row, strconv.Itoa(repo.PullRequests.TotalCount))
		row = append(row, calcStatus(repo))
		writer.Write(row)
	}
	// remember to flush!
	writer.Flush()

}

/// calculate last activity
/// last activity is rather last commit in contrary to PushedAt field
func lastActivity(repo repository) time.Time {

	if len(repo.Refs.Edges) == 0 {
		return repo.PushedAt
	}

	var lastActivity time.Time

	for _, node := range repo.Refs.Edges {
		if node.Node.Target.Commit.AuthoredDate.After(lastActivity) {
			lastActivity = node.Node.Target.Commit.AuthoredDate
		}
	}
	return lastActivity
}

/// based on last activity, PRs set different comments for status
func calcStatus(repo repository) string {

	if repo.IsArchived {
		return "Already archived"
	}
	if repo.IsEmpty {
		return "Delete - empty repo"
	}
	if repo.IsDisabled {
		return "Disabled repo"
	}

	//lasty activity
	lastActivity := lastActivity(repo)
	totalPullRequests := repo.PullRequests.TotalCount

	timespan := time.Since(lastActivity)

	diffYear := timespan.Hours() / 24.0 / 365.0
	diffMonth := timespan.Hours() / 24.0 / 30.0
	// log.Printf("diffYear %v", diffYear)
	// log.Printf("diffMonth %v", diffMonth)

	// immediatley mark as for archive older than 2 years
	if diffYear > 2 {
		return "Archive - inactive for more than 2 years"
	}

	// more than a year and les than twop years - tentative lets take a peek in pull requests
	if diffYear > 1 && diffYear < 2 && totalPullRequests < 5 {
		if repo.IsFork {
			return "Archive - tentative forked repo - inactive more than 1 year but with some PRs"
		} else {
			return "Archive - tentative - inactive more than 1 year and with few PRs"
		}
	}

	if diffYear > 1 && diffYear < 2 && totalPullRequests == 0 {
		if repo.IsFork {
			return "Archive - forked and inactive more than 1 year and no PRs at all"
		} else {
			return "Archive - inactive more than 1 year and no PRs at all"
		}
	}

	if diffMonth > 6 && totalPullRequests < 5 {
		return "OK - WARNING - there was no activity in last 6 months and very few PRs"
	}

	if diffMonth < 6 && totalPullRequests == 0 {
		return "OK - WARNING - there was activity in last 6 months but no PRs(trunk based development?)"
	}

	if diffMonth < 6 && totalPullRequests < 5 {
		return "OK - WARNING - there was activity in last 6 months and very few PRs"
	}

	return "OK"
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
