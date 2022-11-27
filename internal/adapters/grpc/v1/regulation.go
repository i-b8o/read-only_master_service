package grpc_adapter

import (
	"context"
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type regulationStorage struct {
	client wr_pb.WriterGRPCClient
}

func NewRegulationStorage(client wr_pb.WriterGRPCClient) *regulationStorage {
	return &regulationStorage{client: client}
}

func (rs *regulationStorage) Create(ctx context.Context, regulation entity.Regulation) (uint64, error) {
	// Mapping
	req := &wr_pb.CreateRegulationRequest{Name: regulation.Name, Abbreviation: regulation.Abbreviation, Title: regulation.Title}
	resp, err := rs.client.CreateRegulation(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, err
}

func (rs *regulationStorage) Delete(ctx context.Context, regulationID uint64) error {
	req := &wr_pb.DeleteRegulationRequest{ID: regulationID}
	_, err := rs.client.DeleteRegulation(ctx, req)
	return err
}
