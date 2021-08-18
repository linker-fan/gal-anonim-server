package database

import (
	"linker-fan/gal-anonim-server/server/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var user = models.User{
	ID:           1,
	Username:     "test-username",
	PasswordHash: "test-password-hash",
	IsAdmin:      false,
	Created:      time.Now(),
	Updated:      time.Now(),
	Pin:          nil,
}

func TestInsertUser(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := "insert into users(id,username,passwordHash,isAdmin,created,updated) values (default, $1, $2, $3, $4, $5)"

	t.Run("Test Insert User", func(t *testing.T) {
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().WithArgs(user.Username, user.PasswordHash, user.IsAdmin, time.Now(), time.Now()).WillReturnResult(sqlmock.NewResult(0, 1))

		err := dw.InsertUser(user.Username, user.PasswordHash)
		assert.NoError(t, err)
	})

	t.Run("Test Insert User Err", func(t *testing.T) {

	})
}
