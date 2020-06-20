package osm

import (
	"github.com/thanosKontos/gravelmap"
)

type osm2GmWays struct {
	ways           map[int]map[int]gravelmap.Way
	nodeDB         gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd       gravelmap.Osm2LatLngReader
	costEvaluator  gravelmap.CostEvaluator
	pathSimplifier gravelmap.PathSimplifier
}

func NewOsm2GmWays(
	nodeDB gravelmap.Osm2GmNodeReaderWriter,
	gmNodeRd gravelmap.Osm2LatLngReader,
	costEvaluator gravelmap.CostEvaluator,
	pathSimplifier gravelmap.PathSimplifier,
) *osm2GmWays {
	ways := make(map[int]map[int]gravelmap.Way)

	return &osm2GmWays{
		nodeDB:         nodeDB,
		gmNodeRd:       gmNodeRd,
		ways:           ways,
		costEvaluator:  costEvaluator,
		pathSimplifier: pathSimplifier,
	}
}

func (o *osm2GmWays) Add(osmNodeIds []int64, tags map[string]string, osmID int64) {
	prevNodeID := 0
	var wayNodeIds []int
	for i, osmNdID := range osmNodeIds {
		node := o.nodeDB.Read(osmNdID)
		wayNodeIds = append(wayNodeIds, node.ID)

		// First way node
		if i == 0 {
			prevNodeID = node.ID

			continue
		}

		// Edge node with a connection or last node
		if i == len(osmNodeIds)-1 || node.ConnectionCnt > 1 {
			o.AddBackAndForthEdgesToWays(prevNodeID, node.ID, wayNodeIds, tags, osmID)

			prevNodeID = node.ID
			wayNodeIds = []int{prevNodeID}

			continue
		}
	}
}

func (o *osm2GmWays) AddBackAndForthEdgesToWays(edgeNodeFrom, edgeNodeTo int, wayNodeIds []int, tags map[string]string, osmID int64) {
	points := o.getWayPoints(wayNodeIds)
	evaluation := o.costEvaluator.Evaluate(points.points, tags)

	if existingEvaluatedWay, ok := o.ways[edgeNodeFrom][edgeNodeTo]; ok {
		if existingEvaluatedWay.Cost > evaluation.BidirectionalCost.Normal {
			if _, ok := o.ways[edgeNodeFrom]; !ok {
				o.ways[edgeNodeFrom] = map[int]gravelmap.Way{}
			}

			// There is another from/to way. Keep only the lower cost one
			o.ways[edgeNodeFrom][edgeNodeTo] = gravelmap.Way{
				Points:        points.points,
				Tags:          tags,
				Distance:      evaluation.Distance,
				Type:          evaluation.WayType,
				ElevationInfo: evaluation.BidirectionalElevationInfo.Normal,
				Cost:          evaluation.BidirectionalCost.Normal,

				OriginalOsmID: osmID,
			}
		}
	} else {
		if _, ok := o.ways[edgeNodeFrom]; !ok {
			o.ways[edgeNodeFrom] = map[int]gravelmap.Way{}
		}

		o.ways[edgeNodeFrom][edgeNodeTo] = gravelmap.Way{
			Points:        points.points,
			Tags:          tags,
			Distance:      evaluation.Distance,
			Type:          evaluation.WayType,
			ElevationInfo: evaluation.BidirectionalElevationInfo.Normal,
			Cost:          evaluation.BidirectionalCost.Normal,

			OriginalOsmID: osmID,
		}
	}

	if existingEvaluatedWay, ok := o.ways[edgeNodeTo][edgeNodeFrom]; ok {
		if existingEvaluatedWay.Cost > evaluation.BidirectionalCost.Reverse {
			if _, ok := o.ways[edgeNodeTo]; !ok {
				o.ways[edgeNodeTo] = map[int]gravelmap.Way{}
			}

			// There is another to/from way. Keep only the lower cost one
			o.ways[edgeNodeTo][edgeNodeFrom] = gravelmap.Way{
				Points:        points.reverse,
				Tags:          tags,
				Distance:      evaluation.Distance,
				Type:          evaluation.WayType,
				ElevationInfo: evaluation.BidirectionalElevationInfo.Reverse,
				Cost:          evaluation.BidirectionalCost.Reverse,

				OriginalOsmID: osmID,
			}
		}
	} else {
		if _, ok := o.ways[edgeNodeTo]; !ok {
			o.ways[edgeNodeTo] = map[int]gravelmap.Way{}
		}

		o.ways[edgeNodeTo][edgeNodeFrom] = gravelmap.Way{
			Points:        points.reverse,
			Tags:          tags,
			Distance:      evaluation.Distance,
			Type:          evaluation.WayType,
			ElevationInfo: evaluation.BidirectionalElevationInfo.Reverse,
			Cost:          evaluation.BidirectionalCost.Reverse,

			OriginalOsmID: osmID,
		}
	}
}

func (o *osm2GmWays) Get() map[int]map[int]gravelmap.Way {
	return o.ways
}

type wayPoints struct {
	points  []gravelmap.Point
	reverse []gravelmap.Point
}

func (o *osm2GmWays) getWayPoints(wayGmNds []int) wayPoints {
	var pts []gravelmap.Point
	var revPts []gravelmap.Point

	for _, ndID := range wayGmNds {
		pt, _ := o.gmNodeRd.Read(ndID)
		pts = append(pts, pt)
	}

	pts = o.pathSimplifier.Simplify(pts)

	for i := len(pts) - 1; i >= 0; i-- {
		revPts = append(revPts, pts[i])
	}

	return wayPoints{pts, revPts}
}
