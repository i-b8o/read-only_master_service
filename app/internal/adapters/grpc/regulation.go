package grpc_adapter

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"
)

type regulationStorage struct {
	client wr_pb.WritableRegulationGRPCClient
}

func NewRegulationStorage(client wr_pb.WritableRegulationGRPCClient) *regulationStorage {
	return &regulationStorage{client: client}
}

func (rs *regulationStorage) GetOne(ctx context.Context, regulationID uint64) (entity.Regulation, error) {

}

func (rs *regulationStorage) GetAll(ctx context.Context) ([]entity.Regulation, error) {

}

func (rs *regulationStorage) Create(ctx context.Context, regulation entity.Regulation) (string, error) {
	// Mapping
	req := &wr_pb.CreateRegulationRequest{Name: regulation.Name, Abbreviation: regulation.Abbreviation, Title: regulation.Title}
	return := rs.client.CreateRegulation(ctx, req)
	
}

func (rs *regulationStorage) DeleteRegulation(ctx context.Context, regulationID uint64) error {

}

func (rs *regulationStorage) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {

}
