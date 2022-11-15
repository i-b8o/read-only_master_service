package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type RegulationStorage interface {
	Create(ctx context.Context, regulation entity.Regulation) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type regulationService struct {
	storage RegulationStorage
}

func NewRegulationService(storage RegulationStorage) *regulationService {
	return &regulationService{storage: storage}
}

func (s *regulationService) Create(ctx context.Context, regulation entity.Regulation) (uint64, error) {
	return s.storage.Create(ctx, regulation)
}

func (s *regulationService) Delete(ctx context.Context, regulationId uint64) error {
	return s.storage.Delete(ctx, regulationId)
}
