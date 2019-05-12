package routing

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type PgRouting struct {
	routingClient *sql.DB
	auxClient     *sql.DB
}

// NewRouting initialize and return an new PgRouting object.
func NewPgRoutingPrep(DBUser, DBPass, DBName, DBPort, AuxDBName string) (*PgRouting, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, DBName, DBPort)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	connStr = fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, AuxDBName, DBPort)
	auxDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PgRouting{
		routingClient: DB,
		auxClient:     auxDB,
	}, nil
}

func (r *PgRouting) Prepare() error {
	_, err := r.auxClient.Exec(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity  WHERE datname = 'routing'`)
	if err != nil {
		return err
	}

	_, err = r.auxClient.Exec("DROP DATABASE IF EXISTS " + os.Getenv("DBNAME"))
	if err != nil {
		return err
	}

	_, err = r.auxClient.Exec("CREATE DATABASE " + os.Getenv("DBNAME"))
	if err != nil {
		return err
	}

	_, err = r.routingClient.Exec("CREATE EXTENSION postGIS")
	if err != nil {
		return err
	}

	_, err = r.routingClient.Exec("CREATE EXTENSION pgRouting")

	return err
}

func (r *PgRouting) Close() error {
	err := r.routingClient.Close()
	if err != nil {
		return err
	}

	err = r.auxClient.Close()
	return err
}
