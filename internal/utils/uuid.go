package utils

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type BinaryUUID uuid.UUID

// NewBinaryUUID generates a new random BinaryUUID.
func NewBinaryUUID() BinaryUUID {
	return BinaryUUID(uuid.New())
}

// ParseBinaryUUID parses a string into a BinaryUUID.
// It returns an error if the string is not a valid UUID.
// This function MUST start with an uppercase 'P' to be exported.
func ParseBinaryUUID(s string) (BinaryUUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return BinaryUUID{}, fmt.Errorf("failed to parse UUID string '%s': %w", s, err)
	}
	return BinaryUUID(u), nil
}

// String returns the string representation of the BinaryUUID.
func (b BinaryUUID) String() string {
	return uuid.UUID(b).String()
}

// GormDataType returns the GORM data type for BinaryUUID.
func (b BinaryUUID) GormDataType() string {
	return "binary(16)"
}

// MarshalJSON implements the json.Marshaler interface for BinaryUUID.
func (b BinaryUUID) MarshalJSON() ([]byte, error) {
	s := uuid.UUID(b)
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for BinaryUUID.
func (b *BinaryUUID) UnmarshalJSON(by []byte) error {
	// Remove quotes from the JSON string
	s, err := uuid.ParseBytes(by[1 : len(by)-1])
	if err != nil {
		return fmt.Errorf("failed to unmarshal UUID from JSON: %w", err)
	}
	*b = BinaryUUID(s)
	return nil
}

// Value implements the driver.Valuer interface for BinaryUUID.
func (b BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(b).MarshalBinary()
}

// Scan implements the sql.Scanner interface for BinaryUUID.
func (b *BinaryUUID) Scan(value interface{}) error {
	var u uuid.UUID
	if err := u.Scan(value); err != nil {
		return fmt.Errorf("failed to scan UUID from database: %w", err)
	}
	*b = BinaryUUID(u)
	return nil
}
