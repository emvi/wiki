package testutil

import (
	auth "emviwiki/auth/model"
	dashboard "emviwiki/dashboard/model"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	backend "emviwiki/shared/model"
	"fmt"
	"time"
)

func CheckOpenConnectionsNull(conn *db.Connection) {
	// wait for goroutines to finish
	time.Sleep(time.Second)

	if conn.Stats().OpenConnections != 0 {
		panic(fmt.Sprintf("OpenConnections must be 0, but was: %v", conn.Stats().OpenConnections))
	}
}

func ConnectBackend(mig bool) *db.Connection {
	c := config.Get()

	if mig {
		migrate(c.BackendDB, "emviwikitest")
	}

	connection := db.NewConnection(db.ConnectionData{
		Host:               c.BackendDB.Host,
		Port:               c.BackendDB.Port,
		User:               c.BackendDB.User,
		Password:           c.BackendDB.Password,
		Schema:             "emviwikitest",
		SSLMode:            "disable",
		MaxOpenConnections: 1,
	})
	connection.SetMaxIdleConns(0)
	backend.SetConnection(connection)
	return connection
}

func ConnectAuth(mig bool) *db.Connection {
	c := config.Get()

	if mig {
		migrate(c.AuthDB, "emviauthtest")
	}

	connection := db.NewConnection(db.ConnectionData{
		Host:               c.AuthDB.Host,
		Port:               c.AuthDB.Port,
		User:               c.AuthDB.User,
		Password:           c.AuthDB.Password,
		Schema:             "emviauthtest",
		SSLMode:            "disable",
		MaxOpenConnections: 1,
	})
	connection.SetMaxIdleConns(0)
	auth.SetConnection(connection)
	return connection
}

func ConnectDashboard(mig bool) *db.Connection {
	c := config.Get()

	if mig {
		migrate(c.DashboardDB, "emviwikidashboardtest")
	}

	connection := db.NewConnection(db.ConnectionData{
		Host:               c.DashboardDB.Host,
		Port:               c.DashboardDB.Port,
		User:               c.DashboardDB.User,
		Password:           c.DashboardDB.Password,
		Schema:             "emviwikidashboardtest",
		SSLMode:            "disable",
		MaxOpenConnections: 1,
	})
	connection.SetMaxIdleConns(0)
	dashboard.SetDashboardConnection(connection)
	return connection
}

func migrate(database config.Database, schema string) {
	if config.Get().Migrate.Dir == "" {
		config.Get().Migrate.Dir = "../schema"
	}

	config.Get().Migrate.Database = database
	config.Get().Migrate.Database.Schema = schema
	config.Get().Migrate.Database.SSLMode = "disable"
	db.Migrate()
}
