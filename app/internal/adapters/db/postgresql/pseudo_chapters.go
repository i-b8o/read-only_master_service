package postgressql

import (
	"context"
	"fmt"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type pseudoChapterStorage struct {
	client client.PostgreSQLClient
}

func NewPseudoChapter(client client.PostgreSQLClient) *pseudoChapterStorage {
	return &pseudoChapterStorage{client: client}
}

func (ps *pseudoChapterStorage) Create(ctx context.Context, ID uint64, pseudoId string) error {
	const sql = `INSERT INTO pseudo_chapters ("c_id", "pseudo") VALUES ($1, $2)`
	if _, err := ps.client.Exec(ctx, sql, ID, pseudoId); err != nil {
		return err
	}
	return nil
}

func (ps *pseudoChapterStorage) Delete(ctx context.Context, chapterID uint64) error {
	sql := `delete from pseudo_chapters where c_id=$1`
	_, err := ps.client.Exec(ctx, sql, chapterID)
	return err
}

func (ps *pseudoChapterStorage) GetIDByPseudo(ctx context.Context, pseudoId string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT c_id FROM "pseudo_chapters" WHERE pseudo = '%s' LIMIT 1`, pseudoId)
	row := ps.client.QueryRow(ctx, sql)
	chapterID := uint64(0)
	err := row.Scan(&chapterID)
	return chapterID, err
}
