package matrix

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrInvalidGeoURI represents an error in parsing geo URI.
	ErrInvalidGeoURI = errors.New("invalid geo URI")
	// ErrAltitudeNotPresent is returned when altitude is requested but not present.
	ErrAltitudeNotPresent = errors.New("altitude not present")
)

// URL is a URI which is likely to be MXC URI.
type URL string

// TODO Maybe provide a method to convert it to a normal HTTP.

// GeoURI is a geographic URI.
type GeoURI string

// Parse returns the lat, long and altitude (if present).
// This only implements simple parsing. For more details, use the specific functions.
func (g GeoURI) Parse() (float64, float64, *float64, error) {
	if !strings.HasPrefix(string(g), "geo:") || g == "geo:" {
		return 0, 0, nil, ErrInvalidGeoURI
	}
	split := strings.Split(strings.Split(string(g)[4:], ";")[0], ",")
	if len(split) < 2 || len(split) > 3 {
		return 0, 0, nil, ErrInvalidGeoURI
	}
	parsed := make([]float64, len(split))
	for k, v := range split {
		var err error
		parsed[k], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, 0, nil, ErrInvalidGeoURI
		}
	}
	//nolint:gomnd // 2 is present in longtitude and latitude.
	if len(parsed) == 2 {
		return parsed[0], parsed[1], nil, nil
	}
	return parsed[0], parsed[1], &parsed[2], nil
}

// Lat returns the latitude specified in the GeoURI.
func (g GeoURI) Lat() (float64, error) {
	lat, _, _, err := g.Parse()
	return lat, err
}

// Long returns the longitude specified in the GeoURI.
func (g GeoURI) Long() (float64, error) {
	_, long, _, err := g.Parse()
	return long, err
}

// Altitude returns the altitude specified in the GeoURI.
func (g GeoURI) Altitude() (float64, error) {
	_, _, alt, err := g.Parse()
	if alt == nil {
		return 0, ErrAltitudeNotPresent
	}
	return *alt, err
}
