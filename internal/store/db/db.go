package db

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ineverbee/bitlybutnot/internal/store"
	"github.com/jackc/pgx/v4"
)

//Options type holds connection parameters
type Options struct {
	Username string
	Password string
	Host     string
	Port     int
	Timeout  int
}

//LinkDB type represents postgresql storage
type LinkDB struct {
	db *pgx.Conn
}

//LinkDB.Set method inserts rows in the storage
func (s *LinkDB) Set(key uint32, shortURL, longURL string) error {
	_, err := s.db.Exec(context.Background(), "insert into links (id, shortURL, longURL) values ($1, $2, $3)", key, shortURL, longURL)
	if err != nil {
		return err
	}
	return nil
}

//LinkDB.Get method gets specific row from the storage
func (s *LinkDB) Get(key uint32) (string, error) {
	var url string
	err := s.db.QueryRow(context.Background(), "select longURL from links where id=$1", key).Scan(&url)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("no record found")
		}
		return "", err
	}
	return url, nil
}

//NewLinkDB func creates new connection to db and returns store instance
func NewLinkDB(ctx context.Context, l func(format string, v ...interface{}), o *Options) (store.Store, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/?sslmode=disable&connect_timeout=%d",
		url.QueryEscape(o.Username),
		url.QueryEscape(o.Password),
		o.Host,
		o.Port,
		o.Timeout,
	)

	l("Trying to connect to %s\n", connStr)
	var (
		conn *pgx.Conn
		err  error
	)

	ticker := time.NewTicker(1 * time.Second)
	timeout := 20 * time.Second
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
LOOP:
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			conn, err = pgx.Connect(ctx, connStr)
			if err == nil {
				break LOOP
			}
			l("Failed to connect to db %s, %v", connStr, err)
			l("Trying to reconnect..")
		}
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	l("Connect success!")

	return &LinkDB{db: conn}, nil
}
