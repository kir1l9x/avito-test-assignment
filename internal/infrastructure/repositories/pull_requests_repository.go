package repositories

import (
	"context"
	"errors"
	"fmt"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/valueObjects"
	"go.uber.org/zap"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/dbModels"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/mappers"
	repoErrors "github.com/kir1l9x/avito-test-assignment/internal/infrastructure/repositories/errors"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

const (
	tablePRName          = "pull_requests"
	tablePRReviewersName = "pr_reviewers"
	schemaName           = "public"
	dbDialect            = "postgres"
)

type PullRequestsRepository struct {
	db     *pgxpool.Pool
	getter *pgxTm.CtxGetter
	mapper *mappers.PullRequestStatusesMapper
}

func NewPullRequestsRepository(p *pgxpool.Pool, g *pgxTm.CtxGetter) *PullRequestsRepository {
	prr := &PullRequestsRepository{
		db:     p,
		getter: g,
		mapper: &mappers.PullRequestStatusesMapper{
			IDToStatus: make(map[int16]string),
			StatusToID: make(map[string]int16),
		},
	}

	prr.initStatusMapper(context.Background())

	return prr
}

func (prr *PullRequestsRepository) Create(ctx context.Context, pullRequest *models.PullRequest) (*models.PullRequest, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	statusID, err := prr.mapper.FromDomain(pullRequest.Status.String())
	if err != nil {
		logg.Error("failed mapping pull request status", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	ds := dialect.
		Insert(goqu.T(tablePRName).Schema(schemaName)).
		Rows(goqu.Record{
			"id":                pullRequest.ID,
			"pull_request_name": pullRequest.Name,
			"author_id":         pullRequest.AuthorID,
			"status_id":         statusID,
			"created_at":        pullRequest.CreatedAt,
			"merged_at":         pullRequest.MergedAt,
		}).
		Returning(goqu.T(tablePRName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build create pull request query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var createdPullRequest dbModels.PullRequest

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(
		&createdPullRequest.ID,
		&createdPullRequest.Name,
		&createdPullRequest.AuthorID,
		&createdPullRequest.StatusID,
		&createdPullRequest.CreatedAt,
		&createdPullRequest.MergedAt,
	)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan created pull request", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	reviewers, err := prr.createReviewers(ctx, pullRequest.ID, pullRequest.Reviewers)
	if err != nil {
		return nil, err
	}

	createdPullRequest.Reviewers = make([]string, 0, len(reviewers))
	for _, reviewer := range reviewers {
		createdPullRequest.Reviewers = append(createdPullRequest.Reviewers, reviewer.UserID)
	}

	pr, err := prr.mapper.ToDomainModel(&createdPullRequest)
	if err != nil {
		logg.Error("failed to map created pull request to domain", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	return pr, nil
}

func (prr *PullRequestsRepository) GetByID(ctx context.Context, pullRequestID string) (*models.PullRequest, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T(tablePRName).Schema(schemaName).All()).
		From(goqu.T(tablePRName).Schema(schemaName)).
		Where(goqu.Ex{"id": pullRequestID})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build get pull request by id query",
			zap.Error(err),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var pullRequest dbModels.PullRequest

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(
		&pullRequest.ID,
		&pullRequest.Name,
		&pullRequest.AuthorID,
		&pullRequest.StatusID,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.New(domainErrors.CodePRNotFound, domainErrors.ErrPRNotFound)
		}

		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to scan pull request by id",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, mapped
	}

	reviewers, err := prr.getReviewersByPRID(ctx, pullRequest.ID)
	if err != nil {
		return nil, err
	}

	pullRequest.Reviewers = make([]string, 0, len(reviewers))
	for _, reviewer := range reviewers {
		pullRequest.Reviewers = append(pullRequest.Reviewers, reviewer.UserID)
	}

	pr, err := prr.mapper.ToDomainModel(&pullRequest)
	if err != nil {
		logg.Error(
			"failed to map pull request to domain",
			zap.Error(err),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	return pr, nil
}

func (prr *PullRequestsRepository) Update(ctx context.Context, pullRequest *models.PullRequest) (*models.PullRequest, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	statusID, err := prr.mapper.FromDomain(pullRequest.Status.String())
	if err != nil {
		logg.Error("failed mapping pull request status", zap.Error(err), zap.String("pull_request_id", pullRequest.ID))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	ds := dialect.
		Update(goqu.T(tablePRName).Schema(schemaName)).
		Set(goqu.Record{
			"pull_request_name": pullRequest.Name,
			"author_id":         pullRequest.AuthorID,
			"status_id":         statusID,
			"created_at":        pullRequest.CreatedAt,
			"merged_at":         pullRequest.MergedAt,
		}).
		Where(goqu.Ex{"id": pullRequest.ID}).
		Returning(goqu.T(tablePRName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build update pull request query",
			zap.Error(err),
			zap.String("pull_request_id", pullRequest.ID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var updatedPullRequest dbModels.PullRequest

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	row := conn.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&updatedPullRequest.ID,
		&updatedPullRequest.Name,
		&updatedPullRequest.AuthorID,
		&updatedPullRequest.StatusID,
		&updatedPullRequest.CreatedAt,
		&updatedPullRequest.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.New(domainErrors.CodePRNotFound, domainErrors.ErrPRNotFound)
		}

		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to scan updated pull request",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("pull_request_id", pullRequest.ID),
		)
		return nil, mapped
	}

	if err = prr.deleteReviewersByPRID(ctx, pullRequest.ID); err != nil {
		return nil, err
	}

	reviewers, err := prr.createReviewers(ctx, pullRequest.ID, pullRequest.Reviewers)
	if err != nil {
		return nil, err
	}

	updatedPullRequest.Reviewers = make([]string, 0, len(reviewers))
	for _, r := range reviewers {
		updatedPullRequest.Reviewers = append(updatedPullRequest.Reviewers, r.UserID)
	}

	pr, err := prr.mapper.ToDomainModel(&updatedPullRequest)
	if err != nil {
		logg.Error(
			"failed to map updated pull request to domain",
			zap.Error(err),
			zap.String("pull_request_id", pullRequest.ID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	return pr, nil
}

func (prr *PullRequestsRepository) GetListByUser(ctx context.Context, userID string) ([]models.PullRequest, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T(tablePRName).Schema(schemaName).All()).
		From(goqu.T(tablePRName).Schema(schemaName)).
		Join(
			goqu.T(tablePRReviewersName).Schema(schemaName),
			goqu.On(goqu.Ex{
				"public.pull_requests.id": goqu.I("public.pr_reviewers.pr_id"),
			}),
		).
		Where(goqu.Ex{"public.pr_reviewers.user_id": userID})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build get pull requests by reviewer query",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to query pull requests by reviewer",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("user_id", userID),
		)
		return nil, mapped
	}
	defer rows.Close()

	pullRequests := make([]models.PullRequest, 0)

	for rows.Next() {
		var dbPR dbModels.PullRequest
		err = rows.Scan(
			&dbPR.ID,
			&dbPR.Name,
			&dbPR.AuthorID,
			&dbPR.StatusID,
			&dbPR.CreatedAt,
			&dbPR.MergedAt,
		)
		if err != nil {
			logg.Error(
				"failed to scan pull request row in list by user",
				zap.Error(err),
				zap.String("user_id", userID),
			)
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		reviewers, err := prr.getReviewersByPRID(ctx, dbPR.ID)
		if err != nil {
			return nil, err
		}

		dbPR.Reviewers = make([]string, 0, len(reviewers))
		for _, r := range reviewers {
			dbPR.Reviewers = append(dbPR.Reviewers, r.UserID)
		}

		pr, err := prr.mapper.ToDomainModel(&dbPR)
		if err != nil {
			logg.Error(
				"failed to map pull request to domain in list by user",
				zap.Error(err),
				zap.String("user_id", userID),
				zap.String("pull_request_id", dbPR.ID),
			)
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		pullRequests = append(pullRequests, *pr)
	}

	return pullRequests, nil
}

func (prr *PullRequestsRepository) GetReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	ds := dialect.
		Select(
			goqu.I("user_id"),
			goqu.COUNT("*").As("review_amount"),
		).
		From(goqu.T(tablePRReviewersName).Schema(schemaName)).
		GroupBy(goqu.I("user_id")).
		Order(goqu.I("review_amount").Desc())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get reviews amount by user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query reviews amount by user query", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	defer rows.Close()

	var stats []models.UserReviewStats

	for rows.Next() {
		var stat models.UserReviewStats
		err = rows.Scan(
			&stat.UserID,
			&stat.ReviewCount)
		if err != nil {
			logg.Error("failed to scan reviews amount by user row", zap.Error(err))
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (prr *PullRequestsRepository) GetTop5ReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	ds := dialect.
		Select(
			goqu.I("user_id"),
			goqu.COUNT("*").As("review_amount"),
		).
		From(goqu.T(tablePRReviewersName).Schema(schemaName)).
		GroupBy(goqu.I("user_id")).
		Order(goqu.I("review_amount").Desc()).
		Limit(5)

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get top 5 reviews amount by user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query reviews amount by user query", zap.Error(err))
		return nil, mapped
	}

	defer rows.Close()

	stats := make([]models.UserReviewStats, 0, 5)

	for rows.Next() {
		var stat models.UserReviewStats
		err = rows.Scan(
			&stat.UserID,
			&stat.ReviewCount)
		if err != nil {
			logg.Error("failed to scan reviews amount by user row", zap.Error(err))
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (prr *PullRequestsRepository) GetOpenReviewsAmountByUser(ctx context.Context) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	statusID, err := prr.mapper.FromDomain(valueObjects.OpenPullRequest.String())
	if err != nil {
		logg.Error("failed to map open status", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	ds := dialect.
		Select(
			goqu.I(tablePRReviewersName+".user_id"),
			goqu.COUNT("*").As("review_amount"),
		).
		From(goqu.T(tablePRReviewersName).Schema(schemaName)).
		Join(
			goqu.T(tablePRName).Schema(schemaName),
			goqu.On(goqu.Ex{
				tablePRReviewersName + ".pr_id": goqu.I(tablePRName + ".id"),
			}),
		).
		Where(goqu.Ex{tablePRName + ".status_id": statusID}).
		GroupBy(goqu.I(tablePRReviewersName + ".user_id")).
		Order(goqu.I("review_amount").Desc())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get open reviews amount by user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query reviews amount by user query", zap.Error(err))
		return nil, mapped
	}

	defer rows.Close()

	var stats []models.UserReviewStats
	for rows.Next() {
		var stat models.UserReviewStats
		err = rows.Scan(
			&stat.UserID,
			&stat.ReviewCount)
		if err != nil {
			logg.Error("failed to scan reviews amount by user row", zap.Error(err))
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (prr *PullRequestsRepository) GetTop5OpenReviewsAmountByUser(
	ctx context.Context,
) ([]models.UserReviewStats, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	statusID, err := prr.mapper.FromDomain(valueObjects.OpenPullRequest.String())
	if err != nil {
		logg.Error("failed to map open status", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	ds := dialect.
		Select(
			goqu.I(tablePRReviewersName+".user_id"),
			goqu.COUNT("*").As("review_amount"),
		).
		From(goqu.T(tablePRReviewersName).Schema(schemaName)).
		Join(
			goqu.T(tablePRName).Schema(schemaName),
			goqu.On(goqu.Ex{
				tablePRReviewersName + ".pr_id": goqu.I(tablePRName + ".id"),
			}),
		).
		Where(goqu.Ex{tablePRName + ".status_id": statusID}).
		GroupBy(goqu.I(tablePRReviewersName + ".user_id")).
		Order(goqu.I("review_amount").Desc()).
		Limit(5)

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get open reviews amount by user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, err)
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query reviews amount by user query", zap.Error(err))
		return nil, mapped
	}

	defer rows.Close()

	stats := make([]models.UserReviewStats, 0, 5)

	for rows.Next() {
		var stat models.UserReviewStats
		err = rows.Scan(
			&stat.UserID,
			&stat.ReviewCount)
		if err != nil {
			logg.Error("failed to scan reviews amount by user row", zap.Error(err))
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (prr *PullRequestsRepository) createReviewers(
	ctx context.Context,
	pullRequestID string,
	reviewers []string,
) ([]dbModels.Reviewer, error) {
	logg := logger.FromContext(ctx)

	if len(reviewers) == 0 {
		return []dbModels.Reviewer{}, nil
	}

	dialect := goqu.Dialect(dbDialect)
	records := make([]goqu.Record, 0, len(reviewers))
	for _, reviewer := range reviewers {
		records = append(records, goqu.Record{
			"pr_id":   pullRequestID,
			"user_id": reviewer,
		})
	}

	ds := dialect.
		Insert(goqu.T(tablePRReviewersName).Schema(schemaName)).
		Rows(records).
		Returning(goqu.T(tablePRReviewersName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build insert reviewers query",
			zap.Error(err),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to insert reviewers",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, mapped
	}
	defer rows.Close()

	prReviewers := make([]dbModels.Reviewer, 0, len(reviewers))

	for rows.Next() {
		var prReviewer dbModels.Reviewer
		err = rows.Scan(
			&prReviewer.PullRequestID,
			&prReviewer.UserID,
		)
		if err != nil {
			logg.Error(
				"failed to scan reviewer row after insert",
				zap.Error(err),
				zap.String("pull_request_id", pullRequestID),
			)
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		prReviewers = append(prReviewers, prReviewer)
	}

	return prReviewers, nil
}

func (prr *PullRequestsRepository) getReviewersByPRID(
	ctx context.Context,
	pullRequestID string,
) ([]dbModels.Reviewer, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T(tablePRReviewersName).Schema(schemaName).All()).
		From(goqu.T(tablePRReviewersName).Schema(schemaName)).
		Where(goqu.Ex{"pr_id": pullRequestID})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build get reviewers by pr_id query",
			zap.Error(err),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to query reviewers by pr_id",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("pull_request_id", pullRequestID),
		)
		return nil, mapped
	}
	defer rows.Close()

	prReviewers := make([]dbModels.Reviewer, 0)

	for rows.Next() {
		var prReviewer dbModels.Reviewer
		err = rows.Scan(
			&prReviewer.PullRequestID,
			&prReviewer.UserID,
		)
		if err != nil {
			logg.Error(
				"failed to scan reviewer row by pr_id",
				zap.Error(err),
				zap.String("pull_request_id", pullRequestID),
			)
			return nil, domainErrors.New(domainErrors.CodeInternal, err)
		}

		prReviewers = append(prReviewers, prReviewer)
	}

	return prReviewers, nil
}

func (prr *PullRequestsRepository) deleteReviewersByPRID(ctx context.Context, pullRequestID string) error {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Delete(goqu.T(tablePRReviewersName).Schema(schemaName)).
		Where(goqu.Ex{"pr_id": pullRequestID})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error(
			"failed to build delete reviewers query",
			zap.Error(err),
			zap.String("pull_request_id", pullRequestID),
		)
		return domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := prr.getter.DefaultTrOrDB(ctx, prr.db)

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error(
			"failed to delete reviewers by pr_id",
			zap.Error(err),
			zap.Error(mapped),
			zap.String("pull_request_id", pullRequestID),
		)
		return mapped
	}

	return nil
}

func (prr *PullRequestsRepository) initStatusMapper(ctx context.Context) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T("pr_statuses").Schema(schemaName).All()).
		From(goqu.T("pr_statuses").Schema(schemaName))

	sql, _, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build query for init status mapper", zap.Error(err))
		return
	}

	rows, err := prr.db.Query(ctx, sql)
	if err != nil {
		logg.Error("failed to query for init status mapper", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var id int16
		err = rows.Scan(&id, &status)
		if err != nil {
			logg.Error("failed to scan row for init status mapper", zap.Error(err))
			continue
		}

		prr.mapper.StatusToID[status] = id
		prr.mapper.IDToStatus[id] = status
	}
}
