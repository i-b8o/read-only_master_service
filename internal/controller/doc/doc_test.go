package doc_controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"read-only_master_service/internal/domain/entity"
	"read-only_master_service/pkg/client/sqlite"
	"testing"

	pb "github.com/i-b8o/read-only_contracts/pb/master/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const (
	dbPath = ""
)

func setupDB() *sql.DB {
	sqliteConfig := sqlite.NewSqliteConfig(
		dbPath,
	)

	pgClient, err := sqlite.NewClient(sqliteConfig)
	if err != nil {
		log.Fatal(err)
	}

	return pgClient
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30002"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewMasterDocGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input    *pb.CreateDocRequest
		expected *pb.CreateDocResponse
		err      error
	}{
		{
			input:    &pb.CreateDocRequest{PseudoId: "372952", DocName: "Имя правила", Title: "Заголовок"},
			expected: &pb.CreateDocResponse{ID: 2},
			err:      nil,
		},
	}

	for _, test := range tests {
		e, err := client.Create(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		// doc
		sql := fmt.Sprintf("select id, name, abbreviation from doc where id=%d", e.ID)
		rows, err := pgClient.Query(sql)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var id uint64
		var name, abbreviation string
		for rows.Next() {
			if err = rows.Scan(
				&id, &name, &abbreviation,
			); err != nil {
				t.Log(err)
			}
		}

		assert.Equal(e.ID, id)
		assert.Equal(test.input.DocName, name)

		assert.True(proto.Equal(test.expected, e), fmt.Sprintf("CreateDoc(%v)=%v want: %v", test.input, e, test.expected))
		assert.Equal(test.err, err)
		// pseudo doc
		sqlP := fmt.Sprintf("select * from pseudo_doc where r_id=%d", e.ID)
		rows, err = pgClient.Query(sqlP)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var cId uint64
		var pseudo string
		for rows.Next() {
			if err = rows.Scan(
				&cId, &pseudo,
			); err != nil {
				t.Log(err)
			}
		}

		assert.Equal(cId, e.ID)
		assert.Equal(pseudo, test.input.PseudoId)
	}
	_, err = pgClient.Exec(resetDB)
	if err != nil {
		t.Log(err)
	}
}
func TestGetAbsents(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30002"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewMasterDocGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input    *pb.Empty
		expected *pb.GetAbsentsResponse
		err      error
	}{
		{
			input:    &pb.Empty{},
			expected: &pb.GetAbsentsResponse{Absents: []*pb.MasterAbsent{&pb.MasterAbsent{ID: 1, Pseudo: "aaaaa", Done: false, ParagraphId: 1}, &pb.MasterAbsent{ID: 2, Pseudo: "bbbbb", Done: true, ParagraphId: 2}, &pb.MasterAbsent{ID: 3, Pseudo: "ccccc", Done: false, ParagraphId: 3}}},
			err:      nil,
		},
	}

	for _, test := range tests {
		resp, err := client.GetAbsents(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.True(proto.Equal(test.expected, resp), fmt.Sprintf("GetAbsents(%v)=%v want: %v", test.input, resp, test.expected))
		assert.Equal(test.err, err, err)
	}
}

func TestGetAll(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30002"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewMasterDocGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input    *pb.Empty
		expected *pb.GetAllDocsResponse
		err      error
	}{
		{
			input:    &pb.Empty{},
			expected: &pb.GetAllDocsResponse{Docs: []*pb.Doc{&pb.Doc{ID: 1, DocName: "Имя первой записи"}}},
			err:      nil,
		},
	}

	for _, test := range tests {
		resp, err := client.GetAll(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.True(proto.Equal(test.expected, resp), fmt.Sprintf("GetAll(%v)=%v want: %v", test.input, resp, test.expected))
		assert.Equal(test.err, err, err)
	}
}

