package regulation_dto

import (
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

func RegulationFromCreateRegulationRequest(req *pb.CreateRegulationRequest) entity.Regulation {
	return entity.Regulation{
		Name:         req.RegulationName,
		Pseudo:       req.PseudoId,
		Abbreviation: req.Abbreviation,
		Title:        &req.Title,
	}
}

func RegulationsFromRegulations(domainRegulations []entity.Regulation) (regulations []*pb.Regulation) {
	for _, r := range domainRegulations {
		regulation := pb.Regulation{ID: r.Id, RegulationName: r.Name, Abbreviation: r.Abbreviation, Title: *r.Title}
		regulations = append(regulations, &regulation)
	}
	return regulations
}

func GetAbsentsResponseFromAbsents(domainAbsents []*entity.Absent) (response *pb.GetAbsentsResponse) {
	var absents []*pb.MasterAbsent
	for _, a := range domainAbsents {
		absent := pb.MasterAbsent{ID: a.ID, Pseudo: a.Pseudo, Done: a.Done, ParagraphId: a.ParagraphID}
		absents = append(absents, &absent)
	}
	return &pb.GetAbsentsResponse{Absents: absents}
}
