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

func (rs *regulationStorage) Create(ctx context.Context, regulation entity.Regulation) (uint64, error) {
	// Mapping
	req := &wr_pb.CreateRegulationRequest{Name: regulation.Name, Abbreviation: regulation.Abbreviation, Title: regulation.Title}
	resp, err := rs.client.CreateRegulation(ctx, req)
	return resp.ID, err
}

func (rs *regulationStorage) DeleteRegulation(ctx context.Context, regulationID uint64) (string, error) {
	req := &wr_pb.DeleteRegulationRequest{ID: regulationID}
	resp, err := rs.client.DeleteRegulation(ctx, req)

	return resp.Status, err
}
