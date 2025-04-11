package service

import (
	"context"
	"fmt"
	"testing"

	urlTool "github.com/av-ugolkov/bitly-copy/pkg/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func hashURL(url string) string {
	return urlTool.GetHashURL(url)
}

func TestService_Add(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	t.Run("add new url", func(t *testing.T) {
		var (
			uid     = "1"
			longURL = "https://www.google.com/"
			hashURL = hashURL(longURL)
		)

		db := NewMockdb(t)
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid, hashURL)).Return("", ErrUrlExists)
		db.EXPECT().SetNX(ctx, mock.Anything, longURL).Return(true, nil)
		db.EXPECT().SetNX(ctx, mock.Anything, mock.Anything).Return(true, nil)
		svc := New(db)
		url, err := svc.Add(ctx, uid, longURL, "")
		if err != nil {
			t.Error(err)
		}
		assert.NotEmpty(t, url)
	})
	t.Run("add duplicate url", func(t *testing.T) {
		var (
			uid     = "1"
			longURL = "https://www.google.com/"
			hashURL = hashURL(longURL)
		)

		db := NewMockdb(t)
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid, hashURL)).Return("", ErrUrlExists).Once()
		db.EXPECT().SetNX(ctx, mock.Anything, longURL).Return(true, nil)
		db.EXPECT().SetNX(ctx, mock.Anything, mock.Anything).Return(true, nil)
		svc := New(db)
		url, err := svc.Add(ctx, uid, longURL, "")
		if err != nil {
			t.Error(err)
		}

		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid, hashURL)).
			RunAndReturn(func(ctx context.Context, key string) (string, error) {
				return url, nil
			}).Once()
		newUrl, err := svc.Add(ctx, uid, longURL, "")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, url, newUrl)
	})
	t.Run("same url - different users", func(t *testing.T) {
		var (
			uid1    = "1"
			uid2    = "2"
			longURL = "https://www.google.com/"
			hashURL = hashURL(longURL)
		)

		db := NewMockdb(t)
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid1, hashURL)).Return("", ErrUrlExists).Once()
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid2, hashURL)).Return("", ErrUrlExists).Once()
		db.EXPECT().SetNX(ctx, mock.Anything, longURL).Return(true, nil)
		db.EXPECT().SetNX(ctx, mock.Anything, mock.Anything).Return(true, nil)
		svc := New(db)
		url1, err := svc.Add(ctx, uid1, longURL, "")
		if err != nil {
			t.Error(err)
		}
		url2, err := svc.Add(ctx, uid2, longURL, "")
		if err != nil {
			t.Error(err)
		}

		assert.NotEqual(t, url1, url2)
	})
}
