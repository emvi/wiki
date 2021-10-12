package model

import (
	"emviwiki/shared/db"
)

var (
	connection *db.Connection
)

// SetConnection sets the connection pool to work with.
func SetConnection(c *db.Connection) {
	connection = c
}

// GetConnection returns the current connection pool.
func GetConnection() *db.Connection {
	return connection
}
