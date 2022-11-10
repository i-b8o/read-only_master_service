package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type RegulationStorage interface {
	GetOne(ctx context.Context, regulationID uint64) (entity.Regulation, error)
	GetAll(ctx context.Context) ([]entity.Regulation, error)
	Create(ctx context.Context, regulation entity.Regulation) (string, error)
	DeleteRegulation(ctx context.Context, regulationID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type regulationService struct {
	storage RegulationStorage
}

func NewRegulationService(storage RegulationStorage) *regulationService {
	return &regulationService{storage: storage}
}

func (s *regulationService) Create(ctx context.Context, regulation entity.Regulation) (string, error) {
	return s.storage.Create(ctx, regulation)
}

func (s *regulationService) DeleteRegulation(ctx context.Context, regulationID uint64) error {
	return s.storage.DeleteRegulation(ctx, regulationID)
}

// TODO make it from pseudo_regulation
func (s *regulationService) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	return s.storage.GetIDByPseudo(ctx, pseudoId)
}

// func (s *regulationService) GetOne(ctx context.Context, regulationID uint64) (entity.Regulation, error) {
// 	return s.storage.GetOne(ctx, regulationID)
// }

// func (s *regulationService) GetAll(ctx context.Context) ([]entity.Regulation, error) {
// 	return s.storage.GetAll(ctx)
// }
