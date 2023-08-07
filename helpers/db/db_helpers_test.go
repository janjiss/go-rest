package dbHelpers

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateULID(t *testing.T) {
	// generate two ULIDs
	ulid1, err1 := GenerateULID()
	ulid2, err2 := GenerateULID()

	// there should be no errors
	if err1 != nil {
		t.Errorf("There should be no error generating the first ULID, got: %v", err1)
	}
	if err2 != nil {
		t.Errorf("There should be no error generating the second ULID, got: %v", err2)
	}

	// the ULIDs should not be nil
	if ulid1 == uuid.Nil {
		t.Error("The first ULID should not be nil")
	}
	if ulid2 == uuid.Nil {
		t.Error("The second ULID should not be nil")
	}

	// the ULIDs should not be equal (this might fail if the function is broken, or with extremely low probability due to random collision)
	if ulid1 == ulid2 {
		t.Error("The ULIDs should not be equal")
	}

	// the timestamp portion (first 48 bits / 6 bytes) of the ULIDs should be increasing
	ts1 := ulid1[:6]
	ts2 := ulid2[:6]
	if string(ts2) < string(ts1) {
		t.Error("The timestamp portion of the second ULID should be greater than or equal to the first one")
	}
}
