package postgressql

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	client "read-only_master_service/pkg/client/postgresql"
)

type pseudoChapterStorage struct {
	client client.PostgreSQLClient
}

func NewPseudoChapterStorage(client client.PostgreSQLClient) *pseudoChapterStorage {
	return &pseudoChapterStorage{client: client}
}

func (pcs *pseudoChapterStorage) CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error {
	const sql = `INSERT INTO pseudo_chapter ("c_id", "pseudo") VALUES ($1, $2)`
	if _, err := pcs.client.Exec(ctx, sql, pseudoChapter.ID, pseudoChapter.PseudoId); err != nil {
		return err
	}
	return nil
}

func (pcs *pseudoChapterStorage) DeleteRelationship(ctx context.Context, chapterID uint64) error {
	sql := `delete from pseudo_chapter where c_id=$1`
	_, err := pcs.client.Exec(ctx, sql, chapterID)
	return err
}

func (pcs *pseudoChapterStorage) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT c_id FROM "pseudo_chapter" WHERE pseudo = '%s' LIMIT 1`, pseudoId)
	row := pcs.client.QueryRow(ctx, sql)
	chapterID := uint64(0)
	err := row.Scan(&chapterID)
	return chapterID, err
}
