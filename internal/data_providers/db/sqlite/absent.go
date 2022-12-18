package postgressql

import (
	"context"
	"fmt"
	"read-only_master_service/internal/domain/entity"
	client "read-only_master_service/pkg/client/sqlite"
)

type absentStorage struct {
	client client.SQLiteClient
}

func NewAbsentStorage(client client.SQLiteClient) *absentStorage {
	return &absentStorage{client: client}
}

// Create
func (as *absentStorage) Create(ctx context.Context, absent entity.Absent) error {
	sql := `INSERT INTO absent_reg ("pseudo", "paragraph_id") VALUES ($1,$2) `
	if _, err := as.client.Exec(sql, absent.Pseudo, absent.ParagraphID); err != nil {
		return err
	}
	return nil
}

// Create
func (as *absentStorage) GetAll(ctx context.Context) ([]*entity.Absent, error) {
	sql := `select id, pseudo, done, paragraph_id from absent_reg `
	var absents []*entity.Absent

	rows, err := as.client.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		absent := &entity.Absent{}
		err = rows.Scan(&absent.ID, &absent.Pseudo, &absent.Done, &absent.ParagraphID)
		if err != nil {
			return nil, err
		}

		absents = append(absents, absent)
	}

	return absents, nil
}

// Delete
func (as *absentStorage) DeleteForParagraph(ctx context.Context, paragraphID uint64) error {
	sql := `delete from absent_reg where paragraph_id=$1`
	_, err := as.client.Exec(sql, paragraphID)
	return err
}

func (as *absentStorage) Done(ctx context.Context, pseudo string) error {
	sql := fmt.Sprintf(`UPDATE absent_reg SET done = true where pseudo= '%s'`, pseudo)
	_, err := as.client.Exec(sql)
	return err
}
