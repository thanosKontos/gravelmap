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

	defer graphFile.Close()
	return dataEncoder.Encode(graph)
}

func (r *gobRepo) Fetch(name string) (*WeightedBidirectionalGraph, error) {
	g := NewWeightedBidirectionalGraph()
	dataFile, err := os.Open(fmt.Sprintf("%s/graph_%s.gob", r.storageDir, name))
	if err != nil {
		return nil, err
	}

	dataDecoder := gob.NewDecoder(dataFile)
	if err = dataDecoder.Decode(&g); err != nil {
		return nil, err
	}
	err = dataFile.Close()

	return g, err
}
