package routing_prepare

import (
	"database/sql"
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"os"

	_ "github.com/lib/pq"
)

type Gravelmap struct {
	routingClient *sql.DB
	auxClient     *sql.DB
	logger        gravelmap.Logger
}

// NewRouting initialize and return an new PgRouting object.
func NewGravelmapPreparer(DBUser, DBPass, DBName, DBPort, AuxDBName string, logger gravelmap.Logger) (*Gravelmap, error) {
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

	return &Gravelmap{
		routingClient: DB,
		auxClient:     auxDB,
		logger:        logger,
	}, nil
}

func (g *Gravelmap) Prepare() error {
	_, err := g.auxClient.Exec(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity  WHERE datname = 'routing'`)
	if err != nil {
		return err
	}

	_, err = g.auxClient.Exec("CREATE DATABASE IF NOT EXISTS " + os.Getenv("DBNAME"))
	if err != nil {
		g.logger.Info(fmt.Sprintf("database not created: %s", err))
	}

	_, err = g.routingClient.Exec("CREATE EXTENSION postGIS")
	if err != nil {
		g.logger.Info(fmt.Sprintf("extension postGIS not created: %s", err))
	}

	_, err = g.routingClient.Exec("CREATE EXTENSION pgRouting")
	if err != nil {
		g.logger.Info(fmt.Sprintf("extension pgRouting not created: %s", err))
	}

	return nil
}

func (g *Gravelmap) Close() error {
	err := g.routingClient.Close()
	if err != nil {
		return err
	}

	err = g.auxClient.Close()
	return err
}
