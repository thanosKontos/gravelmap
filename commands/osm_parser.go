package commands

import (
	"github.com/qedus/osmpbf"

	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// dijkstraCommand defines the version command.
func osmParserCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "osm-parse",
		Short: "an osm parse test",
		Long:  "an osm parse test",
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open("/Users/thanoskontos/Downloads/greece_for_routing_from_pbf.osm.pbf")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			d := osmpbf.NewDecoder(f)

			// use more memory from the start, it is faster
			d.SetBufferSize(osmpbf.MaxBlobSize)

			// start decoding with several goroutines, it is faster
			err = d.Start(runtime.GOMAXPROCS(-1))
			if err != nil {
				log.Fatal(err)
			}

			var nc, wc, rc uint64
			for {
				if v, err := d.Decode(); err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				} else {
					switch v := v.(type) {
					case *osmpbf.Node:
						// Process Node v.
						nc++
					case *osmpbf.Way:
						// Process Way v.
						wc++
					case *osmpbf.Relation:
						// Process Relation v.
						rc++
					default:
						log.Fatalf("unknown type %T\n", v)
					}
				}
			}

			fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)
		},
	}
}
