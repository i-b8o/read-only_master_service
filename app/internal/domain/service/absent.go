package service

import (
	"context"
)

type AbsentStorage interface {
	Done(ctx context.Context, pseudo string) error
}

type absentService struct {
	storage AbsentStorage
}

func NewAbsentService(storage AbsentStorage) *absentService {
	return &absentService{storage: storage}
}

func (s absentService) Done(ctx context.Context, pseudo string) error {
	return s.storage.Done(ctx, pseudo)
}
