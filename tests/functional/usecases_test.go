package functional

import (
	"context"
	"testing"

	"github.com/mohammadne/porsesh/internal/entities"
)

func TestUsecaseFeeds(t *testing.T) {

}

func TestUsecasePolls(t *testing.T) {
	t.Run("create_poll", func(t *testing.T) {
		poll := entities.Poll{
			Title:  "some poll2 title???",
			UserID: 1,
			Options: []entities.PollOption{
				{Content: "opt1", Sort: 1},
				{Content: "opt2", Sort: 2},
			},
			Tags: []entities.PollTag{
				{Name: "tag10"},
			},
		}

		err = pollsUsecase.CreatePoll(context.TODO(), &poll)
		if err != nil {
			t.Fatalf("create poll has error %s", err.Error())
		}
	})
}
