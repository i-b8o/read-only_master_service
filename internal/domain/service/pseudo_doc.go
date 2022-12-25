package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"
)

type PseudoDocStorage interface {
	Exist(ctx context.Context, pseudoID string) (bool, error)
	CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error
	DeleteRelationship(ctx context.Context, docID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type pseudoDocService struct {
	storage PseudoDocStorage
}

func NewPseudoDocService(storage PseudoDocStorage) *pseudoDocService {
	return &pseudoDocService{storage: storage}
}

func (prs pseudoDocService) CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error {
	return prs.storage.CreateRelationship(ctx, pseudoDoc)
}

func (prs pseudoDocService) DeleteRelationship(ctx context.Context, docID uint64) error {
	return prs.storage.DeleteRelationship(ctx, docID)
}

func (prs pseudoDocService) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	return prs.storage.GetIDByPseudo(ctx, pseudoId)
}

func (prs pseudoDocService) Exist(ctx context.Context, pseudoID string) (bool, error) {
	return prs.storage.Exist(ctx, pseudoID)
}
