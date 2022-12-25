package postgressql

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	client "read-only_master_service/pkg/client/sqlite"
)

type pseudoDocStorage struct {
	client client.SQLiteClient
}

func NewPseudoDocStorage(client client.SQLiteClient) *pseudoDocStorage {
	return &pseudoDocStorage{client: client}
}

func (prs *pseudoDocStorage) Exist(ctx context.Context, pseudoID string) (bool, error) {
	const sql = `SELECT pseudo FROM pseudo_doc WHERE pseudo='$1'";`

	rows, err := prs.client.Query(sql)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	var pseudo string
	for rows.Next() {
		if err = rows.Scan(
			&pseudo,
		); err != nil {
			return false, err
		}

	}
	return pseudo == pseudoID, err
}

func (prs *pseudoDocStorage) CreateRelationship(ctx context.Context, pseudoDoc entity.PseudoDoc) error {
	const sql = `INSERT INTO pseudo_doc ("doc_id", "pseudo") VALUES ($1, $2)`
	if _, err := prs.client.Exec(sql, pseudoDoc.ID, pseudoDoc.PseudoId); err != nil {
		return err
	}
	return nil
}

func (prs *pseudoDocStorage) DeleteRelationship(ctx context.Context, docID uint64) error {
	const sql = `delete from pseudo_doc where doc_id=$1`
	_, err := prs.client.Exec(sql, docID)
	return err
}

func (prs *pseudoDocStorage) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT doc_id FROM "pseudo_doc" WHERE pseudo = '%s'`, pseudoId)
	rows, err := prs.client.Query(sql)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	docID := uint64(0)
	for rows.Next() {
		if err = rows.Scan(
			&docID,
		); err != nil {
			return 0, err
		}

	}
	return docID, err
}
