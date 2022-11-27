package grpc_adapter

import (
	"context"
	"read-only_master_service/internal/adapters/grpc/v1/dto"
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type chapterStorage struct {
	client wr_pb.WriterGRPCClient
}

func NewChapterStorage(client wr_pb.WriterGRPCClient) *chapterStorage {
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

func (cs *chapterStorage) GetAllIds(ctx context.Context, ID uint64) ([]uint64, error) {
	req := &wr_pb.GetAllChaptersIdsRequest{ID: ID}
	resp, err := cs.client.GetAllChaptersIds(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.IDs, err
}

func (cs *chapterStorage) GetRegulationIdByChapterId(ctx context.Context, ID uint64) (uint64, error) {
	req := &wr_pb.GetRegulationIdByChapterIdRequest{ID: ID}
	resp, err := cs.client.GetRegulationIdByChapterId(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, err
}
