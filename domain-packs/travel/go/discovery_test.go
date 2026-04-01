package travel

import (
	"os"
	"strings"
	"testing"

	"github.com/decisionbox-io/decisionbox/libs/go-common/domainpack"
)

func init() {
	// Tests run from domain-packs/travel/go/ — DOMAIN_PACK_PATH is the domain-packs root
	os.Setenv("DOMAIN_PACK_PATH", "../..")
}

func TestTravelPackImplementsDiscoveryPack(t *testing.T) {
	pack := NewPack()
	dp, ok := domainpack.AsDiscoveryPack(pack)
	if !ok {
		t.Fatal("TravelPack does not implement DiscoveryPack")
	}
	if dp == nil {
		t.Fatal("AsDiscoveryPack returned nil")
	}
}

func TestName(t *testing.T) {
	pack := NewPack()
	if pack.Name() != "travel" {
		t.Errorf("Name() = %q, want %q", pack.Name(), "travel")
	}
}

func TestDomainCategories(t *testing.T) {
	pack := NewPack()
	cats := pack.DomainCategories()

	if len(cats) != 1 {
		t.Fatalf("DomainCategories returned %d categories, want 1", len(cats))
	}

	c := cats[0]
	if c.ID != "vacation_rental" {
		t.Errorf("category ID = %q, want %q", c.ID, "vacation_rental")
	}
	if c.Name == "" {
		t.Error("category has empty Name")
	}
	if c.Description == "" {
		t.Error("category has empty Description")
	}
}

func TestAnalysisAreasBase(t *testing.T) {
	pack := NewPack()
	areas := pack.AnalysisAreas("")

	if len(areas) != 5 {
		t.Errorf("base areas = %d, want 5", len(areas))
	}

	ids := make(map[string]bool)
	for _, a := range areas {
		ids[a.ID] = true
		if !a.IsBase {
			t.Errorf("area %q should be IsBase=true", a.ID)
		}
	}

	for _, expected := range []string{"booking_conversion", "revenue", "guest_retention", "cancellation", "destination"} {
		if !ids[expected] {
			t.Errorf("missing base area: %s", expected)
		}
	}
}

func TestAnalysisAreasVacationRental(t *testing.T) {
	pack := NewPack()
	areas := pack.AnalysisAreas("vacation_rental")

	if len(areas) != 7 {
		t.Errorf("vacation_rental areas = %d, want 7 (5 base + 2 category)", len(areas))
	}

	ids := make(map[string]bool)
	for _, a := range areas {
		ids[a.ID] = true
	}

	for _, expected := range []string{"booking_conversion", "revenue", "guest_retention", "cancellation", "destination", "host_performance", "guest_satisfaction"} {
		if !ids[expected] {
			t.Errorf("missing area: %s", expected)
		}
	}
}

func TestAnalysisAreasUnknownCategory(t *testing.T) {
	pack := NewPack()
	areas := pack.AnalysisAreas("nonexistent")

	// Should return only base areas
	if len(areas) != 5 {
		t.Errorf("unknown category areas = %d, want 5 (base only)", len(areas))
	}
}

func TestAnalysisAreaKeywords(t *testing.T) {
	pack := NewPack()
	areas := pack.AnalysisAreas("vacation_rental")

	for _, a := range areas {
		if len(a.Keywords) == 0 {
			t.Errorf("area %q has no keywords", a.ID)
		}
		if a.Name == "" {
			t.Errorf("area %q has empty Name", a.ID)
		}
		if a.Priority == 0 {
			t.Errorf("area %q has zero Priority", a.ID)
		}
	}
}

