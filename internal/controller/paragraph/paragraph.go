package paragraph_controller

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	controller_dto "read-only_master_service/internal/controller/dto"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

type ParagraphUsecase interface {
	CreateParagraphs(ctx context.Context, paragraphs []entity.Paragraph) error
	UpdateOne(ctx context.Context, id uint64, content string) error
}

type ParagraphGrpcController struct {
	paragraphUsecase ParagraphUsecase
	pb.UnimplementedMasterParagraphGRPCServer
}

func NewParagraphGrpcController(paragraphUsecase ParagraphUsecase) *ParagraphGrpcController {
	return &ParagraphGrpcController{
		paragraphUsecase: paragraphUsecase,
	}
}
func (s *ParagraphGrpcController) Update(ctx context.Context, req *pb.UpdateParagraphRequest) (*pb.Empty, error) {
	ID := req.GetID()
	content := req.Content
	err := s.paragraphUsecase.UpdateOne(ctx, ID, content)
	return &pb.Empty{}, err

}

func (s *ParagraphGrpcController) Create(ctx context.Context, req *pb.CreateParagraphsRequest) (*pb.Empty, error) {
	paragraphs := controller_dto.ParagraphsFromCreateParagraphsRequest(req)
	// cretae paragraphs, create links and speechs for paragraphs
	err := s.paragraphUsecase.CreateParagraphs(ctx, paragraphs)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}