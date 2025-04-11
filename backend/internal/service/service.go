package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	urlTool "github.com/av-ugolkov/bitly-copy/pkg/url"
	"github.com/redis/go-redis/v9"
)

type (
	db interface {
		Get(ctx context.Context, key string) (string, error)
		GetAll(ctx context.Context, key string) ([]string, error)
		SetNX(ctx context.Context, key, value string) (bool, error)
		Increment(ctx context.Context, key string) error
		Remove(ctx context.Context, key string) error
	}
)

type Service struct {
	db db
}

func New(db db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Get(ctx context.Context, shortURL string) (string, error) {
	url, err := s.db.Get(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("service.Get - key [%s]: %w", shortURL, err)
	}

	go func() {
		ctx := context.Background()
		err = s.db.Increment(ctx, fmt.Sprintf("%s:%s", "statistics", shortURL))
		if err != nil {
			slog.Error("service.Get - key [%s]: %w", shortURL, err)
		}
	}()

	return url, nil
}

func (s *Service) Add(ctx context.Context, uid, url, aliase string) (string, error) {
	hashURL := urlTool.GetHashURL(url)
	shortUrl, err := s.db.Get(ctx, fmt.Sprintf("%s:%s", uid, hashURL))
	if err == nil {
		slog.Warn("service.Add - key [%s]: %w", url, ErrUrlExists)
		return shortUrl, nil
	}

	shortUrl = aliase
	if shortUrl == "" {
		shortUrl = urlTool.GenShortURL(url)
	}

	success, count := false, 0
	for !success && count < CountRetryGenShortUrl {
		success, err = s.db.SetNX(ctx, shortUrl, url)
		if !success && shortUrl != aliase {
			shortUrl = urlTool.GenShortURL(url)
			count++
		}
	}
	if err != nil {
		return "", fmt.Errorf("service.Add - key [%s]: %w", shortUrl, err)
	}

	success, err = s.db.SetNX(ctx, fmt.Sprintf("%s:%s", uid, hashURL), shortUrl)
	if err != nil {
		return "", fmt.Errorf("service.Add - key [%s]: %w", hashURL, err)
	}
	return shortUrl, nil
}

func (s *Service) GetUrls(ctx context.Context, uid string) ([]string, error) {
	urls, err := s.db.GetAll(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("service.GetUrls - key [%s]: %w", uid, err)
	}

	return urls, nil
}

func (s *Service) Remove(ctx context.Context, uid, shortUrl string) error {
	urlLink, err := url.Parse(shortUrl)
	if err == nil && urlLink.Host != "" {
		shortUrl = urlLink.Path[1:]
	}

	url, err := s.db.Get(ctx, shortUrl)
	if err != nil {
		return fmt.Errorf("service.Remove - key [%s]: %w", shortUrl, err)
	}

	hashURL := urlTool.GetHashURL(url)

	err = s.db.Remove(ctx, fmt.Sprintf("%s:%s", uid, hashURL))
	if err != nil {
		return fmt.Errorf("service.Remove - key [%s]: %w", hashURL, err)
	}

	err = s.db.Remove(ctx, shortUrl)
	if err != nil {
		return fmt.Errorf("service.Remove - key [%s]: %w", shortUrl, err)
	}

	err = s.db.Remove(ctx, fmt.Sprintf("%s:%s", "statistics", shortUrl))
	if err != nil {
		return fmt.Errorf("service.Remove - key [%s]: %w", shortUrl, err)
	}

	return nil
}

func (s *Service) Statistics(ctx context.Context, shortUrl string) (string, error) {
	urlLink, err := url.Parse(shortUrl)
	if err == nil && urlLink.Host != "" {
		shortUrl = urlLink.Path[1:]
	}

	value, err := s.db.Get(ctx, fmt.Sprintf("%s:%s", "statistics", shortUrl))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "0", nil
		}
		return "", fmt.Errorf("service.Statistics - key [%s]: %w", shortUrl, err)
	}

	return value, nil
}
