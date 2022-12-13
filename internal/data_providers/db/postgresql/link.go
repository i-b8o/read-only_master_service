package postgressql

import (
	"context"
	"read-only_master_service/internal/domain/entity"
	client "read-only_master_service/pkg/client/postgresql"
)

type linkStorage struct {
	client client.PostgreSQLClient
}

func NewLinkStorage(client client.PostgreSQLClient) *linkStorage {
	return &linkStorage{client: client}
}

// GetAll returns all links
func (ps *linkStorage) GetAll(ctx context.Context) ([]*entity.Link, error) {
	const sql = `SELECT id,c_id,paragraph_num FROM "link" ORDER BY c_id`

	var links []*entity.Link

	rows, err := ps.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		link := &entity.Link{}
		err = rows.Scan(&link.ID, &link.ChapterID, &link.ParagraphNum)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

// Create
func (ps *linkStorage) Create(ctx context.Context, link entity.Link) error {
	sql := `INSERT INTO link ("id","c_id","paragraph_num","r_id") VALUES ($1,$2,$3,$4)`
	_, err := ps.client.Exec(ctx, sql, link.ID, link.ChapterID, link.ParagraphNum, link.RID)
	return err
}

// Create
func (ps *linkStorage) CreateForChapter(ctx context.Context, link entity.Link) error {
	sql := `INSERT INTO link ("id","c_id","paragraph_num","r_id") VALUES ($1,$2,$3,$4) ON CONFLICT ("id") DO NOTHING`
	_, err := ps.client.Exec(ctx, sql, link.ID, link.ChapterID, link.ParagraphNum, link.RID)
	return err
}

// Delete
func (ps *linkStorage) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	sql := `delete from link where c_id =$1`
	_, err := ps.client.Exec(ctx, sql, chapterID)
	return err
}
