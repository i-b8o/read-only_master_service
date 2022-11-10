package service

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
)

type SpeechStorage interface {
	// GetAllById(ctx context.Context, paragraphID uint64) ([]entity.Speech, error)
	Create(ctx context.Context, speech entity.Speech) (string, error)
	DeleteForParagraph(ctx context.Context, paragraphID uint64) error
}

type speechService struct {
	storage SpeechStorage
}

func NewSpeechService(storage SpeechStorage) *speechService {
	return &speechService{storage: storage}
}

// func (s speechService) GetAllById(ctx context.Context, paragraphID uint64) ([]entity.Speech, error) {
// 	return s.storage.GetAllById(ctx, paragraphID)
// }

func (s speechService) Create(ctx context.Context, speech entity.Speech) (string, error) {
	return s.storage.Create(ctx, speech)
}

func (s speechService) DeleteForParagraph(ctx context.Context, paragraphID uint64) error {
	return s.storage.DeleteForParagraph(ctx, paragraphID)
}
