package postgres

import (
	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"

	// This imports is required to add Postgres driver and make sqlx postgres calls
	_ "github.com/lib/pq"
)

type postgresConfig struct {
	ConnString string `env:"POSTGRES_ADDRESS,required"`
}

// Client is a postgres client
type Client struct {
	dbConn *sqlx.DB
}

// Start creates a new client and connects postgresConfig.ConnString
func Start() (storageClient *Client, err error) {
	// Retreive Env config
	config := postgresConfig{}
	err = env.Parse(&config)
	if err != nil {
		return
	}

	dbConn, err := sqlx.Connect("postgres", config.ConnString)
	if err != nil {
		return
	}

	storageClient = &Client{dbConn}
	return
}

// Begin begins transaction
func (c *Client) Begin() *sqlx.Tx {
	return c.dbConn.MustBegin()
}

// Commit commits current transaction
func (c *Client) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

// Rollback rollback current transaction
func (c *Client) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}

// Close closes current connexion
func (c *Client) Close() {
	c.dbConn.Close()
}
