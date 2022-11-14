package controller_dto

import (
	"regulations_supreme_service/internal/domain/entity"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
)

func RegulationFromCreateRegulationRequset(req *pb.CreateRegulationRequest) entity.Regulation {
	return entity.Regulation{
		Name:         req.RegulationName,
		Pseudo:       req.PseudoId,
		Abbreviation: req.Abbreviation,
		Title:        req.Title,
	}
}
