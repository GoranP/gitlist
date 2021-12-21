package main

import "github.com/shurcooL/githubv4"
import "time"

// graphql search structs
type repository struct {
	Url             githubv4.String
	IsArchived      bool
	IsEmpty         bool
	IsFork          bool
	IsDisabled      bool
	PrimaryLanguage struct {
		Name githubv4.String
	}
	Languages struct {
		Nodes []struct {
			Name githubv4.String
		}
	} `graphql:"languages(first: 10)"`
	PushedAt     time.Time
	PullRequests struct {
		TotalCount int
	}
	Refs struct {
		Edges []struct {
			Node struct {
				Target struct {
					Commit struct {
						AuthoredDate time.Time
						//Url          githubv4.String //commit url
					} `graphql:"... on Commit"`
				}
			}
		}
	} `graphql:"refs(refPrefix: \"refs/\", first: 100)"`
}

type MainQuery struct {
	Organization struct {
		Repositories struct {
			TotalCount int
			Nodes      []repository
			PageInfo   struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: 20,after: $cursor)"`
	} `graphql:"organization(login: $repo)"`
}
