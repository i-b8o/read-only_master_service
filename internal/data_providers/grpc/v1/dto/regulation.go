package dto

import (
	"read-only_master_service/internal/domain/entity"

	wr_pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

func CreateRegulationsFromGetRegulationsResponse(resp *wr_pb.GetRegulationsResponse) (regulations []entity.Regulation) {
	for _, r := range resp.Regulations {
		regulation := entity.Regulation{Id: r.ID, Name: r.Name, Abbreviation: r.Abbreviation, Title: &r.Title}
		regulations = append(regulations, regulation)
	}
	return regulations
}
