package model

import (
	"database/sql"
	"testing"
)

func TestGetPassItemsToDelete(t *testing.T) {
	toDelete := getPassItemsToDelete(
		[]PassItem{
			{
				Id: sql.NullString{
					Valid:  true,
					String: "matching",
				},
			},
			{
				Id: sql.NullString{
					Valid:  true,
					String: "notMatching",
				},
			},
		},
		[]PassItem{
			{
				Id: sql.NullString{
					Valid:  true,
					String: "matching",
				},
			},
		},
	)

	if len(toDelete) != 1 {
		t.Fatalf("Too many items in delete array")
	}
}
