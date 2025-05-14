package models

import "github.com/mohammadne/porsesh/internal/entities"

// CreatePoll

type CreatePollRequestParams struct {
	UserID entities.UserID `json:"userId"`
}

type CreatePollRequest struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Tags    []string `json:"tags"`
}

// RetrieveFeed

type RetrieveFeedRequestParams struct {
	Tag    string          `json:"tag"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	UserID entities.UserID `json:"userId"`
}

// Vote

type VoteRequest struct {
	UserID      entities.UserID `json:"userId"`
	OptionIndex int             `json:"optionIndex"`
}

// Skip

type SkipRequest struct {
	UserID entities.UserID `json:"userId"`
}
