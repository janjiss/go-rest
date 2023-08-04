package dbHelpers

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
)

func GenerateULID() (uuid.UUID, error) {
	// Get the current Unix timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// We want 48 bits for the timestamp. However, timestamps are generally
	// 64 bits in Go. So, let's take only the least significant 48 bits.
	timestamp = timestamp & ((1 << 48) - 1)

	// Generate 80 random bits
	randomBits := make([]byte, 10)
	_, err := rand.Read(randomBits)
	if err != nil {
		return uuid.UUID{}, err
	}

	// Create a byte slice of the timestamp
	tsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBytes, uint64(timestamp))

	// Build a 16-byte slice with the timestamp and random bits
	ulidBytes := append(tsBytes[2:], randomBits...) // only take the least significant 6 bytes of the timestamp

	// Create a UUID from the byte slice
	u, err := uuid.FromBytes(ulidBytes)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u, nil
}
