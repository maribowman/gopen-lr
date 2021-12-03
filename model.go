package gopen_lr

const (
	/*
		integer values of location type flags in 4 bits of the first byte:
		point flag (bit 5), area flag (bit 6 and 4), and attributes flag (3).
	*/
	LineLocation               = 1
	GeoCoordinateLocation      = 4
	PointAlongLineLocation     = 5
	PoiWithAccessPointLocation = 5
	CircleLocation             = 0
	RectangleLocation          = 8
	GridLocation               = 8
	PolygonLocation            = 2
	ClosedLineLocation         = 11
)

const (
	// DecaMicroDegFactor -> accuracy of wgs84 nodes
	DecaMicroDegFactor = 100000.0
	// DistancePerInterval -> distance interval DNP
	DistancePerInterval = 58.6
	// BearSector -> size of bearing sector in degrees
	BearSector = 11.25
)

type OpenLRDecoder interface {
	DecodeBase64Encoded(data string) (LocationReference, error)
	DecodeBinary(data []byte) (LocationReference, error)
}

type LocationReference interface {
}

type LocationReferencePoint struct {
	// wgs84 geo-coordinates
	Lon, Lat float64
	// functional road class -> describes the road classification
	FRC int
	// form of way -> describes the physical road type
	FOW int
	// sector of bearing -> bearing is the angle between the direction to a point in the network and a reference direction (here: the true North)
	Bear int
	// lowest FCR to next point -> used to limit scanning alternatives
	LFRCNP int
	// distance to next LR-point -> distance in meters
	DNP int
	// direction of point location -> 0: not applicable (default), 1: A to B, 2: B to A, 3: both ways
	//Orientation int
	// side of road -> 0: not applicable (default), 1: right side, 2: left side, 3: both sides
	//SOR int
}

type LineLocationReference struct {
	Points []LocationReferencePoint
	// positive offset -> distance in meters from start node to actual start location
	POffs float64
	// negative offset -> distance in meters from actual end location to end node
	NOffs float64
}
