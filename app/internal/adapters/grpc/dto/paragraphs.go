package dto

import (
	"regulations_supreme_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"
)

func CreateAllParagraphsRequestFromParagraphs(paragraphs []entity.Paragraph) (req wr_pb.CreateAllParagraphsRequest) {
	var wrParagraphs []*wr_pb.WritableParagraph
	for _, paragraph := range paragraphs {
		p := &wr_pb.WritableParagraph{ID: paragraph.ID, Num: paragraph.Num, HasLinks: paragraph.HasLinks, IsTable: paragraph.IsTable, IsNFT: paragraph.IsNFT, Class: paragraph.Class, Content: paragraph.Content}
		wrParagraphs = append(wrParagraphs, p)
	}
	req.Paragraphs = wrParagraphs
	return req
}
