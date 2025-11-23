package storages

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/config"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
)

func CreateConnectionPool() (*pgxpool.Pool, error) {
	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w ", err)
	}

	dsn := dbCfg.GetDBDSNPostgres()

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	if dbCfg.DBMaxCons > math.MaxInt32 || dbCfg.DBMinCons > math.MaxInt32 {
		return nil, fmt.Errorf("invalid connection limits")
	}

	pgCfg.MaxConns = int32(dbCfg.DBMaxCons)
	pgCfg.MinConns = int32(dbCfg.DBMinCons)
	pgCfg.MaxConnLifetime = dbCfg.DBConLifetime
	pgCfg.ConnConfig.ConnectTimeout = dbCfg.DBConnTimeout

	ctx, cancel := context.WithTimeout(context.Background(), dbCfg.DBConnTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool with config: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	logg := logger.FromContext(ctx)

	logg.Info("connected to postgres")

	return pool, nil
}
