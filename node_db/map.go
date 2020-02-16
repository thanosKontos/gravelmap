package node_db

import "github.com/thanosKontos/gravelmap"

type NodeMap map[int64]*gravelmap.NodeOsm2GM

func NewNodeMapDB() NodeMap {
	return make(NodeMap)
}

func (nm NodeMap) Write(gm *gravelmap.NodeOsm2GM) error {
	nm[gm.OldID] = gm

	return nil
}

func (nm NodeMap) Read(osmNdID int64) *gravelmap.NodeOsm2GM {
	if val, ok := nm[osmNdID]; ok {
		return val
	} else {
		return nil
	}
}
