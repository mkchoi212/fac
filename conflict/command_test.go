package conflict_test

import (
	"strings"
	"testing"

	"github.com/mkchoi212/fac/conflict"
)

var TestGitPath = "assets/dummy_repo"
var TestGitSubPath = "../assets/dummy_repo/assets"
var CorrectNumMarkers = 109

func TestMarkerLocations(t *testing.T) {
	markers, ok := conflict.MarkerLocations("../" + TestGitPath)
	if ok != nil {
		t.Errorf("git diff --check failed with error: %s", ok.Error())
	}

	numMarkers := len(markers)
	if numMarkers != CorrectNumMarkers {
		t.Errorf("MarkerLocations was incorrect: got %d, want, %d", numMarkers, CorrectNumMarkers)
	}
}

func TestMarkerLocationsFromSubPath(t *testing.T) {
	markers, ok := conflict.MarkerLocations(TestGitSubPath)
	if ok != nil {
		t.Errorf("git diff --check from sub-directory failed with error: %s", ok.Error())
	}

	numMarkers := len(markers)
	if numMarkers != CorrectNumMarkers {
		t.Errorf("MarkerLocations was incorrect: got %d, want, %d", numMarkers, CorrectNumMarkers)
	}
}

func TestTopLevelPath(t *testing.T) {
	topPath, ok := conflict.TopLevelPath(TestGitSubPath)

	if ok != nil {
		t.Errorf("git rev-parse --show-toplevel failed")
	}

	if !(strings.Contains(topPath, TestGitPath)) {
		t.Errorf("git rev-parse --show-toplevel was incorrect: got %s, want, %s", topPath, "../"+TestGitPath)
	}
}
