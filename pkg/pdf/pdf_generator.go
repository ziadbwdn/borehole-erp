package pdf

import (
	"boreholedata-ms/internal/models"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PDFGenerator provides methods for generating various PDF reports.
type PDFGenerator struct{}

// NewPDFGenerator creates and returns a new instance of PDFGenerator.
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// GenerateBoreholeLogPDF generates a PDF report for a borehole log.
// It takes station details, lithology logs, lab samples, and UCS results.
func (g *PDFGenerator) GenerateBoreholeLogPDF(
	station *models.Station,
	lithologyLogs []*models.LithologyLog,
	labSamples []*models.LabSample,
	ucsResults []*models.UCSResult,
) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Borehole Log: %s", station.StationName))
	pdf.Ln(12)

	// Station Details
	pdf.SetFont("Arial", "", 10)

	pdf.Cell(0, 7, fmt.Sprintf("Project: %s", station.ProjectID.String()))
	pdf.Ln(5)
	// Corrected: Using GormDecimal.Value for formatting
	pdf.Cell(0, 7, fmt.Sprintf("Location: Lat %s, Lon %s",
		formatDecimalString(station.Latitude.Internal.Value, 6),
		formatDecimalString(station.Longitude.Internal.Value, 6)))
	pdf.Ln(5)
	pdf.Cell(0, 7, fmt.Sprintf("Total Depth: %s m",
		formatDecimalString(station.TotalDepth.Internal.Value, 2)))
	pdf.Ln(5)
	pdf.Cell(0, 7, fmt.Sprintf("Drilling Date: %s",
		station.DrillingDate.Format("2006-01-02")))
	pdf.Ln(5)
	pdf.Cell(0, 7, fmt.Sprintf("Geologist: %s", station.GeologistName))
	pdf.Ln(10)

	// Lithology Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 7, "Lithology Log:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 9)
	for _, lith := range lithologyLogs {
		// Corrected: Use GormDecimal.Value for formatting
		pdf.Cell(0, 5, fmt.Sprintf("  %s - %s m: %s (%s) RQD: %s%%, Recovery: %s%%",
			formatDecimalString(lith.DepthFrom.Internal.Value, 2),
			formatDecimalString(lith.DepthTo.Internal.Value, 2),
			lith.LithologyType,
			truncateString(lith.Description, 60),
			formatDecimalString(lith.RQDPercentage.Internal.Value, 2),
			formatDecimalString(lith.RecoveryPercentage.Internal.Value, 2)))
		pdf.Ln(5)
	}
	pdf.Ln(10)

	// Samples Section
	if len(labSamples) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 7, "Samples:")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 9)
		for _, sample := range labSamples {
			// Corrected: Use GormDecimal.Value for formatting
			pdf.Cell(0, 5, fmt.Sprintf("  %s: %s-%s m (%s) - %s",
				sample.SampleCode,
				formatDecimalString(sample.DepthFrom.Internal.Value, 2),
				formatDecimalString(sample.DepthTo.Internal.Value, 2),
				sample.SampleType,
				sample.SamplingDate.Format("2006-01-02")))
			pdf.Ln(5)
		}
		pdf.Ln(10)
	}

	// UCS Results Section
	if len(ucsResults) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 7, "UCS Test Results:")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 9)
		for _, ucs := range ucsResults {
			// Corrected: Use GormDecimal.Value for formatting
			pdf.Cell(0, 5, fmt.Sprintf("  %s %s - %s (Sample: %s)",
				formatDecimalString(ucs.UCSValue.Internal.Value, 2),
				ucs.Unit,
				ucs.TestMethod,
				ucs.SampleID.String()))
			pdf.Ln(5)
		}
	}

	// Footer
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 6, fmt.Sprintf("Report generated: %s", time.Now().Format("2006-01-02 15:04")))

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF to buffer: %w", err)
	}

	return buf, nil
}

// Helper to format decimal strings
func formatDecimalString(s string, precision int) string {
	if !strings.Contains(s, ".") {
		if precision > 0 {
			return s + "." + strings.Repeat("0", precision)
		}
		return s
	}

	parts := strings.Split(s, ".")
	integer := parts[0]
	decimal := parts[1]

	if len(decimal) > precision {
		decimal = decimal[:precision]
	} else {
		decimal = decimal + strings.Repeat("0", precision-len(decimal))
	}

	return integer + "." + decimal
}

// Helper to format floats (retained for completeness if needed elsewhere, though not used for model decimals here)
func formatFloat(f float64, precision int) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}

// Helper to truncate long strings
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
