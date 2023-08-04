package dbHelpers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateULID(t *testing.T) {
	// generate two ULIDs
	ulid1, err1 := GenerateULID()
	ulid2, err2 := GenerateULID()

	// there should be no errors
	assert.NoError(t, err1, "There should be no error generating the first ULID")
	assert.NoError(t, err2, "There should be no error generating the second ULID")

	// the ULIDs should not be nil
	assert.NotEqual(t, uuid.Nil, ulid1, "The first ULID should not be nil")
	assert.NotEqual(t, uuid.Nil, ulid2, "The second ULID should not be nil")

	// the ULIDs should not be equal (this might fail if the function is broken, or with extremely low probability due to random collision)
	assert.NotEqual(t, ulid1, ulid2, "The ULIDs should not be equal")

	// the timestamp portion (first 48 bits / 6 bytes) of the ULIDs should be increasing
	ts1 := ulid1[:6]
	ts2 := ulid2[:6]
	assert.GreaterOrEqual(t, ts2, ts1, "The timestamp portion of the second ULID should be greater than or equal to the first one")
}
