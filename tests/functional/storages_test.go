package functional

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mohammadne/porsesh/internal/repository/storage"
)

func TestStoragePollOptions(t *testing.T) {
	var pollID int64 = 3

	t.Run("create_poll_options", func(t *testing.T) {
		tx, err := pollsStorage.StartTransaction(context.TODO())
		if err != nil {
			t.Fatalf("start transaction has error %s", err.Error())
		}

		err = pollsStorage.CreatePollOptions(context.TODO(), tx, pollID, []storage.PollOption{
			{Content: "option 1", Sort: 1},
			{Content: "option 2", Sort: 2},
			{Content: "option 3", Sort: 3},
		})
		if err != nil {
			t.Fatalf("create poll_options has error %s", err.Error())
		}

		tx.Commit()
	})

	t.Run("get_poll_options_by_poll_id", func(t *testing.T) {
		result, err := pollsStorage.GetPollOptionsByPollID(context.TODO(), pollID)
		if err != nil {
			t.Fatalf("get poll_options has error %s", err.Error())
		}

		bytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(bytes))
	})
}

func TestStoragePollTags(t *testing.T) {
	var pollID int64 = 3
	var tagIDs []int64 = []int64{1, 2}

	t.Run("create_poll_tags", func(t *testing.T) {
		tx, err := pollsStorage.StartTransaction(context.TODO())
		if err != nil {
			t.Fatalf("start transaction has error %s", err.Error())
		}

		err = pollsStorage.CreatePollTags(context.TODO(), tx, pollID, tagIDs)
		if err != nil {
			t.Fatalf("create poll_tags has error %s", err.Error())
		}

		tx.Commit()
	})
}

func TestStoragePolls(t *testing.T) {
	var creatorUserID int64 = 1

	t.Run("create_poll", func(t *testing.T) {
		tx, err := pollsStorage.StartTransaction(context.TODO())
		if err != nil {
			t.Fatalf("start transaction has error %s", err.Error())
		}

		id, err := pollsStorage.CreatePoll(context.TODO(), tx, &storage.Poll{
			UserID: creatorUserID,
			Title:  "some poll title",
		})
		if err != nil {
			t.Fatalf("create poll has error %s", err.Error())
		}

		tx.Commit()
		fmt.Println(id)
	})
}

func TestStorageTags(t *testing.T) {
	t.Run("create_tag", func(t *testing.T) {
		tx, err := pollsStorage.StartTransaction(context.TODO())
		if err != nil {
			t.Fatalf("start transaction has error %s", err.Error())
		}

		mapIds, err := tagsStorage.CreateTags(context.TODO(), tx, []storage.Tag{
			{Name: "tag20"},
			{Name: "tag30"},
		})
		if err != nil {
			t.Fatalf("create tag has error %s", err.Error())
		}

		tx.Commit()

		bytes, _ := json.MarshalIndent(mapIds, "", "  ")
		fmt.Println(string(bytes))
	})

	t.Run("get_tags_by_names", func(t *testing.T) {
		tx, err := pollsStorage.StartTransaction(context.TODO())
		if err != nil {
			t.Fatalf("start transaction has error %s", err.Error())
		}

		result, err := tagsStorage.GetTagsByNames(context.TODO(), tx, []string{"tag1", "tag2"})
		if err != nil {
			t.Fatalf("get tag has error %s", err.Error())
		}

		tx.Commit()

		bytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(bytes))
	})
}

func TestStorageVotes(t *testing.T) {
	var voterUserID int64 = 2
	var pollID int64 = 3
	var optionID int64 = 3

	t.Run("create_vote", func(t *testing.T) {
		t.Run("with_option", func(t *testing.T) {
			err = votesStorage.CreateVote(context.TODO(), &storage.Vote{
				UserID:   voterUserID,
				PollID:   pollID,
				OptionID: sql.NullInt64{Valid: true, Int64: 2},
			})
			if err != nil {
				t.Fatalf("create vote has error %s", err.Error())
			}
		})

		t.Run("no_option", func(t *testing.T) {
			err = votesStorage.CreateVote(context.TODO(), &storage.Vote{
				UserID:   voterUserID,
				PollID:   pollID,
				OptionID: sql.NullInt64{Valid: false},
			})
			if err != nil {
				t.Fatalf("create vote has error %s", err.Error())
			}
		})

	})

	t.Run("get_poll_option_votes_count", func(t *testing.T) {
		result, err := votesStorage.GetPollOptionVotesCount(context.TODO(), optionID)
		if err != nil {
			t.Fatalf("get option's vote count has error %s", err.Error())
		}
		fmt.Println(result)
	})

	t.Run("get_current_date_user_vote_count", func(t *testing.T) {
		result, err := votesStorage.GetCurrentDateUserVoteCount(context.TODO(), voterUserID)
		if err != nil {
			t.Fatalf("get voter's votes count has error %s", err.Error())
		}
		fmt.Println(result)
	})
}
