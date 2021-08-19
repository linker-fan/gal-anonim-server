package database

import (
	"linker-fan/gal-anonim-server/server/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var member = models.Member{
	ID:     1,
	RoomID: 1,
	UserID: 1,
}

func TestInsertMember(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("insert into members (id, roomID, userID, joined) values (default, $1, $2, $3)")

	t.Run("Test InsertMember", func(t *testing.T) {
		created := time.Now()
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(member.RoomID, member.UserID, created).WillReturnResult(sqlmock.NewResult(0, 1))
		err := dw.InsertMember(member.RoomID, member.UserID)
		assert.NoError(t, err)
	})
}
