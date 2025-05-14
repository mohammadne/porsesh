package functional

import (
	"context"
	"fmt"
	"testing"

	"github.com/mohammadne/porsesh/internal/repository/storage"
)

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

}

func TestStorageVotes(t *testing.T) {

}
