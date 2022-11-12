package usecase_regulation

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type RegulationService interface {
	Create(ctx context.Context, regulation entity.Regulation) (uint64, error)
	Delete(ctx context.Context, regulationId uint64) error
}

type ChapterService interface {
	DeleteAll(ctx context.Context, ID uint64) error
	GetAll(ctx context.Context, ID uint64) ([]uint64, error)
}
type ParagraphService interface {
	DeleteForRegulation(ctx context.Context, chaptersIDs []uint64) error
}

type AbsentService interface {
	Done(ctx context.Context, pseudo string) error
}

type PseudoRegulationService interface {
	CreateRelationship(ctx context.Context, pseudoRegulation entity.PseudoRegulation) error
	DeleteRelationship(ctx context.Context, regulationID uint64) error
}

type regulationUsecase struct {
	regulationService       RegulationService
	chapterService          ChapterService
	paragraphService        ParagraphService
	absentService           AbsentService
	pseudoRegulationService PseudoRegulationService
	logging                 logging.Logger
}

func NewRegulationUsecase(regulationService RegulationService, chapterService ChapterService, paragraphService ParagraphService, absentService AbsentService, pseudoRegulationService PseudoRegulationService, logging logging.Logger) *regulationUsecase {
	return &regulationUsecase{regulationService: regulationService, chapterService: chapterService, paragraphService: paragraphService, absentService: absentService, pseudoRegulationService: pseudoRegulationService, logging: logging}
}

func (u regulationUsecase) CreateRegulation(ctx context.Context, regulation entity.Regulation) (uint64, error) {
	// create a regulation
	ID, err := u.regulationService.Create(ctx, regulation)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// create an id-pseudoId relationship
	err = u.pseudoRegulationService.CreateRelationship(ctx, entity.PseudoRegulation{ID: ID, PseudoId: regulation.Pseudo})
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	// mark the regulation as done
	err = u.absentService.Done(ctx, regulation.Pseudo)
	if err != nil {
		u.logging.Error(err)
		return 0, err
	}

	return ID, nil
}

func (u regulationUsecase) DeleteRegulation(ctx context.Context, ID uint64) error {
	// delete a regulation
	err := u.regulationService.Delete(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}

	// delete all paragraphs for the regulation
	IDs, err := u.chapterService.GetAll(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	err = u.paragraphService.DeleteForRegulation(ctx, IDs)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	// delete all chapters for the regulation
	err = u.chapterService.DeleteAll(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}

	// delete the id-pseudoId relationship
	err = u.pseudoRegulationService.DeleteRelationship(ctx, ID)
	if err != nil {
		u.logging.Error(err)
		return err
	}
	return nil
}
