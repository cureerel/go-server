package dbtypes

import (
	"context"
	"io"
)

type DBClient interface {
	io.Closer
	Ping(ctx context.Context) error
}

type SQLDB interface {
	DBClient
	GormDB() interface{}
}

type MongoDB interface {
	DBClient
	Database() interface{}
	Collection(name string) interface{}
}

type Redis interface {
	DBClient
	Client() interface{}
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}
