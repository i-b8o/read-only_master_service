package grpc_controller

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"

	controller_dto "regulations_supreme_service/internal/controller/grpc/dto"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
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

type SupremeRegulationGRPCService struct {
	regulationUsecase RegulationUsecase
	chapterUsecase    ChapterUsecase
	paragraphUsecase  ParagraphUsecase
	pb.UnimplementedSupremeRegulationGRPCServer
}

func NewSupremeRegulationGRPCService(regulationUsecase RegulationUsecase, chapterUsecase ChapterUsecase, paragraphUsecase ParagraphUsecase) *SupremeRegulationGRPCService {
	return &SupremeRegulationGRPCService{
		regulationUsecase: regulationUsecase,
		chapterUsecase:    chapterUsecase,
		paragraphUsecase:  paragraphUsecase,
	}
}

func (s *SupremeRegulationGRPCService) CreateRegulation(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	regulation := controller_dto.RegulationFromCreateRegulationRequset(req)
	ID, err := s.regulationUsecase.CreateRegulation(ctx, regulation)
	return &pb.CreateRegulationResponse{ID: ID}, err
}

func (s *SupremeRegulationGRPCService) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	chapter := controller_dto.ChapterFromCreateChapterRequest(req)
	id, err := s.chapterUsecase.CreateChapter(ctx, chapter)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: id}, nil
}

func (s *SupremeRegulationGRPCService) CreateParagraphs(ctx context.Context, req *pb.CreateParagraphsRequest) (*pb.Empty, error) {
	paragraphs := controller_dto.ParagraphsFromCreateParagraphsRequest(req)
	err := s.paragraphUsecase.CreateParagraphs(ctx, paragraphs)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *SupremeRegulationGRPCService) GenerateLinks(ctx context.Context, req *pb.GenerateLinksRequest) (*pb.GenerateLinksResponse, error) {
	ID := req.GetID()
	err := s.regulationUsecase.GenerateLinks(ctx, ID)
	return &pb.GenerateLinksResponse{ID: ID}, err
}

func (s *SupremeRegulationGRPCService) DeleteRegulation(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.regulationUsecase.DeleteRegulation(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, err
}
