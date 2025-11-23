package handlersTestsHelp

import (
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kir1l9x/avito-test-assignment/api/handlers"
	"github.com/kir1l9x/avito-test-assignment/api/router"
	"github.com/kir1l9x/avito-test-assignment/internal/appplication/services"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/repositories"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

func SetupTestServer(t *testing.T) (*httptest.Server, *pgxpool.Pool) {
	t.Helper()

	pool := repositories.SetupTestDB(t)
	repositories.TruncateAll(t, pool)

	logg, err := logger.New("debug")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	trm := manager.Must(pgxv5.NewDefaultFactory(pool))

	usersRepo := repositories.NewUsersRepository(pool, pgxv5.DefaultCtxGetter)
	teamsRepo := repositories.NewTeamsRepository(pool, pgxv5.DefaultCtxGetter)
	prRepo := repositories.NewPullRequestsRepository(pool, pgxv5.DefaultCtxGetter)

	usersService := services.NewUsersService(usersRepo)
	teamsService := services.NewTeamService(teamsRepo, usersRepo, trm)
	prService := services.NewPullRequestService(prRepo, usersRepo, trm)

	usersHandler := handlers.NewUsersHandler(usersService)
	teamsHandler := handlers.NewTeamHandler(teamsService)
	prHandler := handlers.NewPullRequestHandler(prService)

	r := router.NewRouter(logg.Logger, teamsHandler, usersHandler, prHandler)
	r.Use(logger.Middleware(logg.Logger))

	engine := binding.Validator.Engine()
	v, ok := engine.(*validator.Validate)
	if !ok {
		t.Fatalf("unexpected validator type %T", engine)
	}

	err = validators.InitializeInValidator(v)
	if err != nil {
		t.Fatalf("failed to init validator: %v", err)
	}
	err = validators.InitDomainValidator()
	if err != nil {
		t.Fatalf("failed to init domain validator: %v", err)
	}

	ts := httptest.NewServer(r)

	return ts, pool
}

func TestDataPath(rel string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	base := filepath.Dir(file)
	return filepath.Join(base, "testdata", rel)
}
