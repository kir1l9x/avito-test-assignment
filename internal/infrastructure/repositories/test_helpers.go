package repositories

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/config"
)

var (
	testPool  *pgxpool.Pool //nolint:gochecknoglobals // shared pool for integration tests
	testDBURL string
	testDBMu  sync.Mutex
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	testDBMu.Lock()
	t.Cleanup(func() {
		testDBMu.Unlock()
	})

	if testPool != nil {
		return testPool
	}

	if testDBURL == "" {
		dbCfg := config.LoadTestDBConfig()
		testDBURL = dbCfg.GetDBDSNPostgres()
	}

	cfg, err := pgxpool.ParseConfig(testDBURL)
	if err != nil {
		t.Fatalf("cannot parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatalf("cannot connect pgxpool: %v", err)
	}

	waitForDB(t, pool)
	runMigrations(t)

	testPool = pool

	return pool
}

func waitForDB(t *testing.T, pool *pgxpool.Pool) {
	deadline := time.After(30 * time.Second)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := pool.Ping(ctx)
		cancel()

		if err == nil {
			return
		}

		select {
		case <-deadline:
			t.Fatalf("DB NOT READY: %v", err)
		default:
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func runMigrations(t *testing.T) {
	t.Helper()

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller failed")
	}
	projectRoot := filepath.Join(filepath.Dir(filePath), "..", "..", "..")

	migrationsPath := "file://" + filepath.Join(projectRoot, "migrations")

	m, err := migrate.New(
		migrationsPath,
		testDBURL,
	)
	if err != nil {
		t.Fatalf("migrate new: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up: %v", err)
	}
}

func TruncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE 
			pr_reviewers,
			pull_requests,
			users,
			teams 
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}