func TestUpdateLinks(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30002"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewMasterDocGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = pgClient.Exec(resetDB)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		input      *pb.UpdateLinksRequest
		expected   *pb.UpdateLinksResponse
		paragraphs []entity.Paragraph
		absent     []entity.Absent
		err        error
	}{
		{
			input:      &pb.UpdateLinksRequest{ID: 1},
			expected:   &pb.UpdateLinksResponse{ID: 1},
			paragraphs: []entity.Paragraph{entity.Paragraph{ID: 1, Num: 1, HasLinks: true, IsTable: false, IsNFT: false, Class: "any-class", Content: "Содержимое <a id=\"dst101675\"></a> первого <a href='1/3/111'>параграфа</a>", ChapterID: 1}, entity.Paragraph{ID: 2, Num: 2, IsTable: true, IsNFT: true, HasLinks: true, Class: "any-class", Content: "Содержимое второго <a href='372952/4e92c731969781306ebd1095867d2385f83ac7af/335104'>пункта 5.14</a> параграфа", ChapterID: 1}, entity.Paragraph{ID: 3, Num: 3, IsTable: false, IsNFT: false, HasLinks: true, Class: "any-class", Content: "<a id='335050'></a>Содержимое третьего параграфа<a href='/document/cons_doc_LAW_2875/'>таблицей N 2</a>.", ChapterID: 1}},
			absent:     []entity.Absent{entity.Absent{Pseudo: "aaaaa", Done: false, ParagraphID: 1}, entity.Absent{Pseudo: "bbbbb", Done: false, ParagraphID: 2}, entity.Absent{Pseudo: "ccccc", Done: false, ParagraphID: 3}, entity.Absent{Pseudo: "372952", Done: false, ParagraphID: 2}, entity.Absent{Pseudo: "2875", Done: false, ParagraphID: 3}},

			err: nil,
		},
	}

	for _, test := range tests {
		e, err := client.UpdateLinks(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.True(proto.Equal(test.expected, e), fmt.Sprintf("CreateDoc(%v)=%v want: %v", test.input, e, test.expected))
		assert.Equal(test.err, err)

		// paragraph
		sql := "select paragraph_id, order_num, is_table, is_nft, has_links, class, content, c_id from paragraph"
		rows, err := pgClient.Query(sql)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var paragraphs []entity.Paragraph
		for rows.Next() {
			var p entity.Paragraph
			if err = rows.Scan(
				&p.ID, &p.Num, &p.IsTable, &p.IsNFT, &p.HasLinks, &p.Class, &p.Content, &p.ChapterID,
			); err != nil {
				t.Log(err)
			}
			paragraphs = append(paragraphs, p)
		}
		assert.Equal(test.paragraphs, paragraphs)

		// abscent
		sql1 := "select pseudo, paragraph_id from absent_reg"
		rows, err = pgClient.Query(sql1)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var abscents []entity.Absent
		for rows.Next() {
			var a entity.Absent
			if err = rows.Scan(
				&a.Pseudo, &a.ParagraphID,
			); err != nil {
				t.Log(err)
			}
			abscents = append(abscents, a)
		}
		assert.Equal(test.absent, abscents)
	}
	_, err = pgClient.Exec(resetDB)
	if err != nil {
		t.Log(err)
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30002"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewMasterDocGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input *pb.DeleteDocRequest
		err   error
	}{
		{
			input: &pb.DeleteDocRequest{ID: 1},
			err:   nil,
		},
	}

	for _, test := range tests {
		_, err := client.Delete(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(test.err, err)

		// paragraph
		sql := "select id from paragraph"
		rows, err := pgClient.Query(sql)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var paragraphs []entity.Paragraph
		for rows.Next() {
			var p entity.Paragraph
			if err = rows.Scan(
				&p.ID,
			); err != nil {
				t.Log(err)
			}
			paragraphs = append(paragraphs, p)
		}
		assert.True(len(paragraphs) == 0)

		// chapter
		sql1 := "select id from chapter"
		rows, err = pgClient.Query(sql1)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var chapters []entity.Chapter
		for rows.Next() {
			var c entity.Chapter
			if err = rows.Scan(
				c.ID,
			); err != nil {
				t.Log(err)
			}
			chapters = append(chapters, c)
		}
		assert.True(len(chapters) == 0)
		// abscent
		sql2 := "select id from doc"
		rows, err = pgClient.Query(sql2)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()

		var docs []entity.Doc
		for rows.Next() {
			var r entity.Doc
			if err = rows.Scan(
				&r.Id,
			); err != nil {
				t.Log(err)
			}
			docs = append(docs, r)
		}
		assert.True(len(docs) == 0)
	}
	_, err = pgClient.Exec(resetDB)
	if err != nil {
		t.Log(err)
	}
}

const resetDB = `
DROP TABLE IF EXISTS absent_reg;
DROP TABLE IF EXISTS pseudo_chapter;
DROP TABLE IF EXISTS pseudo_doc;
DROP TABLE IF EXISTS link;
DROP MATERIALIZED VIEW IF EXISTS reg_search;
DROP INDEX IF EXISTS idx_search;
DROP TABLE IF EXISTS paragraph;
DROP TABLE IF EXISTS chapter;
DROP TABLE IF EXISTS doc;


CREATE TABLE doc (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (NAME != '') UNIQUE,
    abbreviation TEXT,
    title TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE doc ADD COLUMN ts tsvector GENERATED ALWAYS AS (setweight(to_tsvector('russian', coalesce(name, '')), 'A') || setweight(to_tsvector('russian', coalesce(title, '')), 'B')) STORED;
CREATE INDEX reg_ts_idx ON doc USING GIN (ts);



CREATE TABLE chapter (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name != ''),
    order_num SMALLINT NOT NULL CHECK (order_num >= 0),
    num TEXT,
    r_id integer REFERENCES doc,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE chapter ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', name)) STORED;
CREATE INDEX ch_ts_idx ON chapter USING GIN (ts);

CREATE TABLE paragraph (
    id SERIAL PRIMARY KEY,
    paragraph_id INT NOT NULL CHECK (paragraph_id >= 0),
    order_num INT NOT NULL CHECK (order_num >= 0),
    is_table BOOLEAN NOT NULL,
    is_nft BOOLEAN NOT NULL,
    has_links BOOLEAN NOT NULL,
    class TEXT,
    content TEXT NOT NULL,
    c_id integer REFERENCES chapter
);

ALTER TABLE paragraph ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', content)) STORED;
CREATE INDEX p_ts_idx ON paragraph USING GIN (ts);


CREATE MATERIALIZED VIEW reg_search 
AS SELECT 
r.id AS "r_id", r.name AS "r_name", NULL AS "c_id", NULL AS "c_name", CAST(NULL AS integer) AS "p_id", NULL AS "p_text", r.name AS "text",
to_tsvector('russian', r.name) AS ts FROM doc AS r UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", NULL AS "p_id", NULL AS "p_text", c.name AS "text",
to_tsvector('russian', c.name) AS ts FROM chapter AS c INNER JOIN doc AS r ON r.id= c.r_id
UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", p.paragraph_id AS "p_id", p.content AS "p_text", p.content AS "text",
to_tsvector('russian', content) AS ts 
FROM paragraph AS p INNER JOIN chapter AS c ON p.c_id= c.id INNER JOIN doc AS r ON c.r_id = r.id;

create index idx_search on reg_search using GIN(ts);

CREATE TABLE pseudo_doc (
    r_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE pseudo_chapter (
    c_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE absent_reg (
    id SERIAL PRIMARY KEY,
    pseudo TEXT NOT NULL CHECK (pseudo != ''),
    done BOOLEAN NOT NULL DEFAULT false,
    paragraph_id integer  
);

CREATE TABLE link (
    id INT NOT NULL UNIQUE,
    paragraph_num INT NOT NULL CHECK (paragraph_num >= 0),
    c_id integer,
    r_id integer
);

INSERT INTO doc ("name", "abbreviation", "title", "created_at") VALUES ('Имя первой записи', 'Аббревиатура первой записи', 'Заголовок первой записи', '2023-01-01 00:00:00');
INSERT INTO chapter ("name", "num", "order_num","r_id", "updated_at") VALUES ('Имя первой записи','I',1,1, '2023-01-01 00:00:00'), ('Имя второй записи','II',2,1, '2023-01-01 00:00:00'), ('Имя третьей записи','III',3,1, '2023-01-01 00:00:00');
INSERT INTO paragraph ("paragraph_id","order_num","is_table","is_nft","has_links","class","content","c_id") VALUES (1,1,false,false,true,'any-class','Содержимое <a id="dst101675"></a> первого <a href=''11111/a3a3a3/111''>параграфа</a>', 1), (2,2,true,true,true,'any-class','Содержимое второго <a href=''372952/4e92c731969781306ebd1095867d2385f83ac7af/335104''>пункта 5.14</a> параграфа', 1), (3,3,false,false,true,'any-class','<a id=''335050''></a>Содержимое третьего параграфа<a href=''/document/cons_doc_LAW_2875/''>таблицей N 2</a>.', 1);
INSERT INTO pseudo_doc ("r_id", "pseudo") VALUES (1, 11111);
INSERT INTO pseudo_chapter ("c_id", "pseudo") VALUES (3, 'a3a3a3');
INSERT INTO absent_reg ("pseudo", "done", "paragraph_id") VALUES ('aaaaa', false, 1), ('bbbbb', true, 2), ('ccccc', false, 3);`
