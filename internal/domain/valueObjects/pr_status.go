package valueObjects

type PullRequestStatus string

const (
	OpenPullRequest   PullRequestStatus = "OPEN"
	MergedPullRequest PullRequestStatus = "MERGED"
)

func (s PullRequestStatus) IsValid() bool {
	switch s {
	case OpenPullRequest, MergedPullRequest:
		return true
	}

	return false
}

func (s PullRequestStatus) String() string {
	return string(s)
}
