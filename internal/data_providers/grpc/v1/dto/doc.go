package dto

import (
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

func CreateDocsFromGetDocsResponse(resp *wr_pb.GetDocsResponse) (docs []entity.Doc) {
	for _, r := range resp.Docs {
		doc := entity.Doc{Id: r.ID, Name: r.Name, Title: r.Title, Description: r.Description, Keywords: r.Keywords}
		docs = append(docs, doc)
	}
	return docs
}
