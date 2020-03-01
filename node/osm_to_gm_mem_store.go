package node

import (
	"github.com/thanosKontos/gravelmap"
)

type NodeMap map[int64]*gravelmap.Node

func NewOsm2GmNodeMemoryStore() NodeMap {
	return make(NodeMap)
}

func (nm NodeMap) Write(osmNdID int64, gm *gravelmap.Node) error {
	nm[osmNdID] = gm

	return nil
}

func (nm NodeMap) Read(osmNdID int64) *gravelmap.Node {
	if val, ok := nm[osmNdID]; ok {
		return val
	} else {
		return nil
	}
}
