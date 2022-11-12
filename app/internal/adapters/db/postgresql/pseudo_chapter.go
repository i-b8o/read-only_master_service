package postgressql

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type pseudoChapterStorage struct {
	client client.PostgreSQLClient
}

func NewPseudoChapterStorage(client client.PostgreSQLClient) *pseudoChapterStorage {
	return &pseudoChapterStorage{client: client}
}

func (pcs *pseudoChapterStorage) CreateRelationship(ctx context.Context, pseudoChapter entity.PseudoChapter) error {
	const sql = `INSERT INTO pseudo_chapters ("c_id", "pseudo") VALUES ($1, $2)`
	if _, err := pcs.client.Exec(ctx, sql, pseudoChapter.ID, pseudoChapter.PseudoId); err != nil {
		return err
	}
	return nil
}

func (pcs *pseudoChapterStorage) DeleteRelationship(ctx context.Context, chapterID uint64) error {
	sql := `delete from pseudo_chapters where c_id=$1`
	_, err := pcs.client.Exec(ctx, sql, chapterID)
	return err
}
