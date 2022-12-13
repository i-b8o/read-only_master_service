package postgressql

// TODO delete
// import (
// 	"context"
// 	"read-only_master_service/internal/domain/entity"
// 	client "read-only_master_service/pkg/client/postgresql"
// 	"strconv"
// )

// type speechStorage struct {
// 	client client.PostgreSQLClient
// }

// func NewSpeechStorage(client client.PostgreSQLClient) *speechStorage {
// 	return &speechStorage{client: client}
// }

// // GetOneById returns all speech text associated with the given paragraph ID
// // func (ss *speechStorage) GetAllById(ctx context.Context, paragraphID uint64) ([]entity.Speech, error) {
// // 	const sql = `select order_num, content from speech where paragraph_id = $1 ORDER BY order_num`

// // 	var speechSlice []entity.Speech

// // 	rows, err := ss.client.Query(ctx, sql, paragraphID)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	defer rows.Close()

// // 	for rows.Next() {
// // 		speech := entity.Speech{}
// // 		if err = rows.Scan(
// // 			&speech.OrderNum, &speech.Content,
// // 		); err != nil {
// // 			return nil, err
// // 		}

// // 		speechSlice = append(speechSlice, speech)
// // 	}

// // 	return speechSlice, nil
// // }

// // Create returns the ID of the inserted chapter
// func (ss *speechStorage) Create(ctx context.Context, speech entity.Speech) (string, error) {
// 	const sql = `INSERT INTO speech ("paragraph_id", "content", "order_num") VALUES ($1, $2, $3) RETURNING "id"`

// 	row := ss.client.QueryRow(ctx, sql, speech.ParagraphID, speech.Content, speech.OrderNum)
// 	var pararaphID uint64

// 	err := row.Scan(&pararaphID)
// 	if err != nil {
// 		return "0", err
// 	}

// 	return strconv.FormatUint(pararaphID, 10), nil
// }

// // Delete
// func (ps *speechStorage) DeleteForParagraph(ctx context.Context, paragraphID uint64) error {
// 	sql := `delete from speech where paragraph_id=$1`
// 	_, err := ps.client.Exec(ctx, sql, paragraphID)
// 	return err
// }
