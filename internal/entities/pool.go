package entities

import "time"

type PollID int64

type Poll struct {
	ID        PollID
	Title     string
	UserID    UserID
	CreatedAt time.Time
	Options   []PollOption
	Tags      []PollTag
}

type PollOption struct {
	Content string
	Sort    int
}

type PollTag struct {
	Name string
}

type PollStatistics struct {
	PoolID PollID
	Votes  []PollStatisticsVote
}

type PollStatisticsVote struct {
	Option string
	Count  uint64
}
