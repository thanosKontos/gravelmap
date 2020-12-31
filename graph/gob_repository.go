package graph

import (
	"encoding/gob"
	"fmt"
	"os"
)

type gobRepo struct {
	storageDir string
}

// NewWeightedBidirectionalGraph creates objects of type WeightedBidirectionalGraph
func NewGobRepo(storageDir string) *gobRepo {
	return &gobRepo{
		storageDir: storageDir,
	}
}

func (r *gobRepo) Store(graph *WeightedBidirectionalGraph, name string) error {
	graphFile, err := os.Create(fmt.Sprintf("%s/graph_%s.gob", r.storageDir, name))
	if err != nil {
		return err
	}
	dataEncoder := gob.NewEncoder(graphFile)
	err = dataEncoder.Encode(graph)
	if err != nil {
		return err
	}

	return graphFile.Close()
}

func (r *gobRepo) Fetch(name string) (*WeightedBidirectionalGraph, error) {
	g := NewWeightedBidirectionalGraph()
	dataFile, err := os.Open("_files/graph_mtb.gob")
	if err != nil {
		return nil, err
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&g)
	if err != nil {
		return nil, err
	}
	err = dataFile.Close()

	return g, err
}
