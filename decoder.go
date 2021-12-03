package gopen_lr

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/wadey/go-rounding"
	"math"
	"math/big"
)

type decoder struct {
	reader *bytes.Reader
	offset int
	size   int
}

func NewOpenLRDecoder() OpenLRDecoder {
	return &decoder{}
}

func (decoder *decoder) DecodeBase64Encoded(data string) (LocationReference, error) {
	olrBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return LineLocationReference{}, err
	} else {
		return decoder.DecodeBinary(olrBytes)
	}
}

func (decoder *decoder) DecodeBinary(data []byte) (LocationReference, error) {
	decoder.reader = bytes.NewReader(data)
	decoder.offset = -1
	decoder.size = len(data)

	version, locationType := decoder.readStatus()
	if version != 3 {
		return LineLocationReference{}, fmt.Errorf("only version 3 supported, detected version: %d", version)
	}
	switch locationType {
	case LineLocation:
		return decoder.parseLine(), nil
	//case 	api.GeoCoordinateLocation:
	//case	api.PointAlongLineLocation:
	//case	api.PoiWithAccessPointLocation:
	//case	api.CircleLocation:
	//case	api.RectangleLocation:
	//case	api.GridLocation:
	//case	api.PolygonLocation:
	//case	api.ClosedLineLocation:
	default:
		// TODO: impl additional reference types if funny
		return nil, fmt.Errorf("location type %d not supported", locationType)
	}
}

func (decoder *decoder) read(size int) []byte {
	var tempBytes []byte
	for i := 0; i < size; i++ {
		tempByte, _ := decoder.reader.ReadByte()
		tempBytes = append(tempBytes, tempByte)
		decoder.offset++
	}
	return tempBytes
}

func (decoder *decoder) readStatus() (byte, byte) {
	status := decoder.read(1)[0]
	version := status & 7
	locationType := (status >> 3) & 0b1111
	return version, locationType
}

func (decoder *decoder) parseLine() LineLocationReference {
	var points []LocationReferencePoint
	relativePoints := (decoder.size - 9) / 7
	lon, lat := decoder.readCoordinates()
	fow, frc, _, bear, lfrcnp := decoder.readPointAttributes()
	for i := 1; i <= relativePoints; i++ {
		dnp := decoder.readDNP()
		points = append(points, LocationReferencePoint{Lon: lon, Lat: lat, FRC: frc, FOW: fow, Bear: bear, LFRCNP: lfrcnp, DNP: dnp})
		lon, lat = decoder.readRelativeCoordinates(lon, lat)
		fow, frc, _, bear, lfrcnp = decoder.readPointAttributes()
	}
	points = append(points, LocationReferencePoint{Lon: lon, Lat: lat, FRC: frc, FOW: fow, Bear: bear, LFRCNP: 7, DNP: 0})

	pOffs, nOffs := 0.0, 0.0
	if lfrcnp&0b10 > 0 {
		pOffs = decoder.readOffset()
	}
	if lfrcnp&0b01 > 0 {
		nOffs = decoder.readOffset()
	}

	return LineLocationReference{
		Points: points,
		POffs:  pOffs,
		NOffs:  nOffs,
	}
}

func (decoder *decoder) readCoordinates() (float64, float64) {
	// reads absolute coordinates from the buffer (6 bytes)
	lonInt := bytesToInt(decoder.read(3), true)
	latInt := bytesToInt(decoder.read(3), true)
	lon := intToDeg(lonInt)
	lat := intToDeg(latInt)
	return lon, lat
}

func (decoder *decoder) readRelativeCoordinates(prevLon, prevLat float64) (float64, float64) {
	// read coordinates from buffer relative to the previous ones (4 bytes)
	relLonInt := bytesToInt(decoder.read(2), true)
	relLatInt := bytesToInt(decoder.read(2), true)
	lon := prevLon + float64(relLonInt)/DecaMicroDegFactor
	lat := prevLat + float64(relLatInt)/DecaMicroDegFactor
	return lon, lat
}

func (decoder *decoder) readPointAttributes() (int, int, int, int, int) {
	/*
		read point attributes from buffer (2 bytes)
		- FOW (3 bits)
		- FRC (3 bits)
		- reserved for future use, orientation/side of road (2 bits)
		- bear (5 bits)
		- LFRCNP (3 bits)
	*/
	attributeBytes := decoder.read(2)
	fow := attributeBytes[0] & 0b111
	frc := (attributeBytes[0] >> 3) & 0b111
	reserved := (attributeBytes[0] >> 6) & 0b11
	bear := float64(attributeBytes[1]&0b11111)*BearSector + BearSector/2
	lfrcnp := (attributeBytes[1] >> 5) & 0b111
	return int(fow), int(frc), int(reserved), roundHalfUp(bear), int(lfrcnp)
}

func (decoder *decoder) readDNP() int {
	// reads distance to next point from the buffer (1 byte)
	interval := bytesToInt(decoder.read(1), false)
	dnp := (float64(interval) + 0.5) * DistancePerInterval
	return roundHalfUp(dnp)
}

func (decoder *decoder) readOffset() float64 {
	// reads distance to next point from the buffer (1 byte)
	bucketIndex := bytesToInt(decoder.read(1), false)
	offset := (float64(bucketIndex) + 0.5) / 256
	return offset
}

func bytesToInt(rawBytes []byte, signed bool) int {
	// converts big endian bytes to signed/unsigned int
	resolution := len(rawBytes) * 8
	var freshBytes [8]byte
	copy(freshBytes[8-len(rawBytes):], rawBytes)
	value := int(binary.BigEndian.Uint64(freshBytes[:]))
	if signed && (rawBytes[0]>>7) == 1 {
		value -= 1 << resolution
	}
	return value
}

func intToDeg(value int) float64 {
	// converts integer into degree coordinate
	//resolution := 24 // default
	return ((float64(value) - math.Copysign(float64(1), float64(value))*0.5) * 360) / (1 << 24)
}

func roundHalfUp(value float64) int {
	// as per tomtom impl, rounding half up (e.g. 2.5 -> 3)
	roundedValueFloat, ok := rounding.Round(new(big.Rat).SetFloat64(value), 0, rounding.HalfUp).Float64()
	if !ok {
		return 0
	}
	return int(roundedValueFloat)
}
