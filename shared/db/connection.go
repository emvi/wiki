package db

import (
	"database/sql"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

const (
	connectTimeout            = "60"
	defaultHost               = "localhost"
	defaultPort               = "5432"
	defaultUser               = "postgres"
	defaultPassword           = "postgres"
	defaultSchema             = "postgres"
	defaultSSLMode            = "disable"
	defaultMaxOpenConnections = 1
)

const (
	connectionString = `host=%s
		port=%s
		user=%s
		password=%s
		dbname=%s
		sslmode=%s
		sslcert=%s
		sslkey=%s
		sslrootcert=%s
		connectTimeout=%s
		timezone=%s`
)

// Connection is a connection pool to a database with extended functionality.
type Connection struct {
	sqlx.DB
}

// ConnectionData is the required data to connect to a database.
// If a parameter is left empty, the default value will be used instead.
type ConnectionData struct {
	Host               string
	Port               string
	User               string
	Password           string
	Schema             string
	SSLMode            string
	SSLCert            string
	SSLKey             string
	SSLRootCert        string
	MaxOpenConnections int
}

// NewConnection returns a new database connection using given configuration.
func NewConnection(data ConnectionData) *Connection {
	logbuch.Info("Connecting to database...")
	db, err := sqlx.Connect("postgres", postgresConnection(&data))

	if err != nil {
		logbuch.Fatal("Error connecting to database", logbuch.Fields{"err": err})
		return nil
	}

	if err := db.Ping(); err != nil {
		logbuch.Fatal("Error pinging database", logbuch.Fields{"err": err})
		return nil
	}

	logbuch.Info("Connected")
	setMaxOpenConnections(db, data.MaxOpenConnections)
	return &Connection{DB: *db}
}

// Disconnect closes the database connection.
func (connection *Connection) Disconnect() {
	logbuch.Info("Disconnecting from database...")

	if err := connection.Close(); err != nil {
		logbuch.Warn("Error when closing database connection", logbuch.Fields{"err": err})
	}

	logbuch.Info("Disconnected")
}

// SaveEntity saves the given entity using the insert or update statement,
// depending on whether the ID is set or not.
func (connection *Connection) SaveEntity(tx *sqlx.Tx, entity Entity, insert, update string) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer Commit(tx)
	}

	var err error

	if entity.GetId() == 0 {
		var rows *sqlx.Rows
		rows, err = tx.NamedQuery(insert, entity)

		if err == nil {
			defer CloseRows(rows)
			rows.Next()
			var id hide.ID

			if err := rows.Scan(&id); err != nil {
				logbuch.Error("Error scanning entity id", logbuch.Fields{"err": err, "entity": entity})
				Rollback(tx)
				return err
			}

			entity.SetId(id)
		}
	} else {
		_, err = tx.NamedExec(update, entity)
	}

	if err != nil {
		logbuch.Error("Error saving entity", logbuch.Fields{"err": err, "entity": entity})
		Rollback(tx)
		return err
	}

	return nil
}

// Exec executes a single SQL statement in given transaction or creates a new one if nil.
// This function should not be used for production code and is intended to be used within tests only.
func (connection *Connection) Exec(tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer Commit(tx)
	}

	result, err := tx.Exec(query, args...)

	if err != nil {
		logbuch.Error("Error executing sql statement", logbuch.Fields{"err": err})
		return result, err
	}

	return result, nil
}

func postgresConnection(data *ConnectionData) string {
	if data.Host == "" {
		data.Host = defaultHost
	}

	if data.Port == "" {
		data.Port = defaultPort
	}

	if data.User == "" {
		data.User = defaultUser
	}

	if data.Password == "" {
		data.Password = defaultPassword
	}

	if data.Schema == "" {
		data.Schema = defaultSchema
	}

	if data.SSLMode == "" {
		data.SSLMode = defaultSSLMode
	}

	zone, offset := time.Now().Zone()
	timezone := zone + strconv.Itoa(-offset/3600)
	logbuch.Info("Setting time zone", logbuch.Fields{"timezone": timezone})
	return fmt.Sprintf(connectionString, data.Host, data.Port, data.User, data.Password, data.Schema, data.SSLMode, data.SSLCert, data.SSLKey, data.SSLRootCert, connectTimeout, timezone)
}

func setMaxOpenConnections(db *sqlx.DB, maxOpenConnections int) {
	if maxOpenConnections == 0 {
		maxOpenConnections = defaultMaxOpenConnections
	}

	logbuch.Info("Setting max open connections to database", logbuch.Fields{"max_open_connections": maxOpenConnections})
	db.SetMaxOpenConns(maxOpenConnections)
}
