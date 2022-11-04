package main

import "testing"

func TestWithExistingStation(t *testing.T) {
	_, err := getStationByShortCode("VIA")
	if err != nil {
		t.Fatalf("Station %s should have been found", "VIA")
	}
}

func TestWithNonExistingStation(t *testing.T) {
	_, err := getStationByShortCode("FOOBAR")
	if err == nil {
		t.Fatalf("Station %s should have not been found", "FOOBAR")
	}
}
