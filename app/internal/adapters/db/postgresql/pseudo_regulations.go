package postgressql

import (
	"context"
	"fmt"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type pseudoRegulationStorage struct {
	client client.PostgreSQLClient
}

func NewPseudoRegulation(client client.PostgreSQLClient) *pseudoRegulationStorage {
	return &pseudoRegulationStorage{client: client}
}

func (ps *pseudoRegulationStorage) Create(ctx context.Context, id uint64, pseudoId string) error {
	const sql = `INSERT INTO pseudo_regulations ("r_id", "pseudo") VALUES ($1, $2)`
	if _, err := ps.client.Exec(ctx, sql, id, pseudoId); err != nil {
		return err
	}
	return nil
}

func (ps *pseudoRegulationStorage) Delete(ctx context.Context, regulationID uint64) error {
	sql := `delete from pseudo_regulations where r_id=$1`
	_, err := ps.client.Exec(ctx, sql, regulationID)
	return err
}

func (ps *pseudoRegulationStorage) GetID(ctx context.Context, pseudoId string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT r_id FROM "pseudo_regulations" WHERE pseudo = '%s'`, pseudoId)
	rows, err := ps.client.Query(ctx, sql)
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
