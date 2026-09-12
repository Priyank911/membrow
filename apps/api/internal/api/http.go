package api

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Priyank911/membrow/apps/api/internal/agent"
	"github.com/Priyank911/membrow/apps/api/internal/config"
	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
	"github.com/Priyank911/membrow/apps/api/internal/observability"
	"github.com/Priyank911/membrow/apps/api/internal/search"
)

type Server struct {
	app      *fiber.App
	agentSvc *agent.Service
	cfg      config.Config
}

type searchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type runRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

func NewServer(cfg config.Config, logger *slog.Logger, provider search.Provider, agentSvc *agent.Service) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout: cfg.RequestTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			appErr := apperrors.AsAppError(err)
			logger.Error("request failed", "path", c.Path(), "method", c.Method(), "request_id", observability.RequestIDFromCtx(c), "error", appErr.Error())
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error": fiber.Map{
					"code": appErr.Code,
					"message": appErr.Message,
					"details": appErr.Details,
					"retryable": appErr.Retryable,
					"request_id": observability.RequestIDFromCtx(c),
				},
			})
		},
	})

	app.Use(observability.RequestContextMiddleware())

	api := &Server{app: app, agentSvc: agentSvc, cfg: cfg}
	api.registerRoutes(provider)
	return api
}

func (s *Server) registerRoutes(provider search.Provider) {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "membrow-api"})
	})

	s.app.Post("/search", func(c *fiber.Ctx) error {
		var req searchRequest
		if err := c.BodyParser(&req); err != nil {
			return apperrors.InvalidInput("invalid json body", nil)
		}
		if req.Query == "" {
			return apperrors.InvalidInput("query is required", map[string]string{"field": "query"})
		}
		if req.Limit == 0 {
			req.Limit = s.cfg.DefaultSearchSize
		}
		ctx, cancel := context.WithTimeout(c.UserContext(), s.cfg.RequestTimeout)
		defer cancel()
		results, err := provider.Search(ctx, req.Query, req.Limit)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"data": fiber.Map{"results": results}, "request_id": observability.RequestIDFromCtx(c)})
	})

	s.app.Post("/agent/run", func(c *fiber.Ctx) error {
		var req runRequest
		if err := c.BodyParser(&req); err != nil {
			return apperrors.InvalidInput("invalid json body", nil)
		}
		if req.Limit == 0 {
			req.Limit = s.cfg.DefaultSearchSize
		}
		runID, err := s.agentSvc.StartRun(c.UserContext(), req.Query, req.Limit)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"data": fiber.Map{"run_id": runID, "status": agent.RunRunning},
			"request_id": observability.RequestIDFromCtx(c),
		})
	})

	s.app.Get("/agent/run/:id", func(c *fiber.Ctx) error {
		runID := c.Params("id")
		record, err := s.agentSvc.GetRun(c.UserContext(), runID)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"data": record, "request_id": observability.RequestIDFromCtx(c)})
	})
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	closeErr := make(chan error, 1)
	go func() {
		closeErr <- s.app.Shutdown()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-closeErr:
		if err != nil && !errors.Is(err, fiber.ErrServiceUnavailable) {
			return err
		}
		return nil
	}
}

func sleep(delay time.Duration) {
	<-time.After(delay)
}
