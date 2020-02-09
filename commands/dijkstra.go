package commands

import (
	"log"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

// dijkstraCommand defines the version command.
func dijkstraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dijkstra",
		Short: "a dijkstra test",
		Long:  "a dijkstra test",
		Run: func(cmd *cobra.Command, args []string) {
			graph := dijkstra.NewGraph()

			//Add the verticies
			graph.AddVertex(1111)
			graph.AddVertex(1112)
			graph.AddVertex(1113)
			graph.AddVertex(1114)
			graph.AddVertex(1115)

			//Add the arcs
			graph.AddArc(1111,1112,4)
			graph.AddArc(1111,1115,3)
			graph.AddArc(1111,1114,10)
			graph.AddArc(1112,1113,4)
			graph.AddArc(1113,1114,4)
			graph.AddArc(1115,1114,5)

			best, err := graph.Shortest(1111,1114)
			if err != nil{
				log.Fatal(err)
			}
			fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)
		},
	}
}
