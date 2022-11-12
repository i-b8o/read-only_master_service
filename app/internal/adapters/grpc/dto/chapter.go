package dto

import (
	"regulations_supreme_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"
)

func CreateChapterRequestFromChapter(chapter entity.Chapter) *wr_pb.CreateChapterRequest {
	var paragraphs []*wr_pb.WritableParagraph
	for _, p := range chapter.Paragraphs {
		paragraph := &wr_pb.WritableParagraph{ID: p.ID, Num: p.Num, HasLinks: p.HasLinks, IsTable: p.IsTable, IsNFT: p.IsNFT, Class: p.Class, ChapterID: p.ChapterID, Content: p.Content}
		paragraphs = append(paragraphs, paragraph)
	}
	return &wr_pb.CreateChapterRequest{ID: chapter.ID, Name: chapter.Name, Num: chapter.Num, RegulationID: chapter.RegulationID, OrderNum: chapter.OrderNum, Paragraphs: paragraphs}
}

func GetChaptersIDsFromGetAllResponse(wrChapters []*wr_pb.WritableChapter) (IDs []uint64) {
	for _, wrChapter := range wrChapters {
		IDs = append(IDs, wrChapter.ID)
	}
	return IDs
}
