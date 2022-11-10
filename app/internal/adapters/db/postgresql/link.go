package postgressql

import (
	"context"
	"regulations_supreme_service/internal/domain/entity"
	client "regulations_supreme_service/pkg/client/postgresql"
)

type linkStorage struct {
	client client.PostgreSQLClient
}

func NewLinkStorage(client client.PostgreSQLClient) *linkStorage {
	return &linkStorage{client: client}
}

// GetAll returns all links
func (ps *linkStorage) GetAll(ctx context.Context) ([]*entity.Link, error) {
	const sql = `SELECT id,c_id,paragraph_num FROM "links" ORDER BY c_id`

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

// // GetAllByChapterID returns all links associated with the given chapter ID
// func (ps *linkStorage) GetAllByChapterID(ctx context.Context, chapterID uint64) ([]*entity.Link, error) {
// 	const sql = `SELECT id,c_id,paragraph_num FROM "links" WHERE c_id = $1 ORDER BY c_id`
// 	var links []*entity.Link
// 	rows, err := ps.client.Query(ctx, sql, chapterID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		link := &entity.Link{}
// 		if err = rows.Scan(
// 			&link.ID, &link.ChapterID, &link.ParagraphNum); err != nil {
// 			return nil, err
// 		}

// 		links = append(links, link)
// 	}

// 	return links, nil
// }

// Create
func (ps *linkStorage) Create(ctx context.Context, link entity.Link) error {
	sql := `INSERT INTO links ("id","c_id","paragraph_num","r_id") VALUES ($1,$2,$3,$4)`
	_, err := ps.client.Exec(ctx, sql, link.ID, link.ChapterID, link.ParagraphNum, link.RID)
	return err
}

// Create
func (ps *linkStorage) CreateForChapter(ctx context.Context, link entity.Link) error {
	sql := `INSERT INTO links ("id","c_id","paragraph_num","r_id") VALUES ($1,$2,$3,$4) ON CONFLICT ("id") DO NOTHING`
	_, err := ps.client.Exec(ctx, sql, link.ID, link.ChapterID, link.ParagraphNum, link.RID)
	return err
}

// Delete
func (ps *linkStorage) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	sql := `delete from links where c_id =$1`
	_, err := ps.client.Exec(ctx, sql, chapterID)
	return err
}

// func (ps *linkStorage) GetOneByParagraphID(ctx context.Context, paragraphID, regregulationID uint64) (entity.Link, error) {
// 	const sqlCheck = `SELECT EXISTS(SELECT 1 FROM links WHERE id=$1 AND r_id=$2)`
// 	var exist bool
// 	checkRow := ps.client.QueryRow(ctx, sqlCheck, paragraphID, regregulationID)
// 	var link entity.Link

// 	err := checkRow.Scan(&exist)
// 	if err != nil {
// 		return entity.Link{ID: 0}, err
// 	}

// 	if !exist {
// 		return link, nil
// 	}

// 	const sql = `SELECT  id, paragraph_num, c_id FROM "links" WHERE id=$1 AND r_id=$2 LIMIT 1`
// 	row := ps.client.QueryRow(ctx, sql, paragraphID, regregulationID)

// 	err = row.Scan(&link.ID, &link.ParagraphNum, &link.ChapterID)
// 	if err != nil {
// 		return link, err
// 	}

// 	return link, nil
// }
