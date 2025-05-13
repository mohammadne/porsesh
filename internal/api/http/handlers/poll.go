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
	g.Post("/", handler.todo)
}

type poll struct {
	logger *zap.Logger
}

func (s *poll) todo(c fiber.Ctx) error {
	response := &models.Response{}

	return response.Write(c, http.StatusInternalServerError)
}
