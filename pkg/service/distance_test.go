package service

import (
	"testing"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/external"
)

func TestCalcDistance(t *testing.T) {
	distInKm := int(calcDistance(-34.0, -64.0, 64.0, 26.0, "K"))
	if distInKm != 13361 {
		t.Errorf("from Argentina to Finland [latlng (64.0, 26.0)] there are 13361 km and we calculed: %d km", distInKm)
	}
}

func TestResolveDistance(t *testing.T) {
	var database db.DBInterface = &db.MockDB{
		ExistsMock: func(hash, key string) bool {
			if key == "FIN" {
				return true
			}
			return false
		},
		RetrieveMock: func(hash, key string) string { return "13361" },
	}

	var countryAR external.IP2countryResponse = external.IP2countryResponse{"AR", "ARG", "Argentina"}
	var countryFI external.IP2countryResponse = external.IP2countryResponse{"FI", "FIN", "Finland"}
	var countryEH external.IP2countryResponse = external.IP2countryResponse{"EH", "ESH", "Western Sahara"}

	dist := ResolveDistance(database, countryAR)
	if dist != 0 {
		t.Errorf("from Argentina to Argentina there are 0 km and we calculed: %d km", dist)
	}

	dist = ResolveDistance(database, countryFI)
	if dist != 13361 {
		t.Errorf("from Argentina to Finland there are 13361 km and we calculed: %d km", dist)
	}

	dist = ResolveDistance(database, countryEH)
	if dist != 8444 {
		t.Errorf("from Argentina to Argentina there are 0 km and we calculed: %d km", dist)
	}
}
