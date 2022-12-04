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
		Title:        req.Title,
	}
}

func RegulationsFromRegulations(domainRegulations []entity.Regulation) (regulations []*pb.Regulation) {
	for _, r := range domainRegulations {
		regulation := pb.Regulation{ID: r.Id, RegulationName: r.Name, Abbreviation: r.Abbreviation, Title: r.Title}
		regulations = append(regulations, &regulation)
	}
	return regulations
}
