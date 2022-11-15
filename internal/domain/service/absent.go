package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type AbsentStorage interface {
	Done(ctx context.Context, pseudo string) error
	Create(ctx context.Context, absent entity.Absent) error
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

func (s absentService) Create(ctx context.Context, absent entity.Absent) error {
	return s.storage.Create(ctx, absent)
}
