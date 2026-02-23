package redis

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type Client struct {
    client *redis.Client
}

type Redis interface {
    Ping(ctx context.Context) error
    Set(ctx context.Context, key string, value interface{}, expiration int) error
    Get(ctx context.Context, key string) (string, error)
    Del(ctx context.Context, keys ...string) error
    Close() error
}

func New(addr, username, password string, db int) (Redis, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Username: username,
        Password: password,
        DB:       db,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to redis: %w", err)
    }

    return &Client{client: client}, nil
}

func (c *Client) Ping(ctx context.Context) error {
    return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
    return c.client.Close()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration int) error {
    return c.client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
    return c.client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
    return c.client.Del(ctx, keys...).Err()
}