package doc_controller

import (
	"context"
	doc_dto "read-only_master_service/internal/controller/doc/dto"
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

type DocUsecase interface {
	GetAll(ctx context.Context) ([]entity.Doc, error)
	CreateDoc(ctx context.Context, doc entity.Doc) (uint64, error)
	GenerateLinks(ctx context.Context, docID uint64) error
	DeleteDoc(ctx context.Context, ID uint64) error
	GetAbsents(ctx context.Context) ([]*entity.Absent, error)
}

type DocGrpcController struct {
	docUsecase DocUsecase
	pb.UnimplementedMasterDocGRPCServer
}

func NewDocGrpcController(docUsecase DocUsecase) *DocGrpcController {
	return &DocGrpcController{
		docUsecase: docUsecase,
	}
}

func (s *DocGrpcController) Create(ctx context.Context, req *pb.CreateDocRequest) (*pb.CreateDocResponse, error) {
	doc := doc_dto.DocFromCreateDocRequest(req)
	// create a doc and an id-pseudoId relationship
	ID, err := s.docUsecase.CreateDoc(ctx, doc)
	return &pb.CreateDocResponse{ID: ID}, err
}

func (s *DocGrpcController) GetAll(ctx context.Context, req *pb.Empty) (*pb.GetAllDocsResponse, error) {
	domainDocs, err := s.docUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	docs := doc_dto.DocsFromDocs(domainDocs)
	return &pb.GetAllDocsResponse{Docs: docs}, nil
}

func (s *DocGrpcController) UpdateLinks(ctx context.Context, req *pb.UpdateLinksRequest) (*pb.UpdateLinksResponse, error) {
	ID := req.GetID()
	err := s.docUsecase.GenerateLinks(ctx, ID)
	return &pb.UpdateLinksResponse{ID: ID}, err
}

func (s *DocGrpcController) Delete(ctx context.Context, req *pb.DeleteDocRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.docUsecase.DeleteDoc(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, err
}

func (s *DocGrpcController) GetAbsents(ctx context.Context, req *pb.Empty) (*pb.GetAbsentsResponse, error) {
	absents, err := s.docUsecase.GetAbsents(ctx)
	if err != nil {
		return nil, err
	}
	return doc_dto.GetAbsentsResponseFromAbsents(absents), err
}
