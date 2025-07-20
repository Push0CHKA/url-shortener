package sqlite

import (
	"database/sql"
	"errors"
	"url-shortener/internal/lib/e"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (storage *Storage, err error) {
	defer func() { err = e.WrapIfErr("storage.sqlite.New", err) }()

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias) 
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (id int64, err error) {
	defer func() { err = e.WrapIfErr("storage.sqlite.SaveURL", err) }()

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrURLExist
		}

		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetURLbyAlias(alias string) (resURL string, err error) {
	defer func() { err = e.WrapIfErr("storage.sqlite.GetURLbyAlias", err) }()

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", err
	}

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", err
	}

	return resURL, nil
}

func (s *Storage) DeleteURLbyAlias(alias string) (err error) {
	defer func() { err = e.WrapIfErr("storage.sqlite.DeleteURLbyAlias", err) }()

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return err
	}

	return nil
}
