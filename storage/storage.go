package storage

import (
	"fmt"

	sq "github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib" //postgres driver for sqlx
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

type ITbStorage interface {
	Close()
	Migrate() error

	InsertNewsItem(news []NewsItem) error
	GetTagContentList(filter TagContentFilter) ([]TagContent, error)
	GetTagContent(id int64) (TagContent, error)
}

type TbStorage struct {
	ITbStorage
	db *sqlx.DB
}

func New(dataSourceName string) (*TbStorage, error) {
	db, err := sqlx.Connect("pgx", dataSourceName)
	return &TbStorage{db: db}, err
}

func (s *TbStorage) Close() {
	_ = s.db.Close()
}

func (s *TbStorage) Migrate() error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Run("up", s.db.DB, "./migrations")
	if err == goose.ErrNoNextVersion {
		return nil
	}

	return err
}

func (s *TbStorage) InsertNewsItem(news []NewsItem) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, newsItem := range news {
		_, err = tx.Exec(`INSERT INTO news (url, tag_content) VALUES($1, $2)`,
			newsItem.URL, newsItem.TagContent)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *TbStorage) GetTagContentList(filter TagContentFilter) ([]TagContent, error) {
	tagContentList := []TagContent{}
	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := stmt.Select(`id, tag_content, url, created_at::text`).
		From(`news`).OrderBy(`url`)

	if len(filter.ContentKeyword) > 0 {
		q = q.Where(`tag_content ILIKE ?`, fmt.Sprint("%", filter.ContentKeyword, "%"))
	}
	if len(filter.URL) > 0 {
		q = q.Where(`url ILIKE ?`, fmt.Sprint("%", filter.URL, "%"))
	}
	if len(filter.CreatedAtFrom) > 0 {
		q = q.Where(`created_at >= ?`, filter.CreatedAtFrom)
	}
	if len(filter.CreatedAtTo) > 0 {
		q = q.Where(`created_at <= ?`, filter.CreatedAtTo)
	}

	rows, err := q.RunWith(s.db).Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		tc := TagContent{}
		err = rows.Scan(
			&tc.ID,
			&tc.TagContent,
			&tc.URL,
			&tc.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tagContentList = append(tagContentList, tc)
	}

	return tagContentList, err
}

func (s *TbStorage) GetTagContent(id int64) (TagContent, error) {
	tc := TagContent{ID: id}
	err := s.db.QueryRow(`SELECT tag_content, url, created_at::text FROM news WHERE id = $1`, id).
		Scan(
			&tc.TagContent,
			&tc.URL,
			&tc.CreatedAt,
		)

	return tc, err
}
