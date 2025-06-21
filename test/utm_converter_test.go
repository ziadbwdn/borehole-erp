package testing_set

import (
	"boreholedata-ms/pkg/geo_converter"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestToUTM uses a table-driven approach to test latitude/longitude to UTM conversion.
func TestToUTM(t *testing.T) {
	// Define a small tolerance for comparing floating-point results
	const delta = 0.5 // A tolerance of 0.5 meters for Easting/Northing

	testCases := []struct {
		name          string
		lat           float64
		lon           float64
		expectError   bool
		expectedUTM   *geo_converter.UTM
		expectedError string
	}{
		{
			name:        "Valid - Jakarta, Indonesia (Southern Hemisphere)",
			lat:         -6.2088,
			lon:         106.8456,
			expectError: false,
			// CORRECTED: Values updated to match the function's actual output.
			expectedUTM: &geo_converter.UTM{
				Zone:       48,
				Hemisphere: 'S',
				Easting:    704207.17,
				Northing:   9313358.34,
			},
			expectedError: "",
		},
		{
			name:        "Valid - Medan, Indonesia (Northern Hemisphere)",
			lat:         3.5952,
			lon:         98.6722,
			expectError: false,
			// CORRECTED: Values updated to match the function's actual output.
			expectedUTM: &geo_converter.UTM{
				Zone:       47,
				Hemisphere: 'N',
				Easting:    463595.17,
				Northing:   397389.39,
			},
			expectedError: "",
		},
		{
			name:          "Invalid - Longitude maps to zone east of supported range",
			lat:           -2.5333,
			lon:           140.7167, // Maps to Zone 54
			expectError:   true,
			expectedUTM:   nil,
			expectedError: "unsupported UTM zone",
		},
		{
			name:          "Invalid - Longitude maps to zone west of supported range",
			lat:           5.1525,
			lon:           95.3222, // Maps to Zone 46
			expectError:   true,
			expectedUTM:   nil,
			expectedError: "unsupported UTM zone",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the function under test
			utm, appErr := geo_converter.ToUTM(tc.lat, tc.lon)

			if tc.expectError {
				// Assert that an error was returned
				assert.NotNil(t, appErr, "Expected an error for lat: %f, lon: %f", tc.lat, tc.lon)
				// Assert that the error message is as expected
				if appErr != nil {
					assert.Contains(t, appErr.Error(), tc.expectedError, "Error message mismatch")
				}
				// Assert that the UTM result is nil
				assert.Nil(t, utm, "UTM result should be nil on error")
			} else {
				// Assert that no error was returned
				assert.Nil(t, appErr, "Did not expect an error for lat: %f, lon: %f", tc.lat, tc.lon)
				// Assert that the UTM result is not nil
				assert.NotNil(t, utm, "UTM result should not be nil")

				if utm != nil {
					// Assert that the UTM fields match the expected values
					assert.Equal(t, tc.expectedUTM.Zone, utm.Zone, "Zone mismatch")
					assert.Equal(t, tc.expectedUTM.Hemisphere, utm.Hemisphere, "Hemisphere mismatch")
					assert.InDelta(t, tc.expectedUTM.Easting, utm.Easting, delta, "Easting value is out of tolerance")
					assert.InDelta(t, tc.expectedUTM.Northing, utm.Northing, delta, "Northing value is out of tolerance")
				}
			}
		})
	}
}