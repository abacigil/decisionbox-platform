package gaming

import (
	"os"
	"path/filepath"
	"testing"
)

// --- areas.json loading ---

func TestLoadAreas_Base(t *testing.T) {
	areas := loadAreas(filepath.Join(getPromptsPath(), "base", "areas.json"))

	if len(areas) != 3 {
		t.Fatalf("base areas = %d, want 3", len(areas))
	}

	ids := map[string]bool{}
	for _, a := range areas {
		ids[a.ID] = true
		if a.Name == "" {
			t.Errorf("area %q has empty name", a.ID)
		}
		if a.PromptFile == "" {
			t.Errorf("area %q has empty prompt_file", a.ID)
		}
		if len(a.Keywords) == 0 {
			t.Errorf("area %q has no keywords", a.ID)
		}
		if a.Priority == 0 {
			t.Errorf("area %q has zero priority", a.ID)
		}
	}

	for _, expected := range []string{"churn", "engagement", "monetization"} {
		if !ids[expected] {
			t.Errorf("missing base area: %s", expected)
		}
	}
}

func TestLoadAreas_Match3(t *testing.T) {
	areas := loadAreas(filepath.Join(getPromptsPath(), "categories", "match3", "areas.json"))

	if len(areas) != 2 {
		t.Fatalf("match3 areas = %d, want 2", len(areas))
	}

	ids := map[string]bool{}
	for _, a := range areas {
		ids[a.ID] = true
	}

	if !ids["levels"] || !ids["boosters"] {
		t.Errorf("missing match3 areas: got %v", ids)
	}
}

func TestLoadAreas_NonexistentFile(t *testing.T) {
	areas := loadAreas("/nonexistent/path/areas.json")

	if areas != nil {
		t.Errorf("should return nil for missing file, got %v", areas)
	}
}

func TestLoadAreas_MalformedJSON(t *testing.T) {
	// Create a temp file with bad JSON
	tmp, err := os.CreateTemp("", "areas-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	tmp.WriteString(`{not valid json`)
	tmp.Close()

	areas := loadAreas(tmp.Name())
	if areas != nil {
		t.Errorf("should return nil for malformed JSON, got %v", areas)
	}
}

func TestLoadAreas_EmptyArray(t *testing.T) {
	tmp, err := os.CreateTemp("", "areas-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	tmp.WriteString(`[]`)
	tmp.Close()

	areas := loadAreas(tmp.Name())
	if len(areas) != 0 {
		t.Errorf("should return empty for empty array")
	}
}

func TestLoadAreas_PromptFileLinksCorrectly(t *testing.T) {
	areas := loadAreas(filepath.Join(getPromptsPath(), "base", "areas.json"))

	for _, area := range areas {
		promptPath := filepath.Join(getPromptsPath(), "base", area.PromptFile)
		content := readPromptFile(promptPath)
		if content == "" {
			t.Errorf("area %q prompt_file %q resolves to empty content", area.ID, area.PromptFile)
		}
	}
}

func TestLoadAreas_CategoryPromptFileLinksCorrectly(t *testing.T) {
	areas := loadAreas(filepath.Join(getPromptsPath(), "categories", "match3", "areas.json"))

	for _, area := range areas {
		promptPath := filepath.Join(getPromptsPath(), "categories", "match3", area.PromptFile)
		content := readPromptFile(promptPath)
		if content == "" {
			t.Errorf("area %q prompt_file %q resolves to empty content", area.ID, area.PromptFile)
		}
	}
}

func TestReadPromptFile_Exists(t *testing.T) {
	content := readPromptFile(filepath.Join(getPromptsPath(), "base", "exploration.md"))
	if content == "" {
		t.Error("exploration.md should not be empty")
	}
}

func TestReadPromptFile_NotExists(t *testing.T) {
	content := readPromptFile("/nonexistent/file.md")
	if content != "" {
		t.Error("should return empty for missing file")
	}
}

func TestGetPromptsPath_Default(t *testing.T) {
	old := os.Getenv("DOMAIN_PACK_PATH")
	os.Unsetenv("DOMAIN_PACK_PATH")
	defer os.Setenv("DOMAIN_PACK_PATH", old)

	path := getPromptsPath()
	if path != "domain-packs/gaming/prompts" {
		t.Errorf("default path = %q", path)
	}
}

func TestGetPromptsPath_EnvOverride(t *testing.T) {
	old := os.Getenv("DOMAIN_PACK_PATH")
	os.Setenv("DOMAIN_PACK_PATH", "/custom/path")
	defer os.Setenv("DOMAIN_PACK_PATH", old)

	path := getPromptsPath()
	if path != "/custom/path/prompts" {
		t.Errorf("path = %q, want /custom/path/prompts", path)
	}
}

func TestGetProfilesPath_EnvOverride(t *testing.T) {
	old := os.Getenv("DOMAIN_PACK_PATH")
	os.Setenv("DOMAIN_PACK_PATH", "/custom/path")
	defer os.Setenv("DOMAIN_PACK_PATH", old)

	path := getProfilesPath()
	if path != "/custom/path/profiles" {
		t.Errorf("path = %q, want /custom/path/profiles", path)
	}
}

// --- AnalysisAreas with unknown category ---

func TestAnalysisAreas_UnknownCategory_ReturnsBaseOnly(t *testing.T) {
	pack := NewPack()
	areas := pack.AnalysisAreas("nonexistent_category")

	// Should return only base areas (3)
	if len(areas) != 3 {
		t.Errorf("areas = %d, want 3 (base only for unknown category)", len(areas))
	}
}
