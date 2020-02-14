package commands

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

// searchCommand searches the graph.
func searchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "search test",
		Long:  "search test",
		Run: func(cmd *cobra.Command, args []string) {
			graph := dijkstra.NewGraph()

			// open data file
			dataFile, err := os.Open("_files/graph.gob")

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dataDecoder := gob.NewDecoder(dataFile)
			err = dataDecoder.Decode(&graph)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dataFile.Close()

			best, err := graph.Shortest(2173, 2201)

			fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)
		},
	}
}
