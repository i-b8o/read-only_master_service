package grpc_adapter

import (
	"context"
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
	var paragraphs []*wr_pb.WritableParagraph
	for _, p := range chapter.Paragraphs {
		paragraph := &wr_pb.WritableParagraph{ID: p.ID, Num: p.Num, HasLinks: p.HasLinks, IsTable: p.IsTable, IsNFT: p.IsNFT, Class: p.Class, ChapterID: p.ChapterID, Content: p.Content}
		paragraphs = append(paragraphs, paragraph)
	}
	req := &wr_pb.CreateChapterRequest{ID: chapter.ID, Name: chapter.Name, Num: chapter.Num, RegulationID: chapter.RegulationID, OrderNum: chapter.OrderNum, Paragraphs: paragraphs}
	resp, err := cs.client.CreateChapter(ctx, req)
	return resp.ID, err

}

func (cs *chapterStorage) GetAllById(ctx context.Context, regulationID uint64) ([]entity.Chapter, error) {
	req := &wr_pb.GetAllChaptersRequest{ID: regulationID}
	resp, err := cs.client.GetAllChapters(ctx, req)
	if err != nil {
		return nil, err
	}
	var chapters []entity.Chapter
	// Move to DTO
	for _, c := range resp.Chapters {
		var paragraphs []entity.Paragraph
		for _, p := range c.Paragraphs {
			paragraph := entity.Paragraph{ID: p.ID, Num: p.Num, HasLinks: p.HasLinks, IsTable: p.IsTable, IsNFT: p.IsNFT, Class: p.Class, Content: p.Content, ChapterID: p.ChapterID}
			paragraphs = append(paragraphs, paragraph)
		}
		chapter := &entity.Chapter{ID: c.ID, Name: c.Name, Num: c.Num, RegulationID: c.RegulationID, OrderNum: c.OrderNum, Paragraphs: paragraphs}
		chapters = append(chapters, *chapter)
	}

}

func (cs *chapterStorage) GetOrderNum(ctx context.Context, id uint64) (orderNum uint64, err error) {

}

func (cs *chapterStorage) GetOneById(ctx context.Context, chapterID uint64) (entity.Chapter, error) {

}

func (cs *chapterStorage) DeleteForRegulation(ctx context.Context, regulationID uint64) error {

}
