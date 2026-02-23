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

func (c *PoolConfig) withDefaults() {
	if c == nil {
		return
	}

	if c.MaxConns == 0 {
		c.MaxConns = maxConn
	}

	if c.MinConns == 0 {
		c.MinConns = 2
	}
	if c.MaxConnLifetime == 0 {
		c.MaxConnLifetime = time.Hour
	}
	if c.MaxConnIdleTime == 0 {
		c.MaxConnIdleTime = 30 * time.Minute
	}

	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == "" {
		c.Port = "5432"
	}

	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
}

// dsn returns the DSN string. If DSN is set explicitly, it is returned as-is.
// Otherwise a DSN is built from the individual fields.
func (c *PoolConfig) dsn() string {
	if c.DSN != "" {
		return c.DSN
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DB,
		c.SSLMode,
	)
}

func (c *PoolConfig) isValid() bool {

	if c == nil {
		return false
	}

	return c.Host != "" &&
		c.Port != "" &&
		c.User != "" &&
		c.Password != "" &&
		c.DB != "" &&
		c.SSLMode != ""
}
