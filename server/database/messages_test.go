package database

import (
	"linker-fan/gal-anonim-server/server/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var message = models.Message{
	ID:     1,
	RoomID: 1,
	UserID: 1,
	Text:   "sample-text",
}

func TestInsertMessage(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("insert into messages(id, roomID, userID, messageText, created) values(default, $1, $2, $3, $4)")

	t.Run("Test InsertMessage", func(t *testing.T) {
		created := time.Now()
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(message.RoomID, message.UserID, message.Text, created).WillReturnResult(sqlmock.NewResult(0, 1))
		err := dw.InsertMessage(message.RoomID, message.UserID, message.Text, created)
		assert.NoError(t, err)

	})
}
