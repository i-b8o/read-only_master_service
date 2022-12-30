package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type AbsentStorage interface {
	Done(ctx context.Context, pseudo string) error
	Create(ctx context.Context, absent entity.Absent) error
	GetAll(ctx context.Context) ([]*entity.Absent, error)
}

type absentService struct {
	storage AbsentStorage
	logging logging.Logger
}

func NewAbsentService(storage AbsentStorage, logging logging.Logger) *absentService {
	return &absentService{storage: storage, logging: logging}
}

func (s absentService) Done(ctx context.Context, pseudo string) error {
	err := s.storage.Done(ctx, pseudo)
	if err != nil {
		s.logging.Errorf("%s %v", pseudo, err)
		return err
	}
	return nil
}

func (s absentService) Create(ctx context.Context, absent entity.Absent) error {
	err := s.storage.Create(ctx, absent)
	if err != nil {
		s.logging.Errorf("%v %v", absent, err)
		return err
	}
	return nil
}

func (s absentService) GetAll(ctx context.Context) ([]*entity.Absent, error) {
	absents, err := s.storage.GetAll(ctx)
	if err != nil {
		s.logging.Errorf("%v", err)
		return nil, err
	}
	return absents, nil
}
