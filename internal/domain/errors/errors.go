package errors

import "fmt"

type Code string

const (
	CodeTeamExists   Code = "TEAM_EXISTS"
	CodeTeamNotFound Code = "TEAM_NOT_FOUND"

	CodeUserNotFound Code = "USER_NOT_FOUND"

	CodePRExists   Code = "PR_EXISTS"
	CodePRNotFound Code = "PR_NOT_FOUND"
	CodePRMerged   Code = "PR_MERGED"

	CodeNotAssigned Code = "NOT_ASSIGNED"
	CodeNoCandidate Code = "NO_CANDIDATE"

	CodeInternal Code = "INTERNAL"
)

var (
	ErrTeamExists   = fmt.Errorf("team already exists")
	ErrTeamNotFound = fmt.Errorf("team not found")

	ErrUserNotFound = fmt.Errorf("user not found")

	ErrPRExists   = fmt.Errorf("pull request already exists")
	ErrPRNotFound = fmt.Errorf("pull request not found")
	ErrPRMerged   = fmt.Errorf("pull request already merged")

	ErrReviewerNotAssigned = fmt.Errorf("reviewer is not assigned to this PR")
	ErrNoReplacement       = fmt.Errorf("no active replacement candidate")
)

type DomainError struct {
	Code Code
	Err  error
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
}

func (e *DomainError) Unwrap() error { return e.Err }

func New(code Code, err error) *DomainError {
	return &DomainError{
		Code: code,
		Err:  err,
	}
}
