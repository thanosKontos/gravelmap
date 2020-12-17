package route

import (
	"github.com/thanosKontos/gravelmap"
)

type router struct {
	edgeFinder  gravelmap.EdgeFinder
	graph       gravelmap.ShortestFinder
	edgeReader  gravelmap.EdgeReader
	pathDecoder gravelmap.Decoder
}

func NewGmRouter(edgeFinder gravelmap.EdgeFinder, graph gravelmap.ShortestFinder, edgeReader gravelmap.EdgeReader, pathDecoder gravelmap.Decoder) *router {
	return &router{
		edgeFinder:  edgeFinder,
		graph:       graph,
		edgeReader:  edgeReader,
		pathDecoder: pathDecoder,
	}
}

func (r *router) Route(ptFrom, ptTo gravelmap.Point) ([]gravelmap.RoutingLeg, error) {
	edgeFrom, err := r.edgeFinder.FindClosest(ptFrom)
	if err != nil {
		return []gravelmap.RoutingLeg{}, err
	}

	edgeTo, err := r.edgeFinder.FindClosest(ptTo)
	if err != nil {
		return []gravelmap.RoutingLeg{}, err
	}

	best, err := r.graph.FindShortest(int(edgeFrom), int(edgeTo))
	if err != nil {
		return []gravelmap.RoutingLeg{}, err
	}

	var edges []gravelmap.Edge
	var prevNodeID = 0
	for i, curNodeID := range best.Path {
		if i == 0 {
			prevNodeID = curNodeID
			continue
		}

		edges = append(edges, gravelmap.Edge{NodeFrom: int32(prevNodeID), NodeTo: int32(curNodeID)})
		prevNodeID = curNodeID
	}

	var routingLegs []gravelmap.RoutingLeg
	presentableWays, _ := r.edgeReader.Read(edges)
	for _, pWay := range presentableWays {
		latLngs := r.pathDecoder.Decode(pWay.Polyline)

		var rlEle *gravelmap.ElevationRange
		if pWay.ElevFrom != 0 && pWay.ElevTo != 0 {
			rlEle = &gravelmap.ElevationRange{
				Start: float64(pWay.ElevFrom),
				End:   float64(pWay.ElevTo),
			}
		}

		wayType := "paved"
		if pWay.SurfaceType == gravelmap.WayTypeUnpaved {
			wayType = "unpaved"
		}
		if pWay.SurfaceType == gravelmap.WayTypePath {
			wayType = "path"
		}

		routingLeg := gravelmap.RoutingLeg{
			Coordinates: latLngs,
			Length:      float64(pWay.Distance),
			WayType:     wayType,
			Elevation:   rlEle,
			OsmID:       pWay.OsmID,
		}

		routingLegs = append(routingLegs, routingLeg)
	}

	return routingLegs, nil
}
