package postgresql

import (
	"fmt"
	"time"
)

type PoolConfig struct {
	ServiceName     string
	DSN             string
	Host            string
	Port            string
	User            string
	Password        string
	DB              string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type Config struct {
	Master  PoolConfig
	Replica *PoolConfig // nil = no replica, all queries go to master
}

const (
	maxConn = 20
)

func (c *PoolConfig) withDefaults() PoolConfig {
	out := *c
	if out.MaxConns == 0 {
		out.MaxConns = maxConn
	}

	if out.MinConns == 0 {
		out.MinConns = 2
	}
	if out.MaxConnLifetime == 0 {
		out.MaxConnLifetime = time.Hour
	}
	if out.MaxConnIdleTime == 0 {
		out.MaxConnIdleTime = 30 * time.Minute
	}
	return out
}

// dsn returns the DSN string. If DSN is set explicitly, it is returned as-is.
// Otherwise a DSN is built from the individual fields.
func (c *PoolConfig) dsn() string {
	if c.DSN != "" {
		return c.DSN
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Port
	if port == "" {
		port = "5432"
	}
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		c.User,
		c.Password,
		c.DB,
		c.SSLMode,
	)
}
