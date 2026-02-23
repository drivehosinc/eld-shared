package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is a PostgreSQL client with optional master/replica routing.
type Client struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool // nil when no replica configured
}

// New creates connection pools, pings both, and returns a ready Client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	masterPool, err := newPool(ctx, cfg.Master)
	if err != nil {
		return nil, fmt.Errorf("postgresql: master pool: %w", err)
	}

	var replicaPool *pgxpool.Pool
	if cfg.Replica != nil && cfg.Replica.isValid() {
		replicaPool, err = newPool(ctx, *cfg.Replica)
		if err != nil {
			masterPool.Close()
			return nil, fmt.Errorf("postgresql: replica pool: %w", err)
		}
	}

	return &Client{master: masterPool, replica: replicaPool}, nil
}

func newPool(ctx context.Context, pc PoolConfig) (*pgxpool.Pool, error) {
	pc.withDefaults()

	poolCfg, err := pgxpool.ParseConfig(pc.dsn())
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolCfg.MaxConns = pc.MaxConns
	poolCfg.MinConns = pc.MinConns
	poolCfg.MaxConnLifetime = pc.MaxConnLifetime
	poolCfg.MaxConnIdleTime = pc.MaxConnIdleTime

	if poolCfg.ConnConfig.RuntimeParams == nil {
		poolCfg.ConnConfig.RuntimeParams = make(map[string]string)
	}
	poolCfg.ConnConfig.RuntimeParams["application_name"] = pc.ServiceName

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}

// pool returns the appropriate pool based on context routing.
func (c *Client) pool(ctx context.Context) *pgxpool.Pool {
	if modeFromCtx(ctx) == modeReplica && c.replica != nil {
		return c.replica
	}
	return c.master
}

// Pool provides direct pool access based on context routing.
func (c *Client) Pool(ctx context.Context) *pgxpool.Pool {
	return c.pool(ctx)
}

// Exec executes a query that doesn't return rows, routed by context.
func (c *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.pool(ctx).Exec(ctx, sql, args...)
}

// Query executes a query that returns rows, routed by context.
func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool(ctx).Query(ctx, sql, args...)
}

// QueryRow executes a query that returns at most one row, routed by context.
func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool(ctx).QueryRow(ctx, sql, args...)
}

// Begin starts a transaction on the master pool.
func (c *Client) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.master.Begin(ctx)
}

// BeginTx starts a transaction with options on the master pool.
func (c *Client) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.master.BeginTx(ctx, txOptions)
}

// WithTransaction executes fn inside a transaction with automatic commit/rollback and panic recovery.
func (c *Client) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) (err error) {
	tx, err := c.master.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgresql: begin tx: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	err = fn(ctx, tx)
	return err
}

// CopyFrom performs a bulk copy into a table on the master pool.
func (c *Client) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.master.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// SendBatch sends a batch of queries, routed by context.
func (c *Client) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.pool(ctx).SendBatch(ctx, b)
}

// Ping verifies connectivity to both pools.
func (c *Client) Ping(ctx context.Context) error {
	if err := c.master.Ping(ctx); err != nil {
		return fmt.Errorf("postgresql: ping master: %w", err)
	}
	if c.replica != nil {
		if err := c.replica.Ping(ctx); err != nil {
			return fmt.Errorf("postgresql: ping replica: %w", err)
		}
	}
	return nil
}

// Close closes both connection pools.
func (c *Client) Close() {
	c.master.Close()
	if c.replica != nil {
		c.replica.Close()
	}
}
