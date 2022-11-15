package postgressql

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	client "read-only_master_service/pkg/client/postgresql"
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

func (prs *pseudoRegulationStorage) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT r_id FROM "pseudo_regulation" WHERE pseudo = '%s'`, pseudoId)
	rows, err := prs.client.Query(ctx, sql)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	regulationID := uint64(0)
	for rows.Next() {
		if err = rows.Scan(
			&regulationID,
		); err != nil {
			return 0, err
		}

	}
	return regulationID, err
}
