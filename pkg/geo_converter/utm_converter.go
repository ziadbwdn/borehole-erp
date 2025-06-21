package geo_converter

import (
    "boreholedata-ms/internal/exception"
    "math"
    "strconv"
)

// ---------------------------------------------------------------------------
// 2. UTM conversion with domain constraints + AppError returns
// ---------------------------------------------------------------------------
type UTM struct {
    Zone       int
    Hemisphere rune
    Easting    float64
    Northing   float64
}

// ToUTM converts lat/long (Decimal) → UTM (WGS-84). Valid only for Indonesian zones 47–52; 
func ToUTM(lat, lon float64) (*UTM, *exception.AppError) {
    zone := int(math.Floor((lon+180)/6.0)) + 1
    if zone < 47 || zone > 52 {
        return nil, exception.NewValidationError("longitude maps to unsupported UTM zone", strconv.Itoa(zone))
    }
    hem := 'N'
    if lat < 0 {
        hem = 'S'
    }
    const (
        a  = 6378137.0
        f  = 1 / 298.257223563
        k0 = 0.9996
    )
    eSq := f * (2 - f) // eccentricity squared
    latRad := lat * math.Pi / 180
    lonRad := lon * math.Pi / 180
    lon0 := float64(zone-1)*6 - 180 + 3
    lon0Rad := lon0 * math.Pi / 180

    NVal := a / math.Sqrt(1-eSq*math.Sin(latRad)*math.Sin(latRad))
    T := math.Pow(math.Tan(latRad), 2)
    C := eSq/(1-eSq) * math.Pow(math.Cos(latRad), 2)
    A := math.Cos(latRad) * (lonRad - lon0Rad)

    M := a * ((1 - eSq/4 - 3*eSq*eSq/64 - 5*eSq*eSq*eSq/256) * latRad -
        (3*eSq/8 + 3*eSq*eSq/32 + 45*eSq*eSq*eSq/1024)*math.Sin(2*latRad) +
        (15*eSq*eSq/256 + 45*eSq*eSq*eSq/1024)*math.Sin(4*latRad) -
        (35*eSq*eSq*eSq/3072)*math.Sin(6*latRad))

    easting := k0*NVal*(A+(1-T+C)*math.Pow(A, 3)/6+
        (5-18*T+T*T+72*C-58*(eSq/(1-eSq)))*math.Pow(A, 5)/120) + 500000

    northing := k0*(M+NVal*math.Tan(latRad)*(math.Pow(A, 2)/2+
        (5-T+9*C+4*C*C)*math.Pow(A, 4)/24+
        (61-58*T+T*T+600*C-330*(eSq/(1-eSq)))*math.Pow(A, 6)/720))
    if lat < 0 {
        northing += 10000000
    }

    return &UTM{
        Zone:       zone,
        Hemisphere: hem,
        Easting:    easting,
        Northing:   northing,
    }, nil
}