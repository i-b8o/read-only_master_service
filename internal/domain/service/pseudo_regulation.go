package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type PseudoRegulationStorage interface {
	CreateRelationship(ctx context.Context, pseudoRegulation entity.PseudoRegulation) error
	DeleteRelationship(ctx context.Context, regulationID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type pseudoRegulationService struct {
	storage PseudoRegulationStorage
}

func NewPseudoRegulationService(storage PseudoRegulationStorage) *pseudoRegulationService {
	return &pseudoRegulationService{storage: storage}
}

func (prs pseudoRegulationService) CreateRelationship(ctx context.Context, pseudoRegulation entity.PseudoRegulation) error {
	return prs.storage.CreateRelationship(ctx, pseudoRegulation)
}

func (prs pseudoRegulationService) DeleteRelationship(ctx context.Context, regulationID uint64) error {
	return prs.storage.DeleteRelationship(ctx, regulationID)
}

func (prs pseudoRegulationService) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	return prs.storage.GetIDByPseudo(ctx, pseudoId)
}
