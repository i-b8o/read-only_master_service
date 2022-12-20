package doc_dto

import (
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

func DocFromCreateDocRequest(req *pb.CreateDocRequest) entity.Doc {
	return entity.Doc{
		Name:         req.DocName,
		Pseudo:       req.PseudoId,
		Abbreviation: req.Abbreviation,
		Title:        &req.Title,
	}
}

func DocsFromDocs(domainDocs []entity.Doc) (docs []*pb.Doc) {
	for _, r := range domainDocs {
		doc := pb.Doc{ID: r.Id, DocName: r.Name, Abbreviation: r.Abbreviation, Title: *r.Title}
		docs = append(docs, &doc)
	}
	return docs
}

func GetAbsentsResponseFromAbsents(domainAbsents []*entity.Absent) (response *pb.GetAbsentsResponse) {
	var absents []*pb.MasterAbsent
	for _, a := range domainAbsents {
		absent := pb.MasterAbsent{ID: a.ID, Pseudo: a.Pseudo, Done: a.Done, ParagraphId: a.ParagraphID}
		absents = append(absents, &absent)
	}
	return &pb.GetAbsentsResponse{Absents: absents}
}
