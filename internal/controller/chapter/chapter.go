package chapter_controller

import (
	"context"
	chapter_dto "read-only_master_service/internal/controller/chapter/dto"
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

type ChapterUsecase interface {
	CreateChapter(ctx context.Context, chapter entity.Chapter) (uint64, error)
}

type ChapterGrpcController struct {
	chapterUsecase ChapterUsecase
	pb.UnimplementedMasterChapterGRPCServer
}

func NewChapterGrpcController(chapterUsecase ChapterUsecase) *ChapterGrpcController {
	return &ChapterGrpcController{
		chapterUsecase: chapterUsecase,
	}
}

func (s *ChapterGrpcController) Create(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	chapter := chapter_dto.ChapterFromCreateChapterRequest(req)
	// create a chapter, create a link for the chapter and create an id-pseudoId relationship
	id, err := s.chapterUsecase.CreateChapter(ctx, chapter)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: id}, nil
}

// TODO GetAll
