package service

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type PseudoDocStorage interface {
	Exist(ctx context.Context, pseudoID string) (bool, error)
	CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error
	DeleteRelationship(ctx context.Context, docID uint64) error
	GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error)
}

type pseudoDocService struct {
	storage PseudoDocStorage
	logging logging.Logger
}

func NewPseudoDocService(storage PseudoDocStorage, logging logging.Logger) *pseudoDocService {
	return &pseudoDocService{storage: storage, logging: logging}
}

func (prs pseudoDocService) CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error {
	err := prs.storage.CreateRelationship(ctx, pseudoDoc)
	if err != nil {
		prs.logging.Errorf("%v %v", pseudoDoc, err)
		return err
	}
	return nil
}

func (prs pseudoDocService) DeleteRelationship(ctx context.Context, docID uint64) error {
	err := prs.storage.DeleteRelationship(ctx, docID)
	if err != nil {
		prs.logging.Errorf("%d %v", docID, err)
		return err
	}
	return nil
}

func (prs pseudoDocService) GetIDByPseudo(ctx context.Context, pseudoId string) (*uint64, error) {
	id, err := prs.storage.GetIDByPseudo(ctx, pseudoId)
	if err != nil {
		prs.logging.Errorf("%s %v", pseudoId, err)
		return nil, err
	}
	return &id, nil
}

func (prs pseudoDocService) Exist(ctx context.Context, pseudoID string) (*bool, error) {
	exist, err := prs.storage.Exist(ctx, pseudoID)
	if err != nil {
		return nil, err
	}
	return &exist, nil
}
