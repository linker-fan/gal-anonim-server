package database

import (
	"database/sql"
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
	Email:        "hyperxpizza@domain.com",
}

func TestInsertUser(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := "insert into users(id,username,passwordHash,isAdmin,created,updated,email) values (default, $1, $2, $3, $4, $5, $6)"

	t.Run("Test Insert User", func(t *testing.T) {
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().WithArgs(user.Username, user.PasswordHash, user.IsAdmin, time.Now(), time.Now()).WillReturnResult(sqlmock.NewResult(0, 1))

		err := dw.InsertUser(user.Username, user.PasswordHash, user.Email)
		assert.NoError(t, err)
	})

	t.Run("Test Insert User Err", func(t *testing.T) {

	})
}

func TestCheckIfUsernameExists(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("select username from users where username = $1")

	t.Run("Test CheckIfUsernameExists Err", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(user.ID, user.Username)
		mock.ExpectQuery(query).WithArgs(user.Username).WillReturnRows(rows)
		err := dw.CheckIfUsernameExists(user.Username)
		assert.Error(t, err)
	})

	t.Run("Test CheckIfUsernameExists SqlNoRowsErr", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
		err := dw.CheckIfUsernameExists(user.Username)
		assert.Error(t, err, sql.ErrNoRows)
	})
}

func TestGetIDAndPasswordByUsername(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("select id, passwordHash, isAdmin from users where username=$1")

	t.Run("Test GetIDAndPasswordByUsername No Err", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "passwordHash", "isAdmin"}).AddRow(user.ID, user.PasswordHash, user.IsAdmin)
		mock.ExpectQuery(query).WithArgs(user.Username).WillReturnRows(rows)
		id, passwordHash, isAdmin, err := dw.GetIDAndPasswordByUsername(user.Username)

		assert.Equal(t, id, user.ID)
		assert.Equal(t, passwordHash, user.PasswordHash)
		assert.Equal(t, isAdmin, user.IsAdmin)
		assert.NoError(t, err)

	})
}

func TestGetUserIDByUsername(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("select id from users where username=$1")

	t.Run("Test GetUserIDByUsername", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(user.ID)
		mock.ExpectQuery(query).WithArgs(user.Username).WillReturnRows(rows)
		id, err := dw.GetUserIDByUsername(user.Username)
		assert.Equal(t, id, user.ID)
		assert.NoError(t, err)
	})

	t.Run("Test GetUserIDByUserame Err", func(t *testing.T) {
		//rows := sqlmock.NewRows([]string{"id"}).AddRow(user.ID)
		mock.ExpectQuery(query).WithArgs(user.Username).WillReturnError(sql.ErrNoRows)
		id, err := dw.GetUserIDByUsername(user.Username)
		assert.Equal(t, id, 0)
		assert.Error(t, err, sql.ErrNoRows)
	})
}

func TestSetPin(t *testing.T) {
	db, mock, err := NewMock()
	if err != nil {
		t.Fail()
	}

	dw := DatabaseWrapper{db: db}
	defer dw.db.Close()

	query := regexp.QuoteMeta("update users set pin=$1 where id=$2")

	t.Run("Test SetPin", func(t *testing.T) {
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(user.Pin, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))
		err := dw.SetPin(*user.Pin, user.ID)
		assert.NoError(t, err)
	})
}
