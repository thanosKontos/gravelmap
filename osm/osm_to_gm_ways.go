package osm

import (
	"github.com/thanosKontos/gravelmap"
)

type osm2GmWays struct {
	ways          map[int]map[int]gravelmap.EvaluatedWay
	nodeDB        gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd      gravelmap.GmNodeReader
	costEvaluator gravelmap.CostEvaluator
}

func NewOsm2GmWays(nodeDB gravelmap.Osm2GmNodeReaderWriter, gmNodeRd gravelmap.GmNodeReader, costEvaluator gravelmap.CostEvaluator) *osm2GmWays {
	ways := make(map[int]map[int]gravelmap.EvaluatedWay)

	return &osm2GmWays{
		nodeDB:        nodeDB,
		gmNodeRd:      gmNodeRd,
		ways:          ways,
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
		} else if i == len(osmNodeIds)-1 {
			points := o.getWayPoints(wayNodeIds)
			evaluation := o.costEvaluator.Evaluate(points.points, tags)

			if _, ok := o.ways[prevEdge]; !ok {
				o.ways[prevEdge] = map[int]gravelmap.EvaluatedWay{}
			}
			if _, ok := o.ways[node.Id]; !ok {
				o.ways[node.Id] = map[int]gravelmap.EvaluatedWay{}
			}

			if existingEvaluatedWay, ok := o.ways[prevEdge][node.Id]; ok {
				if existingEvaluatedWay.Cost > evaluation.WayCost.Normal {
					// There is another from/to way. Keep only the lower cost one
					o.ways[prevEdge][node.Id] = gravelmap.EvaluatedWay{
						Points:        points.points,
						Tags:          tags,
						Distance:      evaluation.Distance,
						WayType:       evaluation.WayType,
						ElevationInfo: evaluation.ElevationEvaluation.Normal,
						Cost:          evaluation.WayCost.Normal,
					}
				}
			} else {
				o.ways[prevEdge][node.Id] = gravelmap.EvaluatedWay{
					Points:        points.points,
					Tags:          tags,
					Distance:      evaluation.Distance,
					WayType:       evaluation.WayType,
					ElevationInfo: evaluation.ElevationEvaluation.Normal,
					Cost:          evaluation.WayCost.Normal,
				}
			}

			if existingEvaluatedWay, ok := o.ways[node.Id][prevEdge]; ok {
				if existingEvaluatedWay.Cost > evaluation.WayCost.Reverse {
					// There is another to/from way. Keep only the lower cost one
					o.ways[node.Id][prevEdge] = gravelmap.EvaluatedWay{
						Points:        points.reverse,
						Tags:          tags,
						Distance:      evaluation.Distance,
						WayType:       evaluation.WayType,
						ElevationInfo: evaluation.ElevationEvaluation.Reverse,
						Cost:          evaluation.WayCost.Reverse,
					}
				}
			} else {
				o.ways[node.Id][prevEdge] = gravelmap.EvaluatedWay{
					Points:        points.reverse,
					Tags:          tags,
					Distance:      evaluation.Distance,
					WayType:       evaluation.WayType,
					ElevationInfo: evaluation.ElevationEvaluation.Reverse,
					Cost:          evaluation.WayCost.Reverse,
				}
			}

			wayNodeIds = []int{prevEdge}
		} else {
			if node.Occurrences > 1 {
				points := o.getWayPoints(wayNodeIds)
				evaluation := o.costEvaluator.Evaluate(points.points, tags)

				if _, ok := o.ways[prevEdge]; !ok {
					o.ways[prevEdge] = map[int]gravelmap.EvaluatedWay{}
				}
				if _, ok := o.ways[node.Id]; !ok {
					o.ways[node.Id] = map[int]gravelmap.EvaluatedWay{}
				}

				if existingEvaluatedWay, ok := o.ways[prevEdge][node.Id]; ok {
					if existingEvaluatedWay.Cost > evaluation.WayCost.Normal {
						// There is another from/to way. Keep only the lower cost one
						o.ways[prevEdge][node.Id] = gravelmap.EvaluatedWay{
							Points:        points.points,
							Tags:          tags,
							Distance:      evaluation.Distance,
							WayType:       evaluation.WayType,
							ElevationInfo: evaluation.ElevationEvaluation.Normal,
							Cost:          evaluation.WayCost.Normal,
						}
					}
				} else {
					o.ways[prevEdge][node.Id] = gravelmap.EvaluatedWay{
						Points:        points.points,
						Tags:          tags,
						Distance:      evaluation.Distance,
						WayType:       evaluation.WayType,
						ElevationInfo: evaluation.ElevationEvaluation.Normal,
						Cost:          evaluation.WayCost.Normal,
					}
				}

				if existingEvaluatedWay, ok := o.ways[node.Id][prevEdge]; ok {
					if existingEvaluatedWay.Cost > evaluation.WayCost.Reverse {
						// There is another to/from way. Keep only the lower cost one
						o.ways[node.Id][prevEdge] = gravelmap.EvaluatedWay{
							Points:        points.reverse,
							Tags:          tags,
							Distance:      evaluation.Distance,
							WayType:       evaluation.WayType,
							ElevationInfo: evaluation.ElevationEvaluation.Reverse,
							Cost:          evaluation.WayCost.Reverse,
						}
					}
				} else {
					o.ways[node.Id][prevEdge] = gravelmap.EvaluatedWay{
						Points:        points.reverse,
						Tags:          tags,
						Distance:      evaluation.Distance,
						WayType:       evaluation.WayType,
						ElevationInfo: evaluation.ElevationEvaluation.Reverse,
						Cost:          evaluation.WayCost.Reverse,
					}
				}

				prevEdge = node.Id
				wayNodeIds = []int{prevEdge}
			}
		}
	}
}

func (o *osm2GmWays) Get() map[int]map[int]gravelmap.EvaluatedWay {
	return o.ways
}

type wayPoints struct {
	points  []gravelmap.Point
	reverse []gravelmap.Point
}

func (o *osm2GmWays) getWayPoints(wayGmNds []int) wayPoints {
	var pts []gravelmap.Point
	var revPts []gravelmap.Point

	for i := len(wayGmNds) - 1; i >= 0; i-- {
		node, _ := o.gmNodeRd.Read(wayGmNds[i])
		revPts = append(revPts, node.Point)
	}

	for _, ndID := range wayGmNds {
		node, _ := o.gmNodeRd.Read(ndID)
		pts = append(pts, node.Point)
	}

	return wayPoints{pts, revPts}
}
