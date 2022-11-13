package grpc_adapter

import (
	"context"
	"regulations_supreme_service/internal/adapters/grpc/v1/dto"
	"regulations_supreme_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"
)

type chapterStorage struct {
	client wr_pb.WritableRegulationGRPCClient
}

func NewChapterStorage(client wr_pb.WritableRegulationGRPCClient) *chapterStorage {
	return &chapterStorage{client: client}
}

func (cs *chapterStorage) Create(ctx context.Context, chapter entity.Chapter) (uint64, error) {
	req := dto.CreateChapterRequestFromChapter(chapter)
	resp, err := cs.client.CreateChapter(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, err
}

func (cs *chapterStorage) DeleteAll(ctx context.Context, ID uint64) error {
	req := &wr_pb.DeleteChaptersForRegulationRequest{ID: ID}
	_, err := cs.client.DeleteChaptersForRegulation(ctx, req)
	return err
}

func (cs *chapterStorage) GetAllIds(ctx context.Context, ID uint64) ([]uint64, error) {
	req := &wr_pb.GetAllChaptersRequest{ID: ID}
	resp, err := cs.client.GetAllChapters(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.GetChaptersIDsFromGetAllResponse(resp.Chapters), err
}

func (cs *chapterStorage) GetRegulationIdByChapterId(ctx context.Context, ID uint64) (uint64, error) {
	req := &wr_pb.GetRegulationIdByChapterIdRequest{ID: ID}
	resp, err := cs.client.GetRegulationIdByChapterId(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, err
}
