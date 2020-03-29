package node

import (
	"github.com/thanosKontos/gravelmap"
)

type NodeMap map[int64]*gravelmap.ConnectionNode

func NewOsm2GmNodeMemoryStore() NodeMap {
	return make(NodeMap)
}

func (nm NodeMap) Write(osmNdID int64, gm *gravelmap.ConnectionNode) error {
	nm[osmNdID] = gm

	return nil
}

func (nm NodeMap) Read(osmNdID int64) *gravelmap.ConnectionNode {
	if val, ok := nm[osmNdID]; ok {
		return val
	} else {
		return nil
	}
}
