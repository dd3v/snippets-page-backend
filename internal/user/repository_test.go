// +build integration

package user

import (
	"context"
	"testing"
	"time"

	"github.com/dd3v/snippets.page.backend/internal/entity"
	"github.com/dd3v/snippets.page.backend/internal/test"
	"github.com/dd3v/snippets.page.backend/pkg/dbcontext"
	"github.com/stretchr/testify/assert"
)

var db *dbcontext.DB
var r Repository
var table = "users"

func TestRepositoryMain(t *testing.T) {
	db = test.Database(t)
	test.TruncateTable(t, db, table)
	r = NewRepository(db)
}
func TestRepositoryCreate(t *testing.T) {
	cases := []struct {
		name   string
		entity entity.User
		fail   bool
	}{
		{
			"success",
			entity.User{
				PasswordHash: "hash_100",
				Login:        "user_100",
				Email:        "user_100@mail.com",
				Website:      "user_100.com",
				Banned:       false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			false,
		},
		{
			"duplicate login index",
			entity.User{
				PasswordHash: "hash_200",
				Login:        "user_100",
				Email:        "user_200@mail.com",
				Website:      "user_200.com",
				Banned:       false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			true,
		},
		{
			"duplicate email index",
			entity.User{
				PasswordHash: "hash_300",
				Login:        "user_300",
				Email:        "user_100@mail.com",
				Website:      "user_300.com",
				Banned:       false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := r.Create(context.TODO(), tc.entity)
			assert.Equal(t, tc.fail, err != nil)
		})
	}
}

func TestRepositoryUpdate(t *testing.T) {
	err := r.Update(context.TODO(), entity.User{
		ID:           1,
		PasswordHash: "hash_100",
		Login:        "user_100",
		Email:        "user_100@mail.com",
		Website:      "new_user_website.com",
		Banned:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	assert.Nil(t, err)
	user, err := r.FindByID(context.TODO(), 1)
	assert.Nil(t, err)
	assert.Equal(t, "new_user_website.com", user.Website)
	assert.Equal(t, true, user.Banned)
}

func TestRepositoryCount(t *testing.T) {
	count, err := r.Count(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, true, count != 0)
}

func TestRepositoryFindByID(t *testing.T) {
	_, err := r.FindByID(context.TODO(), 1)
	assert.Nil(t, err)
}

func TestRepositoryDelete(t *testing.T) {
	err := r.Delete(context.TODO(), 1)
	assert.Nil(t, err)
	_, err = r.FindByID(context.TODO(), 1)
	assert.NotNil(t, err)
}