func TestPromptsBase(t *testing.T) {
	pack := NewPack()
	prompts := pack.Prompts("")

	if prompts.Exploration == "" {
		t.Error("Exploration prompt is empty")
	}
	if !strings.Contains(prompts.Exploration, "Travel & Hospitality Analytics Discovery") {
		t.Error("Exploration prompt missing expected header")
	}
	if prompts.Recommendations == "" {
		t.Error("Recommendations prompt is empty")
	}

	// Should have base analysis areas only
	for _, id := range []string{"booking_conversion", "revenue", "guest_retention", "cancellation", "destination"} {
		if _, ok := prompts.AnalysisAreas[id]; !ok {
			t.Errorf("missing base analysis prompt: %s", id)
		}
	}

	// Should NOT have category-specific areas
	if _, ok := prompts.AnalysisAreas["host_performance"]; ok {
		t.Error("base prompts should not include 'host_performance' area")
	}

	// BaseContext should be loaded
	if prompts.BaseContext == "" {
		t.Error("BaseContext is empty — base_context.md not loaded")
	}
	if !strings.Contains(prompts.BaseContext, "{{PROFILE}}") {
		t.Error("BaseContext missing {{PROFILE}} placeholder")
	}
	if !strings.Contains(prompts.BaseContext, "{{PREVIOUS_CONTEXT}}") {
		t.Error("BaseContext missing {{PREVIOUS_CONTEXT}} placeholder")
	}
}

func TestPromptsBase_NoProfileInExploration(t *testing.T) {
	pack := NewPack()
	prompts := pack.Prompts("")

	if strings.Contains(prompts.Exploration, "{{PROFILE}}") {
		t.Error("exploration prompt should not contain {{PROFILE}} — moved to base_context")
	}
	if strings.Contains(prompts.Exploration, "{{PREVIOUS_CONTEXT}}") {
		t.Error("exploration prompt should not contain {{PREVIOUS_CONTEXT}} — moved to base_context")
	}
}

func TestPromptsBase_NoProfileInAnalysis(t *testing.T) {
	pack := NewPack()
	prompts := pack.Prompts("vacation_rental")

	for id, content := range prompts.AnalysisAreas {
		if strings.Contains(content, "{{PROFILE}}") {
			t.Errorf("analysis prompt %q should not contain {{PROFILE}} — moved to base_context", id)
		}
		if strings.Contains(content, "{{PREVIOUS_CONTEXT}}") {
			t.Errorf("analysis prompt %q should not contain {{PREVIOUS_CONTEXT}} — moved to base_context", id)
		}
	}
}

func TestPromptsVacationRentalMerge(t *testing.T) {
	pack := NewPack()
	prompts := pack.Prompts("vacation_rental")

	// Exploration should contain both base and vacation_rental context
	if !strings.Contains(prompts.Exploration, "Travel & Hospitality Analytics Discovery") {
		t.Error("merged exploration missing base content")
	}
	if !strings.Contains(prompts.Exploration, "Vacation Rental Marketplace Context") {
		t.Error("merged exploration missing vacation_rental context")
	}

	// Should have base + category-specific analysis areas
	for _, id := range []string{"booking_conversion", "revenue", "guest_retention", "cancellation", "destination", "host_performance", "guest_satisfaction"} {
		content, ok := prompts.AnalysisAreas[id]
		if !ok {
			t.Errorf("missing analysis prompt: %s", id)
			continue
		}
		if content == "" {
			t.Errorf("empty analysis prompt: %s", id)
		}
	}
}

func TestProfileSchemaBase(t *testing.T) {
	pack := NewPack()
	schema := pack.ProfileSchema("")

	if _, ok := schema["error"]; ok {
		t.Fatalf("ProfileSchema returned error: %v", schema["error"])
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema has no properties")
	}

	for _, expected := range []string{"platform_info", "booking_model", "monetization", "growth_features", "kpis"} {
		if _, ok := props[expected]; !ok {
			t.Errorf("base schema missing property: %s", expected)
		}
	}
}

func TestProfileSchemaVacationRentalMerge(t *testing.T) {
	pack := NewPack()
	schema := pack.ProfileSchema("vacation_rental")

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema has no properties")
	}

	// Should have base properties
	for _, expected := range []string{"platform_info", "booking_model", "monetization", "kpis"} {
		if _, ok := props[expected]; !ok {
			t.Errorf("merged schema missing base property: %s", expected)
		}
	}

	// Should have vacation_rental-specific properties
	for _, expected := range []string{"property_types", "amenities", "host_program", "seasonal_patterns"} {
		if _, ok := props[expected]; !ok {
			t.Errorf("merged schema missing vacation_rental property: %s", expected)
		}
	}
}
