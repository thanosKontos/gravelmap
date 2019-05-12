package routing

import (
	"errors"
	"fmt"
	"os/exec"
)

type Osm2PgRouting struct {
	user         string
	pass         string
	dbname       string
	port         string
	osmFilename  string
	costFilename string
}

func NewOsm2PgRouting(DBUser, DBPass, DBName, DBPort, osmFilename, costFilename string) *Osm2PgRouting {
	return &Osm2PgRouting{
		user:         DBUser,
		pass:         DBPass,
		dbname:       DBName,
		port:         DBPort,
		costFilename: costFilename,
		osmFilename:  osmFilename,
	}
}

func (r *Osm2PgRouting) Import() error {
	cmd := exec.Command("osm2pgrouting", "-c", r.costFilename, "-p", r.port, "-d", r.dbname, "-f", r.osmFilename, "-U", r.user, "-W", r.pass)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	return nil
}

func (r *Osm2PgRouting) Close() error {
	return nil
}
