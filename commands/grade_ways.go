package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/elevation/osm_grade"
	"github.com/thanosKontos/gravelmap/elevation/srtm_ascii"
	"log"
	"os"
	"strings"
	"sync"
)

const batchSize = 500

type geomRow struct {
	id     int64
	length float64
	geom   string
}

var wg sync.WaitGroup

// createRoutingDataCommand defines the create route command.
func createGradeWaysCommand() *cobra.Command {
	var (
		OSMIDs string
	)

	createGradeWaysCmd := &cobra.Command{
		Use:   "grade-ways",
		Short: "grade ways in route database",
		Long:  "grade ways and fill up elevation cost in route database",
	}

	createGradeWaysCmd.Flags().StringVar(&OSMIDs, "osm_ids", "", "The osm input file.")
	createGradeWaysCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createGradeWaysCmdRun(OSMIDs)
	}

	return createGradeWaysCmd
}

// createGradeWaysCmdRun defines the command run actions.
func createGradeWaysCmdRun(OSMIDs string) error {
	eleFinder, _ := srtm_ascii.NewElevationFinder(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	eleGrader, _ := srtm_ascii.NewElevationGrader(eleFinder)
	osmGrader, err := osm_grade.NewSrtmOsmGrader(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"), eleGrader, logger)
	if err != nil {
		log.Fatalln(err)
	}

	if OSMIDs != "" {
		osmGrader.SetOSMIDs(strings.Split(OSMIDs, ","))
	}

	osmGrader.GradeWays()
	log.Println("Roads graded.")

	return nil
}
