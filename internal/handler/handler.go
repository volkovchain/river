package handler

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.midas.dev/back/river/internal/config"
	"gitlab.midas.dev/back/river/internal/service/salary"
)

// Handler handles CLI command execution
type Handler struct {
	db            *sql.DB
	salaryService *salary.Service
	config        *config.Config
}

// New creates a new handler
func New(db *sql.DB, salaryService *salary.Service, config *config.Config) *Handler {
	return &Handler{
		db:            db,
		salaryService: salaryService,
		config:        config,
	}
}

// Pay executes the pay command
func (h *Handler) Pay(ctx context.Context, isRepay bool) error {
	var err error
	if isRepay {
		err = h.salaryService.Repay(ctx)
	} else {
		err = h.salaryService.Pay(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to process salary payment: %w", err)
	}

	return nil
}
