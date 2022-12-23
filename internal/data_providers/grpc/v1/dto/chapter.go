package dto

import (
	"fmt"
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

func CreateChapterRequestFromChapter(chapter entity.Chapter) *wr_pb.CreateChapterRequest {
	var paragraphs []*wr_pb.WriterParagraph
	for _, p := range chapter.Paragraphs {
		paragraph := &wr_pb.WriterParagraph{ID: p.ID, Num: p.Num, HasLinks: p.HasLinks, IsTable: p.IsTable, IsNFT: p.IsNFT, Class: p.Class, ChapterID: p.ChapterID, Content: p.Content}
		paragraphs = append(paragraphs, paragraph)
	}
	fmt.Println("chapter: ", chapter.Title, chapter.Description, chapter.Keywords)
	return &wr_pb.CreateChapterRequest{ID: chapter.ID, Name: chapter.Name, Num: chapter.Num, DocID: chapter.DocID, OrderNum: chapter.OrderNum, Paragraphs: paragraphs, Title: chapter.Title, Description: chapter.Description, Keywords: chapter.Keywords}
}

// func GetChaptersIDsFromGetAllResponse(wrChapters []*wr_pb.) (IDs []uint64) {
// 	for _, wrChapter := range wrChapters {
// 		IDs = append(IDs, wrChapter.ID)
// 	}
// 	return IDs
// }
