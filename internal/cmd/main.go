package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/MatthewAraujo/airCast/internal/cmd/api"
	configs "github.com/MatthewAraujo/airCast/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	api := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.API.Port), logger)
	if err := api.Run(ctx); err != nil {
		logger.Error("failed to start app", slog.Any("error", err))
	}
}
