package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/db"
	"github.com/gebhn/auth-service/internal/db/sqlc"

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

func insertTokenHelper(t *testing.T) sqlc.CreateTokenParams {
	t.Helper()

	params := sqlc.CreateTokenParams{
		Jti:       "jti",
		UserID:    "1",
		Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
		TokenHash: "hash",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	err := testStore.CreateToken(context.Background(), params)
	assert.NoError(t, err)

	return params
}

func TestExecTx_Success(t *testing.T) {
	clearTables(t, testStore.db)

	err := testStore.ExecTx(context.Background(), func(s Store) error {
		return s.CreateUser(context.Background(), sqlc.CreateUserParams{
			UserID:       "1",
			Username:     "username1",
			Email:        "username1@mail.me",
			PasswordHash: "pass",
		})
	})
	assert.NoError(t, err)
}

func TestExecTx_Invalid(t *testing.T) {
	clearTables(t, testStore.db)

	err := testStore.ExecTx(context.Background(), func(s Store) error {
		return s.CreateUser(context.Background(), sqlc.CreateUserParams{
			UserID:       "1",
			Username:     "",
			Email:        "",
			PasswordHash: "",
		})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction failed")
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())

	user, err := testStore.GetUserByID(context.Background(), "1")
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestExecTx_Rollback(t *testing.T) {
	clearTables(t, testStore.db)

	err := testStore.ExecTx(context.Background(), func(s Store) error {
		e := s.CreateUser(context.Background(), sqlc.CreateUserParams{
			UserID:       "1",
			Username:     "username1",
			Email:        "username1@mail.me",
			PasswordHash: "pass",
		})
		assert.NoError(t, e)

		ee := s.CreateUser(context.Background(), sqlc.CreateUserParams{})
		assert.Error(t, ee)

		return ee
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction failed")
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())

	user, err := testStore.GetUserByID(context.Background(), "1")
	assert.Nil(t, user)
	assert.Error(t, err)
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

	cases := []struct {
		tc    sqlc.CreateUserParams
		label string
	}{
		{
			tc: sqlc.CreateUserParams{
				UserID:       "",
				Username:     "username1",
				Email:        "username1@mail.me",
				PasswordHash: "pass",
			},
			label: "Missing UserID",
		},
		{
			tc: sqlc.CreateUserParams{
				UserID:       "1",
				Username:     "",
				Email:        "username1@mail.me",
				PasswordHash: "pass",
			},
			label: "Missing Username",
		},
		{
			tc: sqlc.CreateUserParams{
				UserID:       "1",
				Username:     "username1",
				Email:        "",
				PasswordHash: "pass",
			},
			label: "Missing Email",
		},
		{
			tc: sqlc.CreateUserParams{
				UserID:       "1",
				Username:     "username1",
				Email:        "username1@mail.me",
				PasswordHash: "",
			},
			label: "Missing Password",
		},
	}

	for _, tc := range cases {
		t.Run("Invalid Input "+tc.label, func(t *testing.T) {
			err = testStore.CreateUser(context.Background(), tc.tc)
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

func TestUpdateUser_Invalid(t *testing.T) {
	clearTables(t, testStore.db)
	user := insertUserHelper(t)

	var err error

	err = testStore.UpdateUser(context.Background(), sqlc.UpdateUserParams{UserID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())

	err = testStore.UpdateUser(context.Background(), sqlc.UpdateUserParams{UserID: user.UserID})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetUserByID_Success(t *testing.T) {
	clearTables(t, testStore.db)
	user := insertUserHelper(t)

	var err error
	var res *sqlc.User

	res, err = testStore.GetUserByID(context.Background(), user.UserID)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, res.UserID)
	assert.Equal(t, user.Username, res.Username)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.PasswordHash, res.PasswordHash)
}

func TestGetUserByID_Invalid(t *testing.T) {
	clearTables(t, testStore.db)
	_ = insertUserHelper(t)

	var err error

	_, err = testStore.GetUserByID(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetUserByID_NotFound(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	_, err = testStore.GetUserByID(context.Background(), "1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), sql.ErrNoRows.Error())
}

func TestGetUserByUsername_Success(t *testing.T) {
	clearTables(t, testStore.db)
	user := insertUserHelper(t)

	var err error
	var res *sqlc.User

	res, err = testStore.GetUserByUsername(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, res.UserID)
	assert.Equal(t, user.Username, res.Username)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.PasswordHash, res.PasswordHash)
}

func TestGetUserByUsername_Invalid(t *testing.T) {
	clearTables(t, testStore.db)
	_ = insertUserHelper(t)

	var err error

	_, err = testStore.GetUserByUsername(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetByUsername_NotFound(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	_, err = testStore.GetUserByUsername(context.Background(), "username")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), sql.ErrNoRows.Error())
}

func TestGetUserByEmail_Success(t *testing.T) {
	clearTables(t, testStore.db)
	user := insertUserHelper(t)

	var err error
	var res *sqlc.User

	res, err = testStore.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, res.UserID)
	assert.Equal(t, user.Username, res.Username)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.PasswordHash, res.PasswordHash)
}

func TestGetUserByEmail_Invalid(t *testing.T) {
	clearTables(t, testStore.db)
	_ = insertUserHelper(t)

	var err error

	_, err = testStore.GetUserByEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	clearTables(t, testStore.db)

	var err error

	_, err = testStore.GetUserByEmail(context.Background(), "username@mail.me")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), sql.ErrNoRows.Error())
}

func TestCreateToken_Success(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var token *sqlc.Token

	err = testStore.CreateToken(context.Background(), sqlc.CreateTokenParams{
		Jti:       "jti",
		UserID:    "1",
		Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
		TokenHash: "hash",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	assert.NoError(t, err)

	token, err = testStore.GetTokenByJTI(context.Background(), "jti")
	assert.NoError(t, err)
	assert.Equal(t, token.Jti, "jti")
	assert.Equal(t, token.UserID, "1")
	assert.Equal(t, token.Kind, pb.TokenKind_TOKEN_KIND_REFRESH.String())
	assert.Equal(t, token.TokenHash, "hash")
	assert.WithinDuration(t, token.IssuedAt, time.Now(), time.Second*5)
	assert.WithinDuration(t, token.ExpiresAt, time.Now().Add(time.Hour*24), time.Second*5)
}

func TestCreateToken_Invalid(t *testing.T) {
	clearTables(t, testStore.db)

	cases := []struct {
		tc    sqlc.CreateTokenParams
		label string
	}{
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "",
				UserID:    "1",
				Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
				TokenHash: "hash",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			label: "Missing JTI",
		},
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "jti",
				UserID:    "",
				Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
				TokenHash: "hash",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			label: "Missing UserID",
		},
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "jti",
				UserID:    "1",
				Kind:      "",
				TokenHash: "hash",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			label: "Missing TokenKind",
		},
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "jti",
				UserID:    "1",
				Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
				TokenHash: "",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			label: "Missing TokenHash",
		},
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "jti",
				UserID:    "1",
				Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
				TokenHash: "hash",
				IssuedAt:  time.Now().Add(time.Hour),
				ExpiresAt: time.Now().Add(time.Hour * 24),
			},
			label: "Invalid IssuedAt",
		},
		{
			tc: sqlc.CreateTokenParams{
				Jti:       "jti",
				UserID:    "1",
				Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
				TokenHash: "hash",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now(),
			},
			label: "Invalid ExpiresAt",
		},
	}

	for _, tc := range cases {
		t.Run("Invalid Input "+tc.label, func(t *testing.T) {
			err := testStore.CreateToken(context.Background(), tc.tc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), ErrInvalidInput.Error())
		})
	}
}

func TestGetTokenByJTI_Success(t *testing.T) {
	clearTables(t, testStore.db)
	params := insertTokenHelper(t)

	var err error
	var token *sqlc.Token

	token, err = testStore.GetTokenByJTI(context.Background(), params.Jti)
	assert.NoError(t, err)
	assert.Equal(t, params.Jti, token.Jti)
	assert.Equal(t, params.UserID, token.UserID)
	assert.Equal(t, params.Kind, token.Kind)
	assert.Equal(t, params.TokenHash, token.TokenHash)
	assert.True(t, token.IssuedAt.Before(time.Now()))
	assert.True(t, token.ExpiresAt.After(time.Now()))
}

func TestGetTokenByJTI_Invalid(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var token *sqlc.Token

	token, err = testStore.GetTokenByJTI(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetTokenByJTI_NotFound(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var token *sqlc.Token

	token, err = testStore.GetTokenByJTI(context.Background(), "does-not-exist")
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), sql.ErrNoRows.Error())
}

func TestGetTokensForUser_Success(t *testing.T) {
	clearTables(t, testStore.db)

	params := []sqlc.CreateTokenParams{
		{
			Jti:       "jti",
			UserID:    "1",
			Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
			TokenHash: "hash",
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		},
		{
			Jti:       "jti2",
			UserID:    "1",
			Kind:      pb.TokenKind_TOKEN_KIND_REFRESH.String(),
			TokenHash: "hash2",
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		},
	}

	var err error
	var tokens []*sqlc.Token

	err = testStore.CreateToken(context.Background(), params[0])
	assert.NoError(t, err)

	err = testStore.CreateToken(context.Background(), params[1])
	assert.NoError(t, err)

	tokens, err = testStore.GetTokensForUser(context.Background(), "1")
	assert.NoError(t, err)
	assert.Len(t, tokens, len(params))
}

func TestGetTokensForUser_Invalid(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var tokens []*sqlc.Token

	tokens, err = testStore.GetTokensForUser(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}

func TestGetTokensForUser_NotFound(t *testing.T) {
	clearTables(t, testStore.db)

	var err error
	var tokens []*sqlc.Token

	tokens, err = testStore.GetTokensForUser(context.Background(), "does-not-exist")
	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), sql.ErrNoRows.Error())
}
