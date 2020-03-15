package osm

import (
	"github.com/thanosKontos/gravelmap"
)

type osm2GmWays struct {
	ways          map[int][]gravelmap.WayTo
	nodeDB        gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd      gravelmap.GmNodeReader
	costEvaluator gravelmap.CostEvaluator
	pathSimplifier gravelmap.PathSimplifier
}

func NewOsm2GmWays(
	nodeDB gravelmap.Osm2GmNodeReaderWriter,
	gmNodeRd gravelmap.GmNodeReader,
	costEvaluator gravelmap.CostEvaluator,
	pathSimplifier gravelmap.PathSimplifier,
	) *osm2GmWays {
	ways := make(map[int][]gravelmap.WayTo)

	return &osm2GmWays{
		nodeDB:        nodeDB,
		gmNodeRd:      gmNodeRd,
		ways:          ways,
		costEvaluator: costEvaluator,
		pathSimplifier: pathSimplifier,
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
		} else if i == len(osmNodeIds)-1 {
			points := o.getWayPoints(wayNodeIds)
			evaluation := o.costEvaluator.Evaluate(points.points, tags)

			o.ways[prevEdge] = append(o.ways[prevEdge], gravelmap.WayTo{
				NdTo:          node.Id,
				Points:        points.points,
				Tags:          tags,
				Distance:      evaluation.Distance,
				WayType:       evaluation.WayType,
				ElevationInfo: evaluation.ElevationEvaluation.Normal,
				Cost:          evaluation.WayCost.Normal,
			})

			o.ways[node.Id] = append(o.ways[node.Id], gravelmap.WayTo{
				NdTo:          prevEdge,
				Points:        points.reverse,
				Tags:          tags,
				Distance:      evaluation.Distance,
				WayType:       evaluation.WayType,
				ElevationInfo: evaluation.ElevationEvaluation.Reverse,
				Cost:          evaluation.WayCost.Reverse,
			})

			wayNodeIds = []int{prevEdge}
		} else {
			if node.Occurrences > 1 {
				points := o.getWayPoints(wayNodeIds)
				evaluation := o.costEvaluator.Evaluate(points.points, tags)

				o.ways[prevEdge] = append(o.ways[prevEdge], gravelmap.WayTo{
					NdTo:          node.Id,
					Points:        points.points,
					Tags:          tags,
					Distance:      evaluation.Distance,
					WayType:       evaluation.WayType,
					ElevationInfo: evaluation.ElevationEvaluation.Normal,
					Cost:          evaluation.WayCost.Normal,
				})

				o.ways[node.Id] = append(o.ways[node.Id], gravelmap.WayTo{
					NdTo:          prevEdge,
					Points:        points.reverse,
					Tags:          tags,
					WayType:       evaluation.WayType,
					ElevationInfo: evaluation.ElevationEvaluation.Reverse,
					Cost:          evaluation.WayCost.Reverse,
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

type wayPoints struct {
	points  []gravelmap.Point
	reverse []gravelmap.Point
}

func (o *osm2GmWays) getWayPoints(wayGmNds []int) wayPoints {
	var pts []gravelmap.Point
	var revPts []gravelmap.Point

	for _, ndID := range wayGmNds {
		node, _ := o.gmNodeRd.Read(ndID)
		pts = append(pts, node.Point)
	}

	pts = o.pathSimplifier.Simplify(pts)

	for i := len(pts) - 1; i >= 0; i-- {
		revPts = append(revPts, pts[i])
	}

	return wayPoints{pts, revPts}
}
