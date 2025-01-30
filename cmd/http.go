package cmd

import (
	"context"
	"fmt"
	"github.com/dhanarc/mobee-wallet/config"
	ledgerRepository "github.com/dhanarc/mobee-wallet/internal/ledger/postgresql"
	"github.com/dhanarc/mobee-wallet/internal/routes"
	"github.com/dhanarc/mobee-wallet/internal/user/auth"
	userRepository "github.com/dhanarc/mobee-wallet/internal/user/postgresql"
	"github.com/dhanarc/mobee-wallet/internal/wallet"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func HTTP(_ *cobra.Command, _ []string) {
	mainConfig := config.InitDefaultConfig()
	port := strconv.Itoa(mainConfig.HTTP.Port)

	db, err := gorm.Open(postgres.Open(mainConfig.Database.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	authClient := auth.NewClient(mainConfig.Auth.Key)

	userService := userRepository.NewClient(db, authClient)
	ledgerService := ledgerRepository.NewClient(db)

	walletService := wallet.NewClient(ledgerService, userService)

	handlers := routes.Init(&routes.Services{
		AuthClient:    authClient,
		WalletService: walletService,
	})

	server := &http.Server{
		Handler:           handlers,
		Addr:              fmt.Sprintf(":%s", port),
		WriteTimeout:      mainConfig.HTTP.WriteTimeout,
		IdleTimeout:       mainConfig.HTTP.IdleTimeout,
		ReadHeaderTimeout: mainConfig.HTTP.ReadHeaderTimeout,
		ReadTimeout:       mainConfig.HTTP.ReadTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Warn(err.Error())
		}
	}()

	slog.Info(fmt.Sprintf("http server started at port: %s", port))

	gracefulShutDown(server, mainConfig)
}

func gracefulShutDown(server *http.Server, cfg *config.Config) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig
	slog.Info("http server graceful shutdown started")

	time.AfterFunc(cfg.HTTP.GracefulShutdownTimeout, func() {
		slog.Error("failed to shutdown container")
		panic("failed to shutdown container")
	})

	err := server.Shutdown(context.Background())
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("http server closing remaining connection")
	err = server.Close()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("http server shutdown gracefully")

	os.Exit(0)
}
