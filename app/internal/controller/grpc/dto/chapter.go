package controller_dto

import (
	"regulations_supreme_service/internal/domain/entity"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
)

func ChapterFromCreateChapterRequest(req *pb.CreateChapterRequest) entity.Chapter {
	return entity.Chapter{
		ID:           req.ChapterId,
		Pseudo:       req.PseudoId,
		Name:         req.ChapterName,
		Num:          req.ChapterNum,
		RegulationID: req.RegulationId,
		OrderNum:     req.OrderNum,
	}
}
