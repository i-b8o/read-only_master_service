package dto

import (
	"read-only_master_service/internal/domain/entity"

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

func ParagraphsFromGetParagraphsWithHrefsResponse(resp *wr_pb.GetParagraphsWithHrefsResponse) (paragraphs []entity.Paragraph) {
	for _, writableParagraph := range resp.Paragraphs {
		paragraph := entity.Paragraph{ID: writableParagraph.ID, Num: writableParagraph.Num, HasLinks: writableParagraph.HasLinks, IsTable: writableParagraph.IsTable, IsNFT: writableParagraph.IsNFT, Class: writableParagraph.Class, Content: writableParagraph.Content, ChapterID: writableParagraph.ChapterID}
		paragraphs = append(paragraphs, paragraph)
	}
	return paragraphs
}
