package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type LinkStorage interface {
	GetAll(ctx context.Context) ([]*entity.Link, error)
	// GetAllByChapterID(ctx context.Context, chapterID uint64) ([]*entity.Link, error)
	Create(ctx context.Context, link entity.Link) error
	CreateForChapter(ctx context.Context, link entity.Link) error
	// GetOneByParagraphID(ctx context.Context, paragraphID, regregulationID uint64) (entity.Link, error)
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type linkService struct {
	storage LinkStorage
}

func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{storage: storage}
}

func (s linkService) GetAll(ctx context.Context) ([]*entity.Link, error) {
	return s.storage.GetAll(ctx)
}

// func (s linkService) GetAllByChapterID(ctx context.Context, chapterID uint64) ([]*entity.Link, error) {
// 	return s.storage.GetAllByChapterID(ctx, chapterID)
// }

func (s linkService) Create(ctx context.Context, link entity.Link) error {
	return s.storage.Create(ctx, link)
}

func (s linkService) CreateForChapter(ctx context.Context, link entity.Link) error {
	return s.storage.CreateForChapter(ctx, link)
}

// func (s linkService) GetOneByParagraphID(ctx context.Context, paragraphID, regregulationID uint64) (entity.Link, error) {
// 	return s.storage.GetOneByParagraphID(ctx, paragraphID, regregulationID)
// }

func (s linkService) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	return s.storage.DeleteForChapter(ctx, chapterID)
}
