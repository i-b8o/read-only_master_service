package postgressql

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type pseudoRegulationStorage struct {
	client client.PostgreSQLClient
}

func NewPseudoRegulationStorage(client client.PostgreSQLClient) *pseudoRegulationStorage {
	return &pseudoRegulationStorage{client: client}
}

func (prs *pseudoRegulationStorage) CreateRelationship(ctx context.Context, pseudoRegulation entity.PseudoRegulation) error {
	const sql = `INSERT INTO pseudo_regulation ("r_id", "pseudo") VALUES ($1, $2)`
	if _, err := prs.client.Exec(ctx, sql, pseudoRegulation.ID, pseudoRegulation.PseudoId); err != nil {
		return err
	}
	return nil
}

func (prs *pseudoRegulationStorage) DeleteRelationship(ctx context.Context, regulationID uint64) error {
	sql := `delete from pseudo_regulations where r_id=$1`
	_, err := prs.client.Exec(ctx, sql, regulationID)
	return err
}
