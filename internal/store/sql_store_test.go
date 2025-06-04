package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/gebhn/auth-service/internal/db"
	"github.com/gebhn/auth-service/internal/db/sqlc"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var testStore *sqlStore

func TestMain(m *testing.M) {
	c, err := sql.Open("libsql", "file::memory:?cache=shared")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer c.Close()

	migrator := db.NewMigrator(c)

	if err := migrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}
	defer func() {
		if err := migrator.Down(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				log.Fatal(err)
			}
		}
	}()

	testStore = NewSqlStore(c)

	os.Exit(m.Run())
}

func clearTables(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec("DELETE FROM tokens; DELETE FROM users;")
	require.NoError(t, err, "failed to clear tables")
}

func insertUserHelper(t *testing.T) sqlc.CreateUserParams {
	t.Helper()

	params := sqlc.CreateUserParams{
		UserID:       "1",
		Username:     "username1",
		Email:        "username1@mail.me",
		PasswordHash: "pass",
	}
	err := testStore.CreateUser(context.Background(), params)
	assert.NoError(t, err)

	return params
}

func TestCreateUser_Success(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var user *sqlc.User

	params := sqlc.CreateUserParams{
		UserID:       "1",
		Username:     "username1",
		Email:        "username1@mail.me",
		PasswordHash: "pass",
	}

	err = testStore.CreateUser(context.Background(), params)
	assert.NoError(t, err)

	user, err = testStore.GetUserByEmail(context.Background(), "username1@mail.me")
	assert.NoError(t, err)
	assert.Equal(t, params.UserID, user.UserID)
	assert.Equal(t, params.Username, user.Username)
	assert.Equal(t, params.Email, user.Email)
	assert.Equal(t, params.PasswordHash, user.PasswordHash)
}

func TestCreateUser_Invalid(t *testing.T) {
	clearTables(t, testStore.db)

	var err error

	cases := []sqlc.CreateUserParams{
		{UserID: "", Username: "username1", Email: "username1@mail.me", PasswordHash: "pass"},
		{UserID: "1", Username: "", Email: "username1@mail.me", PasswordHash: "pass"},
		{UserID: "1", Username: "username1", Email: "", PasswordHash: "pass"},
		{UserID: "1", Username: "username1", Email: "username1@mail.me", PasswordHash: ""},
	}

	for _, tc := range cases {
		t.Run("Invalid Input", func(t *testing.T) {
			err = testStore.CreateUser(context.Background(), tc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), ErrInvalidInput.Error())
		})
	}
}

func TestUpdateUser_Success(t *testing.T) {
	clearTables(t, testStore.db)
	params := insertUserHelper(t)

	var err error
	var user *sqlc.User
	var userTwo *sqlc.User
	var userThree *sqlc.User

	t.Run("Should update Email and Username", func(t *testing.T) {
		err = testStore.UpdateUser(context.Background(), sqlc.UpdateUserParams{
			UserID:   params.UserID,
			Username: "newUsername",
			Email:    "newEmail@mail.me",
		})
		assert.NoError(t, err)

		user, err = testStore.GetUserByID(context.Background(), params.UserID)
		assert.NoError(t, err)
		assert.NotEqual(t, params.Username, user.Username)
		assert.NotEqual(t, params.Email, user.Email)
	})

	t.Run("Should update Email only", func(t *testing.T) {
		err = testStore.UpdateUser(context.Background(), sqlc.UpdateUserParams{
			UserID:   params.UserID,
			Username: "newUsername2",
			Email:    "",
		})
		assert.NoError(t, err)

		userTwo, err = testStore.GetUserByID(context.Background(), params.UserID)
		assert.NoError(t, err)
		assert.NotEqual(t, user.Username, userTwo.Username)
		assert.Equal(t, user.Email, userTwo.Email)
	})

	t.Run("Should update Username only", func(t *testing.T) {
		err = testStore.UpdateUser(context.Background(), sqlc.UpdateUserParams{
			UserID:   params.UserID,
			Username: "",
			Email:    "newEmail2@mail.me",
		})
		assert.NoError(t, err)

		userThree, err = testStore.GetUserByID(context.Background(), params.UserID)
		assert.NoError(t, err)
		assert.NotEqual(t, userTwo.Email, userThree.Email)
		assert.Equal(t, userTwo.Username, userThree.Username)
	})
}

func TestUpdateUser_Invalid(t *testing.T) {}
