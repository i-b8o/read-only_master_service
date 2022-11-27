package controller

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	controller_dto "read-only_master_service/internal/controller/dto"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

type RegulationUsecase interface {
	CreateRegulation(ctx context.Context, regulation entity.Regulation) (uint64, error)
	GenerateLinks(ctx context.Context, regulationID uint64) error
	DeleteRegulation(ctx context.Context, ID uint64) error
}

type ChapterUsecase interface {
	CreateChapter(ctx context.Context, chapter entity.Chapter) (uint64, error)
}

type ParagraphUsecase interface {
	CreateParagraphs(ctx context.Context, paragraphs []entity.Paragraph) error
}

type MasterGRPCService struct {
	regulationUsecase RegulationUsecase
	chapterUsecase    ChapterUsecase
	paragraphUsecase  ParagraphUsecase
	pb.UnimplementedMasterGRPCServer
}

func NewMasterGRPCService(regulationUsecase RegulationUsecase, chapterUsecase ChapterUsecase, paragraphUsecase ParagraphUsecase) *MasterGRPCService {
	return &MasterGRPCService{
		regulationUsecase: regulationUsecase,
		chapterUsecase:    chapterUsecase,
		paragraphUsecase:  paragraphUsecase,
	}
}

func (s *MasterGRPCService) CreateRegulation(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	regulation := controller_dto.RegulationFromCreateRegulationRequset(req)
	// create a regulation and an id-pseudoId relationship
	ID, err := s.regulationUsecase.CreateRegulation(ctx, regulation)
	return &pb.CreateRegulationResponse{ID: ID}, err
}

func (s *MasterGRPCService) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	chapter := controller_dto.ChapterFromCreateChapterRequest(req)
	// create a chapter, create a link for the chapter and create an id-pseudoId relationship
	id, err := s.chapterUsecase.CreateChapter(ctx, chapter)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: id}, nil
}

func (s *MasterGRPCService) CreateParagraphs(ctx context.Context, req *pb.CreateParagraphsRequest) (*pb.Empty, error) {
	paragraphs := controller_dto.ParagraphsFromCreateParagraphsRequest(req)
	// cretae paragraphs, create links and speechs for paragraphs
	err := s.paragraphUsecase.CreateParagraphs(ctx, paragraphs)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *MasterGRPCService) GenerateLinks(ctx context.Context, req *pb.GenerateLinksRequest) (*pb.GenerateLinksResponse, error) {
	ID := req.GetID()
	err := s.regulationUsecase.GenerateLinks(ctx, ID)
	return &pb.GenerateLinksResponse{ID: ID}, err
}

func (s *MasterGRPCService) DeleteRegulation(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.regulationUsecase.DeleteRegulation(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, err
}
