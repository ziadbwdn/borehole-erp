package geo_converter

import (
    "boreholedata-ms/internal/exception"
    "boreholedata-ms/internal/utils"
    "bytes"
    "strings"
)

// Decimal wraps utils.GormDecimal, inheriting Scan/Value automatically.
type Decimal struct {
    utils.GormDecimal
}

// Factory – create from string without touching old code.
func New(s string) (*Decimal, *exception.AppError) {
    gd, err := utils.StringToGormDecimal(s)
    if err != nil {
        return nil, err
    }
    return &Decimal{GormDecimal: *gd}, nil
}

/* ---------- Extra behaviours ONLY this wrapper adds ---------- */

// Tell GORM the default column definition (override per‐field via tag if needed)
func (Decimal) GormDataType() string { return "decimal(10,6)" } // or 8,2 etc.

// JSON out as `"693.51"`
func (d Decimal) MarshalJSON() ([]byte, error) {
    return []byte(`"` + d.Internal.Value + `"`), nil
}

// JSON in from `"693.51"`
func (d *Decimal) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(bytes.TrimSpace(b)), `"`)
    if s == "" {
        d.Internal.Value = "0"
        return nil
    }
    if _, err := utils.StringToGormDecimal(s); err != nil {
        return err
    }
    d.Internal.Value = s
    return nil
}