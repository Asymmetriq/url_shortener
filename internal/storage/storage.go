package storage

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

type ConcurrentMap struct {
	mtx     sync.RWMutex
	con_map map[string]string
}

func CreateMap() *ConcurrentMap {
	return &ConcurrentMap{con_map: make(map[string]string),
		mtx: sync.RWMutex{}}
}

func (c *ConcurrentMap) Get(s string) (string, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	ogUrl, ok := c.con_map[s]
	if !ok {
		return ogUrl, errors.New("no such value")
	}

	return ogUrl, nil
}

func (c *ConcurrentMap) Set(short_url, og_url string) error {
	c.mtx.Lock()
	c.con_map[short_url] = og_url
	c.mtx.Unlock()
	return nil
}

type StorageInMemory struct {
	Map *ConcurrentMap
}

func (s *StorageInMemory) GetURLByShortURL(_ context.Context, shortURL string) (string, error) {
	return s.Map.Get(shortURL)
}

func (s *StorageInMemory) SaveShortURL(_ context.Context, shortURL, ogUrl string) error {
	_, err := s.Map.Get(shortURL)
	if err == nil {
		return errors.New("already exists")
	}

	return s.Map.Set(shortURL, ogUrl)

}

type StorageDB struct {
	DB *sql.DB
}

func (db *StorageDB) GetURLByShortURL(ctx context.Context, s string) (string, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT url FROM Links WHERE short_url = $1", s)

	var url string
	err := row.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (db *StorageDB) SaveShortURL(ctx context.Context, shortUrl, ogUrl string) error {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO Links (short_url, url) VALUES ($1, $2)", shortUrl, ogUrl)
	if err != nil {
		return err
	}
	return nil
}
