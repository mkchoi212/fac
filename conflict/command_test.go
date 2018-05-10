package conflict

import (
	"strings"
	"testing"
)

var TestGitPath = "../test"
var TestGitSubPath = "../test/assets"
var CorrectNumMarkers = 75

func TestMarkerLocations(t *testing.T) {
	markers, ok := MarkerLocations(TestGitPath)
	if ok != nil {
		t.Errorf("git diff --check failed with error: %s", ok.Error())
	}

	numMarkers := len(markers)
	if numMarkers != CorrectNumMarkers {
		t.Errorf("MarkerLocations was incorrect: got %d, want, %d", numMarkers, CorrectNumMarkers)
	}
}

func TestMarkerLocationsFromSubPath(t *testing.T) {
	markers, ok := MarkerLocations(TestGitSubPath)
	if ok != nil {
		t.Errorf("git diff --check from sub-directory failed with error: %s", ok.Error())
	}

	numMarkers := len(markers)
	if numMarkers != CorrectNumMarkers {
		t.Errorf("MarkerLocations was incorrect: got %d, want, %d", numMarkers, CorrectNumMarkers)
	}
}

func TestTopLevelPath(t *testing.T) {
	topPath, ok := TopLevelPath(TestGitSubPath)

	if ok != nil {
		t.Errorf("git rev-parse --show-toplevel failed")
	}

	if !(strings.Contains(topPath, "fac/test")) {
		t.Errorf("git rev-parse --show-toplevel was incorrect: got %s, want, %s", topPath, "../fac/test")
	}
}
