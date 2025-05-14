package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/mohammadne/porsesh/internal/api/http/models"
	"github.com/mohammadne/porsesh/internal/entities"
	"github.com/mohammadne/porsesh/internal/usecases"
)

func NewPoll(r fiber.Router, logger *zap.Logger, feeds usecases.Feeds, pools usecases.Polls) {
	handler := &poll{
		logger: logger,
		feeds:  feeds,
		pools:  pools,
	}

	g := r.Group("poll")
	g.Post("/", handler.createPoll)
	g.Get("/", handler.retrieveFeed)
	g.Post("/:id/vote", handler.vote)
	g.Post("/:id/skip", handler.skip)
	g.Get("/:id/stats", handler.statistics)
}

type poll struct {
	logger *zap.Logger
	// usecases
	feeds usecases.Feeds
	pools usecases.Polls
}

func (s *poll) createPoll(c fiber.Ctx) error {
	response := &models.Response{}

	request := models.CreatePollRequest{}
	if err := c.Bind().Body(&request); err != nil {
		s.logger.Error("invalid create poll request body", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	}

	params := models.CreatePollRequestParams{}
	if err := mapstructure.Decode(c.Queries(), &params); err != nil {
		s.logger.Error("invalid query parameters", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	} else if params.UserID <= 0 {
		s.logger.Error("user-id in parameters should be given", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	}

	poll := entities.Poll{
		Title:   request.Title,
		UserID:  params.UserID,
		Options: make([]entities.PollOption, 0, len(request.Options)),
		Tags:    make([]entities.PollTag, 0, len(request.Tags)),
	}
	for index, option := range request.Options {
		poll.Options = append(poll.Options, entities.PollOption{Content: option, Sort: index + 1})
	}
	for _, tag := range request.Tags {
		poll.Tags = append(poll.Tags, entities.PollTag{Name: tag})
	}

	if err := s.pools.CreatePoll(c.Context(), &poll); err != nil {
		s.logger.Error("error while creating pool", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidCreatePollArguments) {
			return response.Write(c, http.StatusBadRequest)
		}
		return response.Write(c, http.StatusInternalServerError)
	}

	return response.Write(c, http.StatusCreated)
}

func (s *poll) retrieveFeed(c fiber.Ctx) error {
	response := &models.Response{}

	params := models.RetrieveFeedRequestParams{}
	if err := mapstructure.Decode(c.Queries(), &params); err != nil {
		s.logger.Error("invalid query parameters", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	} else if params.UserID <= 0 {
		s.logger.Error("user-id in parameters should be given", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	}

	feed, err := s.feeds.GetUserFeed(c.Context(), params.UserID, params.Tag, params.Page, params.Limit)
	if err != nil {
		s.logger.Error("error while retrieving user feed", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidSkipPollArguments) {
			return response.Write(c, http.StatusBadRequest)
		}
		if errors.Is(err, usecases.ErrSkipPollPollNotExists) {
			return response.Write(c, http.StatusNotFound)
		}
		return response.Write(c, http.StatusInternalServerError)
	}

	response.Data = feed
	return response.Write(c, http.StatusOK)
}

func (s *poll) vote(c fiber.Ctx) error {
	response := &models.Response{}

	idRaw := c.Params("id")
	if len(idRaw) == 0 {
		s.logger.Error("poll id not given for vote")
		return response.Write(c, fiber.StatusBadRequest)
	}
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		s.logger.Error("poll id is invalid in vote")
		return response.Write(c, fiber.StatusBadRequest)
	}

	request := models.VoteRequest{}
	if err := c.Bind().Body(&request); err != nil {
		s.logger.Error("invalid vote request body", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	}

	if err := s.pools.VotePoll(c.Context(), entities.PollID(id), request.UserID, request.OptionIndex); err != nil {
		s.logger.Error("error while vote a pool", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidVotePollArguments) {
			return response.Write(c, http.StatusBadRequest)
		}
		if errors.Is(err, usecases.ErrDailyUserVotesLimit) {
			return response.Write(c, http.StatusForbidden)
		}
		if errors.Is(err, usecases.ErrVotePollPollNotExists) {
			return response.Write(c, http.StatusNotFound)
		}
		return response.Write(c, http.StatusInternalServerError)
	}

	return response.Write(c, http.StatusOK)
}

func (s *poll) skip(c fiber.Ctx) error {
	response := &models.Response{}

	idRaw := c.Params("id")
	if len(idRaw) == 0 {
		s.logger.Error("poll id not given for vote")
		return response.Write(c, fiber.StatusBadRequest)
	}
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		s.logger.Error("poll id is invalid in vote")
		return response.Write(c, fiber.StatusBadRequest)
	}

	request := models.SkipRequest{}
	if err := c.Bind().Body(&request); err != nil {
		s.logger.Error("invalid skip request body", zap.Error(err))
		return response.Write(c, fiber.StatusBadRequest)
	}

	if err := s.pools.SkipPoll(c.Context(), entities.PollID(id), request.UserID); err != nil {
		s.logger.Error("error while skip a pool", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidSkipPollArguments) {
			return response.Write(c, http.StatusBadRequest)
		}
		if errors.Is(err, usecases.ErrSkipPollPollNotExists) {
			return response.Write(c, http.StatusNotFound)
		}
		return response.Write(c, http.StatusInternalServerError)
	}

	return response.Write(c, http.StatusOK)
}

func (s *poll) statistics(c fiber.Ctx) error {
	response := &models.Response{}

	idRaw := c.Params("id")
	if len(idRaw) == 0 {
		s.logger.Error("poll id not given for vote")
		return response.Write(c, fiber.StatusBadRequest)
	}
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		s.logger.Error("poll id is invalid in vote")
		return response.Write(c, fiber.StatusBadRequest)
	}

	statistics, err := s.pools.Statistics(c.Context(), entities.PollID(id))
	if err != nil {
		s.logger.Error("error while skip a pool", zap.Error(err))
		if errors.Is(err, usecases.ErrInvalidSkipPollArguments) {
			return response.Write(c, http.StatusBadRequest)
		}
		if errors.Is(err, usecases.ErrSkipPollPollNotExists) {
			return response.Write(c, http.StatusNotFound)
		}
		return response.Write(c, http.StatusInternalServerError)
	}

	response.Data = statistics
	return response.Write(c, http.StatusOK)
}
