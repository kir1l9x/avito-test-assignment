package errors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
)

func MapPgError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if !ok {
		return err
	}

	switch pgErr.Code {

	case "23505":
		switch pgErr.ConstraintName {
		case "teams_pkey", "teams_name_key":
			return domainErrors.New(domainErrors.CodeTeamExists, domainErrors.ErrTeamExists)

		case "pull_requests_pkey", "pull_requests_id_key":
			return domainErrors.New(domainErrors.CodePRExists, domainErrors.ErrPRExists)

		default:
			return domainErrors.New(domainErrors.CodeInternal, err)
		}

	case "23503":
		switch pgErr.ConstraintName {
		case "pull_requests_author_id_fkey":
			return domainErrors.New(domainErrors.CodeUserNotFound, domainErrors.ErrUserNotFound)

		case "pr_reviewers_user_id_fkey":
			return domainErrors.New(domainErrors.CodeUserNotFound, domainErrors.ErrUserNotFound)

		case "users_team_name_fkey":
			return domainErrors.New(domainErrors.CodeTeamNotFound, domainErrors.ErrTeamNotFound)

		case "pull_requests_pkey":
			return domainErrors.New(domainErrors.CodePRNotFound, domainErrors.ErrPRNotFound)

		default:
			return domainErrors.New(domainErrors.CodeInternal, err)
		}

	case "23514":
		return domainErrors.New(domainErrors.CodeInternal, err)

	default:
		return domainErrors.New(domainErrors.CodeInternal, err)
	}
}
