package osm

import (
	"github.com/thanosKontos/gravelmap"
)

type osm2GmWays struct {
	ways map[int][]gravelmap.WayTo
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd gravelmap.GmNodeReader
	costEvaluator gravelmap.CostEvaluator
}

func NewOsm2GmWays(nodeDB gravelmap.Osm2GmNodeReaderWriter, gmNodeRd gravelmap.GmNodeReader, costEvaluator gravelmap.CostEvaluator) *osm2GmWays {
	ways := make(map[int][]gravelmap.WayTo)

	return &osm2GmWays{
		nodeDB: nodeDB,
		gmNodeRd: gmNodeRd,
		ways: ways,
		costEvaluator: costEvaluator,
	}
}

func (o *osm2GmWays) Add(osmNodeIds []int64, tags map[string]string) {
	prevEdge := 0
	var wayNodeIds []int
	var nodes []gravelmap.Node
	for i, osmNdID := range osmNodeIds {
		node := o.nodeDB.Read(osmNdID)
		nodes = append(nodes, *node)

		wayNodeIds = append(wayNodeIds, node.Id)

		if i == 0 {
			prevEdge = node.Id
		} else if i == len(osmNodeIds) - 1 {
			points := o.getWayPoints(wayNodeIds, false)
			reversePoints := o.getWayPoints(wayNodeIds, true)
			evaluation := o.costEvaluator.Evaluate(points, tags)

			o.ways[prevEdge] = append(o.ways[prevEdge], gravelmap.WayTo{
				NdTo: node.Id,
				Points: points,
				Tags: tags,
				Distance: evaluation.Distance,
				WayType: evaluation.WayType,
				Grade: evaluation.Grade,
				ElevationStart: evaluation.ElevationStart,
				ElevationEnd: evaluation.ElevationEnd,
				Cost: evaluation.Cost,
			})

			o.ways[node.Id] = append(o.ways[node.Id], gravelmap.WayTo{
				NdTo: prevEdge,
				Points: reversePoints,
				Tags: tags,
				Distance: evaluation.Distance,
				WayType: evaluation.WayType,
				Grade: evaluation.Grade,
				ElevationStart: evaluation.ElevationStart,
				ElevationEnd: evaluation.ElevationEnd,
				Cost: evaluation.ReverseCost,
			})

			wayNodeIds = []int{prevEdge}
		} else {
			if node.Occurrences > 1 {
				points := o.getWayPoints(wayNodeIds, false)
				reversePoints := o.getWayPoints(wayNodeIds, true)
				evaluation := o.costEvaluator.Evaluate(points, tags)

				o.ways[prevEdge] = append(o.ways[prevEdge], gravelmap.WayTo{
					NdTo: node.Id,
					Points: points,
					Tags: tags,
					Distance: evaluation.Distance,
					WayType: evaluation.WayType,
					Grade: evaluation.Grade,
					ElevationStart: evaluation.ElevationStart,
					ElevationEnd: evaluation.ElevationEnd,
					Cost: evaluation.Cost,
				})

				o.ways[node.Id] = append(o.ways[node.Id], gravelmap.WayTo{
					NdTo: prevEdge,
					Points: reversePoints,
					Tags: tags,
					WayType: evaluation.WayType,
					Grade: evaluation.Grade,
					ElevationStart: evaluation.ElevationStart,
					ElevationEnd: evaluation.ElevationEnd,
					Cost: evaluation.ReverseCost,
				})

				prevEdge = node.Id
				wayNodeIds = []int{prevEdge}
			}
		}
	}
}

func (o *osm2GmWays) Get() map[int][]gravelmap.WayTo {
	return o.ways
}

func (o *osm2GmWays) getWayPoints(wayGmNds []int, reverse bool) []gravelmap.Point {
	var points []gravelmap.Point

	if reverse {
		for i := len(wayGmNds)-1; i >= 0; i-- {
			gmNode, _ := o.gmNodeRd.Read(wayGmNds[i])
			points = append(points, gmNode.Point)
		}
	} else {
		for _, gmNdID := range wayGmNds {
			gmNode, _ := o.gmNodeRd.Read(gmNdID)
			points = append(points, gmNode.Point)
		}
	}

	return points
}
