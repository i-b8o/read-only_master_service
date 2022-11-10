package grpc_controller

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"

	pb "github.com/i-b8o/regulations_contracts/pb/supreme/v1"
)

type RegulationUsecase interface {
	CreateRegulation(ctx context.Context, regulation entity.Regulation) string
	GenerateLinks(ctx context.Context, regulationID uint64) error
}

type ChapterUsecase interface {
	CreateChapter(ctx context.Context, chapter entity.Chapter) string
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

func (s *SupremeRegulationGRPCService) GenerateLinks(ctx context.Context, req *pb.GenerateLinksRequest) (*pb.GenerateLinksResponse, error) {
	ID := req.GetID()
	err := s.regulationUsecase.GenerateLinks(ctx, ID)
	return &pb.GenerateLinksResponse{ID: ID}, err
}

func (s *SupremeRegulationGRPCService) CreateRegulation(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	// MAPPING
	regulation := entity.Regulation{
		Name:         req.RegulationName,
		Pseudo:       req.PseudoId,
		Abbreviation: req.Abbreviation,
		Title:        req.Title,
	}
	// Usecase
	id := s.regulationUsecase.CreateRegulation(ctx, regulation)
	return &pb.CreateRegulationResponse{Id: id}, nil
}

func (s *SupremeRegulationGRPCService) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	// MAPPING
	chapter := entity.Chapter{
		ID:           req.ChapterId,
		Pseudo:       req.PseudoId,
		Name:         req.ChapterName,
		Num:          req.ChapterNum,
		RegulationID: req.RegulationId,
		OrderNum:     req.OrderNum,
	}
	// Usecase
	id := s.chapterUsecase.CreateChapter(ctx, chapter)
	return &pb.CreateChapterResponse{ID: id}, nil
}
func (s *SupremeRegulationGRPCService) DeleteRegulation(ctx context.Context, regulationID uint64) error {
	return nil
}

func (s *SupremeRegulationGRPCService) CreateParagraphs(ctx context.Context, req *pb.CreateParagraphsRequest) (*pb.CreateParagraphsResponse, error) {
	var paragraphs []entity.Paragraph
	// MAPPING
	for _, p := range req.Paragraphs {
		paragraph := entity.Paragraph{
			ID:        p.ParagraphId,
			Num:       p.ParagraphOrderNum,
			IsTable:   p.IsTable,
			IsNFT:     p.IsNFT,
			HasLinks:  p.HasLinks,
			Class:     p.ParagraphClass,
			Content:   p.ParagraphText,
			ChapterID: p.ChapterId,
		}

		if p.ParagraphId > 0 {
			paragraph.ID = p.ParagraphId
		}
		paragraphs = append(paragraphs, paragraph)
	}
	// Usecase
	err := s.paragraphUsecase.CreateParagraphs(ctx, paragraphs)
	if err != nil {
		return &pb.CreateParagraphsResponse{Status: "not"}, err
	}
	return &pb.CreateParagraphsResponse{Status: "ok"}, nil
}
