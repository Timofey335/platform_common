package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/Timofey335/platform_common/pkg/db"
)

type pgClient struct {
	masterDBC db.DB
}

// New - создает новый клиент к БД
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

// DB - создает соединение с БД, сделан для расширение возможностей,
// если необходимо будет добавить подключение к другой или нескольким БД
func (p *pgClient) DB() db.DB {
	return p.masterDBC
}

// Close - закрывает соединение с БД
func (p *pgClient) Close() error {
	if p.masterDBC != nil {
		p.masterDBC.Close()
	}
	return nil
}
