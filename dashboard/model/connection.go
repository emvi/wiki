package model

import (
	"emviwiki/shared/db"
)

var (
	backendDB   *db.Connection
	authDB      *db.Connection
	dashboardDB *db.Connection
)

// SetConnection sets the connection pool to work with.
func SetConnection(backdb, authdb, dashdb *db.Connection) {
	backendDB = backdb
	authDB = authdb
	dashboardDB = dashdb
}

// SetDashboardConnection sets the dashboard connection pool to work with.
func SetDashboardConnection(connection *db.Connection) {
	dashboardDB = connection
}

// GetConnection returns the dashboard connection pool.
func GetConnection() *db.Connection {
	return dashboardDB
}
