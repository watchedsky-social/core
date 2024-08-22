package utils

import "github.com/paulmach/orb"

func MergeGeometries(zones []orb.Geometry) orb.MultiPolygon {
	collector := orb.MultiPolygon{}

	for _, zone := range zones {
		internalGeometryMerge(zone, &collector)
	}

	return collector
}

func internalGeometryMerge(zone orb.Geometry, collector *orb.MultiPolygon) {
	switch zone.GeoJSONType() {
	case "Polygon":
		p, _ := zone.(orb.Polygon)
		*collector = append(*collector, p)
	case "MultiPolygon":
		p, _ := zone.(orb.MultiPolygon)
		*collector = append(*collector, p...)
	case "GeometryCollection":
		collection, _ := zone.(orb.Collection)
		for _, coll := range collection {
			internalGeometryMerge(coll, collector)
		}
	}
}
