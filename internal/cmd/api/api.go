package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	database "github.com/MatthewAraujo/airCast/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	addr   string
	logger *slog.Logger
	db     *pgxpool.Pool
}

func NewAPIServer(addr string, logger *slog.Logger) *APIServer {
	return &APIServer{
		addr:   addr,
		logger: logger,
	}
}

func (s *APIServer) Run(ctx context.Context) error {

	db, err := database.Connect(ctx, s.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.db = db

	router, err := s.loadRoutes()
	if err != nil {
		return fmt.Errorf("failed when loading routes: %w", err)
	}

	srv := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	errCh := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to listen and serve: %w", err)
		}

		close(errCh)
	}()

	s.logger.Info("server running")

	select {
	// Wait until we receive SIGINT (ctrl+c on cli)
	case <-ctx.Done():
		break
	case err := <-errCh:
		return err
	}

	sCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(sCtx)

	return nil
}
