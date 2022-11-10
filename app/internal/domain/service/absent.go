package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type AbsentStorage interface {
	Create(ctx context.Context, absent entity.Absent) error
	Done(ctx context.Context, pseudo string) error
}

type absentService struct {
	storage AbsentStorage
}

func NewAbsentService(storage AbsentStorage) *absentService {
	return &absentService{storage: storage}
}

func (s absentService) Create(ctx context.Context, absent entity.Absent) error {
	return s.storage.Create(ctx, absent)
}

func (s absentService) Done(ctx context.Context, pseudo string) error {
	return s.storage.Done(ctx, pseudo)
}
