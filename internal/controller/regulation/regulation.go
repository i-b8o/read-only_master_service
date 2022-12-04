package regulation_controller

import (
	"context"
	regulation_dto "read-only_master_service/internal/controller/regulation/dto"
	"read-only_master_service/internal/domain/entity"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
)

type RegulationUsecase interface {
	GetAll(ctx context.Context) ([]entity.Regulation, error)
	CreateRegulation(ctx context.Context, regulation entity.Regulation) (uint64, error)
	GenerateLinks(ctx context.Context, regulationID uint64) error
	DeleteRegulation(ctx context.Context, ID uint64) error
}

type RegulationGRPCService struct {
	regulationUsecase RegulationUsecase
	pb.UnimplementedMasterRegulationGRPCServer
}

func NewRegulationGRPCService(regulationUsecase RegulationUsecase) *RegulationGRPCService {
	return &RegulationGRPCService{
		regulationUsecase: regulationUsecase,
	}
}

func (s *RegulationGRPCService) Create(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	regulation := regulation_dto.RegulationFromCreateRegulationRequest(req)
	// create a regulation and an id-pseudoId relationship
	ID, err := s.regulationUsecase.CreateRegulation(ctx, regulation)
	return &pb.CreateRegulationResponse{ID: ID}, err
}

func (s *RegulationGRPCService) GetAll(ctx context.Context, req *pb.Empty) (*pb.GetAllRegulationsResponse, error) {
	domainRegulations, err := s.regulationUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	regulations := regulation_dto.RegulationsFromRegulations(domainRegulations)
	return &pb.GetAllRegulationsResponse{Regulations: regulations}, nil
}

func (s *RegulationGRPCService) UpdateLinks(ctx context.Context, req *pb.UpdateLinksRequest) (*pb.UpdateLinksResponse, error) {
	ID := req.GetID()
	err := s.regulationUsecase.GenerateLinks(ctx, ID)
	return &pb.UpdateLinksResponse{ID: ID}, err
}

func (s *RegulationGRPCService) Delete(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.regulationUsecase.DeleteRegulation(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, err
}

// TODO UpdateAbsent
