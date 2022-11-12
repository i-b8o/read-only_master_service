package grpc_adapter

import (
	"context"
	"regulations_supreme_service/internal/adapters/grpc/dto"
	"regulations_supreme_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"
)

type paragraphStorage struct {
	client wr_pb.WritableRegulationGRPCClient
}

func NewParagraphStorage(client wr_pb.WritableRegulationGRPCClient) *paragraphStorage {
	return &paragraphStorage{client: client}
}

// Delete
func (ps *paragraphStorage) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	_, err := ps.client.DeleteParagraphsForChapter(ctx, &wr_pb.DeleteParagraphsForChapterRequest{ID: chapterID})
	return err
}

func (ps *paragraphStorage) CreateAll(ctx context.Context, paragraphs []entity.Paragraph) error {
	req := dto.CreateAllParagraphsRequestFromParagraphs(paragraphs)
	_, err := ps.client.CreateAllParagraphs(ctx, &req)
	return err
}
