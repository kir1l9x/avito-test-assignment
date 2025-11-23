package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/kir1l9x/avito-test-assignment/api/handlers"
	"github.com/kir1l9x/avito-test-assignment/api/router"
	"github.com/kir1l9x/avito-test-assignment/internal/appplication/services"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/config"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/repositories"
	"github.com/kir1l9x/avito-test-assignment/internal/infrastructure/storages"
	"github.com/kir1l9x/avito-test-assignment/pkg/logger"
	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

func main() {
	appCfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get app config: %w", err))
	}

	logg, err := logger.New(appCfg.LogLevel)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create logger: %w", err))
	}

	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		logg.Fatal("failed to get db config", zap.Error(err))
	}

	pool, err := storages.CreateConnectionPool()
	if err != nil {
		logg.Fatal("failed to create connection pool", zap.Error(err))
	}

	logg.Info("connected to db", zap.String("host", dbCfg.DBHost))

	defer pool.Close()

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))

	usersRepo := repositories.NewUsersRepository(pool, trmpgx.DefaultCtxGetter)
	teamsRepo := repositories.NewTeamsRepository(pool, trmpgx.DefaultCtxGetter)
	prRepo := repositories.NewPullRequestsRepository(pool, trmpgx.DefaultCtxGetter)

	usersService := services.NewUsersService(usersRepo)
	teamsService := services.NewTeamService(teamsRepo, usersRepo, trManager)
	prService := services.NewPullRequestService(prRepo, usersRepo, trManager)

	usersHandler := handlers.NewUsersHandler(usersService)
	teamsHandler := handlers.NewTeamHandler(teamsService)
	prHandler := handlers.NewPullRequestHandler(prService)

	routers := router.NewRouter(logg.Logger, teamsHandler, usersHandler, prHandler)
	routers.Use(logger.Middleware(logg.Logger))

	engine := binding.Validator.Engine()
	v, ok := engine.(*validator.Validate)
	if !ok {
		logg.Fatal("unexpected validator engine type", zap.String("type", fmt.Sprintf("%T", engine)))
	}

	err = validators.InitializeInValidator(v)
	if err != nil {
		logg.Fatal("failed to initialize validator", zap.Error(err))
	}

	err = validators.InitDomainValidator()
	if err != nil {
		logg.Fatal("failed to initialize domain validator", zap.Error(err))
	}

	srv := http.Server{
		Addr:              appCfg.HTTPPort,
		Handler:           routers,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("url-service is running on %s\n", appCfg.HTTPPort)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("server forced to shutdown: %v", err)
	}
	cancel()

	log.Println("server exited gracefully")
}
