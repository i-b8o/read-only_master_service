package grpc_adapter

import (
	"context"
	"regulations_supreme_service/internal/adapters/grpc/v1/dto"
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

func (ps *paragraphStorage) GetParagraphsWithHrefs(ctx context.Context, chapterId uint64) ([]entity.Paragraph, error) {
	req := &wr_pb.GetParagraphsWithHrefsRequest{ID: chapterId}
	resp, err := ps.client.GetParagraphsWithHrefs(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.ParagraphsFromGetParagraphsWithHrefsResponse(resp), nil
}

func (ps *paragraphStorage) UpdateOne(ctx context.Context, id uint64, content string) error {
	req := &wr_pb.UpdateOneParagraphRequest{ID: id, Content: content}
	_, err := ps.client.UpdateOneParagraph(ctx, req)
	return err
}
