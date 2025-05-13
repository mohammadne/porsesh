package entities

import "time"

type PollID uint64

type Poll struct {
	ID        PollID
	Title     string
	UserID    UserID
	CreatedAt time.Time
	Options   []PollOption
	Tags      []PollTag
}

type PollOption struct {
	ID      uint64
	Content string
	Sort    int
}

type PollTag struct {
	ID   uint64
	Name string
}

type PollStatistics struct {
	PoolID PollID
	Votes  []struct {
		Option string
		Count  int
	}
}
