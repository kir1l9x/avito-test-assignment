package repositories

import (
	"context"
	"errors"
	"fmt"

	pgxTm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	domainErrors "github.com/kir1l9x/avito-test-assignment/internal/domain/errors"
	"github.com/kir1l9x/avito-test-assignment/internal/domain/models"
	repoErrors "github.com/kir1l9x/avito-test-assignment/internal/infrastructure/repositories/errors"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

const (
	tableTeamsName = "teams"
)

type TeamsRepository struct {
	db     *pgxpool.Pool
	getter *pgxTm.CtxGetter
}

func NewTeamsRepository(p *pgxpool.Pool, g *pgxTm.CtxGetter) *TeamsRepository {
	return &TeamsRepository{db: p, getter: g}
}

func (tr *TeamsRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	ds := dialect.
		Insert(goqu.T(tableTeamsName).Schema(schemaName)).
		Rows(goqu.Record{"name": team.Name}).
		Returning(goqu.T(tableTeamsName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build create team query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var teamName string

	conn := tr.getter.DefaultTrOrDB(ctx, tr.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(&teamName)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan created team", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	return &models.Team{
		Name:    teamName,
		Members: []*models.TeamMember{},
	}, nil
}

func (tr *TeamsRepository) GetByName(ctx context.Context, name string) (*models.Team, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	ds := dialect.
		Select(goqu.T(tableTeamsName).Schema(schemaName).All()).
		From(goqu.T(tableTeamsName).Schema(schemaName)).
		Where(goqu.Ex{"name": name})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get team query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("get query: %w", err))
	}

	var teamName string

	conn := tr.getter.DefaultTrOrDB(ctx, tr.db)
	row := conn.QueryRow(ctx, sql, args...)
	err = row.Scan(&teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.New(domainErrors.CodeTeamNotFound, domainErrors.ErrTeamNotFound)
		}

		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan get team query", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	return &models.Team{
		Name:    teamName,
		Members: []*models.TeamMember{},
	}, nil
}
