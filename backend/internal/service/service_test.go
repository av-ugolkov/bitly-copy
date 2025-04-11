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
	t.Run("different url - different users - same aliase", func(t *testing.T) {
		var (
			uid1     = "1"
			uid2     = "2"
			aliase   = "my_link"
			longURL1 = "https://www.google.com/"
			longURL2 = "https://www.mail.ru/"
			hashURL1 = hashURL(longURL1)
			hashURL2 = hashURL(longURL2)
		)

		db := NewMockdb(t)
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid1, hashURL1)).Return("", ErrUrlExists).Once()
		db.EXPECT().Get(ctx, fmt.Sprintf("%s:%s", uid2, hashURL2)).Return("", ErrUrlExists).Once()
		db.EXPECT().SetNX(ctx, mock.Anything, longURL1).Return(true, nil).Once()
		db.EXPECT().SetNX(ctx, mock.Anything, longURL2).Return(true, nil).Once()
		db.EXPECT().SetNX(ctx, mock.Anything, mock.Anything).Return(true, nil)
		svc := New(db)
		url1, err := svc.Add(ctx, uid1, longURL1, aliase)
		if err != nil {
			t.Error(err)
		}
		assert.NotEmpty(t, url1)

		url2, err := svc.Add(ctx, uid2, longURL2, aliase)
		assert.Error(t, err)
		assert.Empty(t, url2)
	})
}
