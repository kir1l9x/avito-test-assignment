package dbModels

import "time"

type PullRequest struct {
	ID        string
	Name      string
	AuthorID  string
	StatusID  int16
	Reviewers []string
	CreatedAt time.Time
	MergedAt  *time.Time
}
