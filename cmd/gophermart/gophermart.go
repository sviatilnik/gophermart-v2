package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	accrual2 "github.com/sviatilnik/gophermart/internal/application/accrual"
	"github.com/sviatilnik/gophermart/internal/application/auth"
	"github.com/sviatilnik/gophermart/internal/application/order"
	"github.com/sviatilnik/gophermart/internal/application/wallet"
	configInfrastructure "github.com/sviatilnik/gophermart/internal/infrastructure/config"
	"github.com/sviatilnik/gophermart/internal/infrastructure/events"
	"github.com/sviatilnik/gophermart/internal/infrastructure/http/handlers"
	middlewareInfrastructure "github.com/sviatilnik/gophermart/internal/infrastructure/http/middleware"
	accrual4 "github.com/sviatilnik/gophermart/internal/infrastructure/persistence/accrual"
	authInfrastructure "github.com/sviatilnik/gophermart/internal/infrastructure/persistence/auth"
	orderInfrastructure "github.com/sviatilnik/gophermart/internal/infrastructure/persistence/order"
	"github.com/sviatilnik/gophermart/internal/infrastructure/persistence/user"
	walletInfrastructure "github.com/sviatilnik/gophermart/internal/infrastructure/persistence/wallet"
	accrual3 "github.com/sviatilnik/gophermart/internal/infrastructure/services/accrual"
	"github.com/sviatilnik/gophermart/internal/infrastructure/services/jwt"
	"go.uber.org/zap"
)

func main() {
	logger := getLogger()
	conf := getConfig()

	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		logger.Fatal(err)
	}

	err = execDBMigrations(db)
	if err != nil {
		logger.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(
		middlewareInfrastructure.GZIPCompress,
		middleware.Logger,
		middleware.RequestID,
		middleware.RealIP,
	)

	eventBus := events.NewInMemoryEventBus(logger)

	tokenGenerator := jwt.NewJWTGenerator(conf.AccessTokenSecret)
	refreshTokenRepo := authInfrastructure.NewRefreshTokenPostgresRepository(db)
	userRepo := user.NewPostgresUserRepository(db)
	regService := auth.NewRegistrationService(userRepo, eventBus)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, tokenGenerator)
	r.Post("/api/user/register", handlers.NewRegistrationHandler(regService, authService, logger).Register)

	authHandler := handlers.NewAuthHandler(authService, logger)
	r.Post("/api/user/login", authHandler.Login)
	r.Post("/api/user/login/refresh", authHandler.LoginByRefreshToken)

	r.Group(func(authRouter chi.Router) {
		authRouter.Use(middlewareInfrastructure.NewAuthMiddleware(jwt.NewVerifier(conf.AccessTokenSecret)).Handle)

		accRepo := accrual4.NewPostgresRepository(db)
		orderRepo := orderInfrastructure.NewOrderPostgresRepository(db)
		orderService := order.NewOrderService(orderRepo, userRepo, accRepo)
		orderHandler := handlers.NewOrderHandler(orderService)
		order.RegisterEventHandlers(eventBus, orderService)

		authRouter.Post("/api/user/orders", orderHandler.Create)
		authRouter.Get("/api/user/orders", orderHandler.GetList)

		walletRepo := walletInfrastructure.NewWalletPostgresRepository(db)
		walletService := wallet.NewWalletService(walletRepo, eventBus)
		walletHandler := handlers.NewWalletHandler(walletService)
		wallet.RegisterEventHandlers(eventBus, walletService)

		authRouter.Get("/api/user/balance", walletHandler.Balance)
		authRouter.Post("/api/user/balance/withdraw", walletHandler.Withdraw)
		authRouter.Get("/api/user/withdrawals", walletHandler.Withdrawals)

		accrual := accrual2.NewService(
			accrual3.NewHTTPClient(conf.AccrualSystemAddress),
			accRepo,
			orderService,
			eventBus,
			logger)

		go accrual.GetAccruals(context.Background())
	})

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    conf.Host,
		Handler: r,
	}

	go func() {
		logger.Info(fmt.Sprintf("start server on %s", server.Addr))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	<-quitChan

	logger.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal(err.Error())
	}

}

func getConfig() configInfrastructure.Config {
	return configInfrastructure.NewConfig(
		configInfrastructure.NewDefaultProvider(),
		configInfrastructure.NewFlagProvider(),
		configInfrastructure.NewEnvProvider(configInfrastructure.NewOSEnvGetter()),
	)
}

func getLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func execDBMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// https://github.com/golang-migrate/migrate/blob/master/source/file/README.md
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infrastructure/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
