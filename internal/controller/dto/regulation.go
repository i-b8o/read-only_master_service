package controller_dto

import (
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

func RegulationFromCreateRegulationRequset(req *pb.CreateRegulationRequest) entity.Regulation {
	return entity.Regulation{
		Name:         req.RegulationName,
		Pseudo:       req.PseudoId,
		Abbreviation: req.Abbreviation,
		Title:        req.Title,
	}
}
