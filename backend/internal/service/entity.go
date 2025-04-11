package service

import "fmt"

const (
	CountRetryGenShortUrl = 10
)

var (
	ErrUrlExists = fmt.Errorf("url already exists")
)
