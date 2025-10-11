package version

import (
	"testing"
)

func TestVersionVariables(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if GitHash == "" {
		t.Error("GitHash should not be empty")
	}

	if GoBuildEnv == "" {
		t.Error("GoBuildEnv should not be empty")
	}

	if GoBuildTime == "" {
		t.Error("GoBuildTime should not be empty")
	}
}

func TestVersionDefault(t *testing.T) {
	// Test default values
	if Version != "dev" && Version == "" {
		t.Errorf("Version should be 'dev' by default or set, got '%s'", Version)
	}

	if GitHash != "dev" && GitHash == "" {
		t.Errorf("GitHash should be 'dev' by default or set, got '%s'", GitHash)
	}

	if GoBuildEnv != "dev" && GoBuildEnv == "" {
		t.Errorf("GoBuildEnv should be 'dev' by default or set, got '%s'", GoBuildEnv)
	}

	if GoBuildTime != "dev" && GoBuildTime == "" {
		t.Errorf("GoBuildTime should be 'dev' by default or set, got '%s'", GoBuildTime)
	}
}
