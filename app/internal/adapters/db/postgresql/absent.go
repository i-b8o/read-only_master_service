package postgressql

import (
	"context"
	"fmt"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type absentStorage struct {
	client client.PostgreSQLClient
}

func NewAbsentStorage(client client.PostgreSQLClient) *absentStorage {
	return &absentStorage{client: client}
}

// Create
func (as *absentStorage) Create(ctx context.Context, pseudo string, paragraphId uint64) error {
	sql := `INSERT INTO absent_reg ("pseudo", "paragraph_id") VALUES ($1,$2) `
	if _, err := as.client.Exec(ctx, sql, pseudo, paragraphId); err != nil {
		return err
	}
	return nil
}

// Delete
func (as *absentStorage) DeleteForParagraph(ctx context.Context, paragraphID uint64) error {
	sql := `delete from absent_reg where paragraph_id=$1`
	_, err := as.client.Exec(ctx, sql, paragraphID)
	return err
}

func (as *absentStorage) Done(ctx context.Context, pseudo string) error {
	sql := fmt.Sprintf(`UPDATE absent_reg SET done = true where pseudo= '%s'`, pseudo)
	_, err := as.client.Exec(ctx, sql)
	return err
}
