package controller_dto

import (
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

func ChapterFromCreateChapterRequest(req *pb.CreateChapterRequest) entity.Chapter {
	return entity.Chapter{
		Pseudo:      req.PseudoId,
		Name:        req.ChapterName,
		Num:         req.ChapterNum,
		DocID:       req.DocId,
		Title:       req.Title,
		Description: req.Description,
		Keywords:    req.Keywords,
	}
}
