package controller_dto

import (
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
)

func ParagraphsFromCreateParagraphsRequest(req *pb.CreateParagraphsRequest) (paragraphs []entity.Paragraph) {
	for _, p := range req.Paragraphs {
		paragraph := entity.Paragraph{
			ID:        p.ParagraphId,
			Num:       p.ParagraphOrderNum,
			IsTable:   p.IsTable,
			IsNFT:     p.IsNFT,
			HasLinks:  p.HasLinks,
			Class:     p.ParagraphClass,
			Content:   p.ParagraphText,
			ChapterID: p.ChapterId,
		}

		if p.ParagraphId > 0 {
			paragraph.ID = p.ParagraphId
		}
		paragraphs = append(paragraphs, paragraph)
	}
	return paragraphs
}
