package gopen_lr

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeBinary(t *testing.T) {
	// given
	decoder := NewOpenLRDecoder()
	tables := []struct {
		baseEncoded string
		expected    LineLocationReference
	}{
		{
			baseEncoded: "CwRn6iRZ8RtnEwGjAlcLPDgHkAX2CxM=",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{6.195806264877319, 51.11905217170715, 3, 3, 84, 3, 1143}, {6.19999626487732, 51.12504217170715, 1, 3, 321, 1, 3311}, {6.2193562648773195, 51.14030217170715, 1, 3, 219, 7, 0}}},
		},
		{
			baseEncoded: "CwUoqSSdQxJfDQC+ArsSDw==",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{7.254592180252075, 51.48885369300842, 2, 2, 354, 2, 791}, {7.256492180252075, 51.495843693008425, 2, 2, 174, 7, 0}}},
		},
		{
			baseEncoded: "CwSzvySU+BtrShdw/l8bFg==",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{6.612364053726196, 51.44329905509949, 3, 3, 129, 3, 4366}, {6.672364053726196, 51.439129055099485, 3, 3, 253, 7, 0}}},
		},
		{
			baseEncoded: "CwS+qySUNhtmKApfBTAbZBAEVgHpG2MWB2L/dRI4Aw==",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{6.6723597049713135, 51.439136266708374, 3, 3, 73, 3, 2373}, {6.698909704971314, 51.452416266708376, 3, 3, 51, 3, 967}, {6.710009704971314, 51.45730626670838, 3, 3, 39, 3, 1319}, {6.728909704971314, 51.45591626670838, 2, 2, 276, 7, 0}}, NOffs: 0.013671875},
		},
		{
			baseEncoded: "Cwqe1SSYshNa67yTFKITWB34Gv1hFAk=",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{14.934979677200317, 51.46376967430115, 2, 3, 298, 2, 13800}, {14.762369677200317, 51.516589674301144, 2, 3, 276, 2, 1729}, {14.742149677200317, 51.509879674301146, 2, 4, 107, 7, 0}}},
		},
		{
			baseEncoded: "CwYa2yYSwRtqAgDr/8QbWXw=",
			expected:    LineLocationReference{Points: []LocationReferencePoint{{8.585010766983032, 53.540507555007935, 3, 3, 118, 3, 147}, {8.587360766983032, 53.539907555007936, 3, 3, 287, 7, 0}}, POffs: 0.486328125},
		},
	}

	for _, table := range tables {
		// when
		actual, err := decoder.DecodeBase64Encoded(table.baseEncoded)
		// then
		assert.NoError(t, err)
		assert.EqualValues(t, table.expected, actual)

		// and when
		olrBytes, byteErr := base64.StdEncoding.DecodeString(table.baseEncoded)
		if byteErr != nil {
			assert.Fail(t, "could not decode base64", byteErr)
		}
		actual, err = decoder.DecodeBinary(olrBytes)
		// then
		assert.NoError(t, err)
		assert.EqualValues(t, table.expected, actual)
	}
}
