package models

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (response *Response) Write(ctx fiber.Ctx, statusCode int) error {
	ctx.Set("Content-Type", "application/json")
	return ctx.Status(statusCode).JSON(&response)
}

type FeedResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Options   []string  `json:"options"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

type StatisticsResponse struct {
	PollID int      `json:"pollId"`
	Votes  []string `json:"votes"`
}
