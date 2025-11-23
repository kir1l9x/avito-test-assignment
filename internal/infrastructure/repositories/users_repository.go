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
	tableUsersName = "users"
)

type UsersRepository struct {
	db     *pgxpool.Pool
	getter *pgxTm.CtxGetter
}

func NewUsersRepository(p *pgxpool.Pool, g *pgxTm.CtxGetter) *UsersRepository {
	return &UsersRepository{db: p, getter: g}
}

func (ur *UsersRepository) CreateOrUpdate(ctx context.Context, user *models.User) (*models.User, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Insert(goqu.T(tableUsersName).Schema(schemaName)).
		Rows(goqu.Record{
			"id":        user.ID,
			"username":  user.Username,
			"team_name": user.TeamName,
			"is_active": user.IsActive,
		}).
		OnConflict(
			goqu.DoUpdate("id",
				goqu.Record{
					"username":  user.Username,
					"team_name": user.TeamName,
					"is_active": user.IsActive,
				}),
		).
		Returning(goqu.T(tableUsersName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build create/update user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var userModel models.User
	conn := ur.getter.DefaultTrOrDB(ctx, ur.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(&userModel.ID, &userModel.Username, &userModel.TeamName, &userModel.IsActive)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan create/update user", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	return &userModel, nil
}

func (ur *UsersRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T(tableUsersName).Schema(schemaName).All()).
		From(goqu.T(tableUsersName).Schema(schemaName)).
		Where(goqu.Ex{"id": id})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get user query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	var userModel models.User

	conn := ur.getter.DefaultTrOrDB(ctx, ur.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(&userModel.ID, &userModel.Username, &userModel.TeamName, &userModel.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.New(domainErrors.CodeUserNotFound, domainErrors.ErrUserNotFound)
		}

		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan user by id", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	return &userModel, nil
}

func (ur *UsersRepository) GetUsersByTeam(ctx context.Context, teamName string) ([]models.User, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)
	ds := dialect.
		Select(goqu.T(tableUsersName).Schema(schemaName).All()).
		From(goqu.T(tableUsersName).Schema(schemaName)).
		Where(goqu.Ex{"team_name": teamName})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get users by team query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := ur.getter.DefaultTrOrDB(ctx, ur.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query users by team", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	userModels := make([]models.User, 0)

	defer rows.Close()
	for rows.Next() {
		var userModel models.User

		err = rows.Scan(&userModel.ID, &userModel.Username, &userModel.TeamName, &userModel.IsActive)
		if err != nil {
			mapped := repoErrors.MapPgError(err)
			logg.Error("failed to scan users by team query", zap.Error(err), zap.Error(mapped))
			return nil, mapped
		}

		userModels = append(userModels, userModel)
	}

	return userModels, nil
}

func (ur *UsersRepository) GetActiveUsersInTeam(ctx context.Context, teamName string) ([]models.User, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Select(goqu.T(tableUsersName).Schema(schemaName).All()).
		From(goqu.T(tableUsersName).Schema(schemaName)).
		Where(goqu.Ex{
			"team_name": teamName,
			"is_active": true,
		})

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build get active users query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build query: %w", err))
	}

	conn := ur.getter.DefaultTrOrDB(ctx, ur.db)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to query active users", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}
	defer rows.Close()

	userModels := make([]models.User, 0)

	for rows.Next() {
		var userModel models.User

		err = rows.Scan(&userModel.ID, &userModel.Username, &userModel.TeamName, &userModel.IsActive)
		if err != nil {
			mapped := repoErrors.MapPgError(err)
			logg.Error("failed to scan active user", zap.Error(err), zap.Error(mapped))
			return nil, mapped
		}

		userModels = append(userModels, userModel)
	}

	return userModels, nil
}

func (ur *UsersRepository) SetIsActive(ctx context.Context, id string, isActive bool) (*models.User, error) {
	logg := logger.FromContext(ctx)
	dialect := goqu.Dialect(dbDialect)

	ds := dialect.
		Update(goqu.T(tableUsersName).Schema(schemaName)).
		Set(goqu.Record{"is_active": isActive}).
		Where(goqu.Ex{"id": id}).
		Returning(goqu.T(tableUsersName).Schema(schemaName).All())

	sql, args, err := ds.ToSQL()
	if err != nil {
		logg.Error("failed to build set isActive query", zap.Error(err))
		return nil, domainErrors.New(domainErrors.CodeInternal, fmt.Errorf("build update: %w", err))
	}

	var userModel models.User

	conn := ur.getter.DefaultTrOrDB(ctx, ur.db)
	row := conn.QueryRow(ctx, sql, args...)

	err = row.Scan(&userModel.ID, &userModel.Username, &userModel.TeamName, &userModel.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.New(domainErrors.CodeUserNotFound, domainErrors.ErrUserNotFound)
		}

		mapped := repoErrors.MapPgError(err)
		logg.Error("failed to scan updated user", zap.Error(err), zap.Error(mapped))
		return nil, mapped
	}

	return &userModel, nil
}
