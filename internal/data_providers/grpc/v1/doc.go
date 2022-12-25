package grpc_adapter

import (
	"context"
	"read-only_master_service/internal/data_providers/grpc/v1/dto"
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type docStorage struct {
	client wr_pb.WriterDocGRPCClient
}

func NewDocStorage(client wr_pb.WriterDocGRPCClient) *docStorage {
	return &docStorage{client: client}
}

func (rs *docStorage) Create(ctx context.Context, doc entity.Doc) (uint64, error) {
	// Mapping
	req := &wr_pb.CreateDocRequest{Name: doc.Name, Title: doc.Title, Description: doc.Description, Keywords: doc.Keywords}
	resp, err := rs.client.Create(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, err
}

func (rs *docStorage) Delete(ctx context.Context, docID uint64) error {
	req := &wr_pb.DeleteDocRequest{ID: docID}
	_, err := rs.client.Delete(ctx, req)
	return err
}

func (rs *docStorage) GetAll(ctx context.Context) ([]entity.Doc, error) {
	resp, err := rs.client.GetAll(ctx, &wr_pb.Empty{})
	if err != nil {
		return nil, err
	}
	return dto.CreateDocsFromGetDocsResponse(resp), nil
}
