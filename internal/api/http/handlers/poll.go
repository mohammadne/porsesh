package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/porsesh/internal/api/http/models"
)

func NewPoll(r fiber.Router, logger *zap.Logger) {
	handler := &poll{
		logger: logger,
	}

	g := r.Group("poll")
	g.Post("/", handler.createPoll)
	g.Get("/", handler.retrieveFeed)
	g.Get("/:id/vote", handler.vote)
	g.Get("/:id/skip", handler.skip)
	g.Get("/:id/stats", handler.statistics)
}

type poll struct {
	logger *zap.Logger
}

func (s *poll) createPoll(c fiber.Ctx) error {
	response := &models.Response{}

	return response.Write(c, http.StatusInternalServerError)
}

func (s *poll) retrieveFeed(c fiber.Ctx) error {
	response := &models.Response{}

	return response.Write(c, http.StatusInternalServerError)
}

func (s *poll) vote(c fiber.Ctx) error {
	response := &models.Response{}

	id := c.Params("id")
	if len(id) == 0 {
		s.logger.Error("poll id not given for vote")
		return response.Write(c, fiber.StatusBadRequest)
	}

	return response.Write(c, http.StatusInternalServerError)
}

func (s *poll) skip(c fiber.Ctx) error {
	response := &models.Response{}

	id := c.Params("id")
	if len(id) == 0 {
		s.logger.Error("poll id not given for skip")
		return response.Write(c, fiber.StatusBadRequest)
	}

	return response.Write(c, http.StatusInternalServerError)
}

func (s *poll) statistics(c fiber.Ctx) error {
	response := &models.Response{}

	id := c.Params("id")
	if len(id) == 0 {
		s.logger.Error("poll id not given for statistics")
		return response.Write(c, fiber.StatusBadRequest)
	}

	return response.Write(c, http.StatusInternalServerError)
}
