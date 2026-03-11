package discovery

import (
	"testing"

	"github.com/decisionbox-io/decisionbox/libs/go-common/domainpack"
	"github.com/decisionbox-io/decisionbox/services/agent/internal/models"
)

// --- Selective Discovery (area filtering) ---

func TestResolvePrompts_NoProjectPrompts(t *testing.T) {
	o := &Orchestrator{projectPrompts: nil}

	dpPrompts := domainpack.PromptTemplates{
		Exploration:     "explore",
		Recommendations: "recommend",
		AnalysisAreas:   map[string]string{"churn": "analyze churn"},
	}
	dpAreas := []domainpack.AnalysisArea{
		{ID: "churn", Name: "Churn", IsBase: true},
	}

	prompts, areas := o.resolvePrompts(dpPrompts, dpAreas)

	if prompts.Exploration != "explore" {
		t.Error("should return domain pack prompts when no project prompts")
	}
	if len(areas) != 1 || areas[0].ID != "churn" {
		t.Error("should return domain pack areas")
	}
}

func TestResolvePrompts_ProjectOverridesExploration(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			Exploration:   "custom exploration",
			AnalysisAreas: map[string]models.AnalysisAreaConfig{},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		Exploration:     "default exploration",
		Recommendations: "default recs",
		AnalysisAreas:   map[string]string{"churn": "default churn"},
	}

	prompts, _ := o.resolvePrompts(dpPrompts, nil)

	if prompts.Exploration != "custom exploration" {
		t.Errorf("exploration = %q, want custom", prompts.Exploration)
	}
	if prompts.Recommendations != "default recs" {
		t.Error("recommendations should keep default when project doesn't override")
	}
}

func TestResolvePrompts_ProjectOverridesRecommendations(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			Recommendations: "custom recs",
			AnalysisAreas:   map[string]models.AnalysisAreaConfig{},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		Exploration:     "default exploration",
		Recommendations: "default recs",
		AnalysisAreas:   map[string]string{},
	}

	prompts, _ := o.resolvePrompts(dpPrompts, nil)

	if prompts.Recommendations != "custom recs" {
		t.Errorf("recommendations = %q, want custom", prompts.Recommendations)
	}
}

func TestResolvePrompts_ProjectOverridesAreaPrompt(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			AnalysisAreas: map[string]models.AnalysisAreaConfig{
				"churn": {
					Name:     "Custom Churn",
					Keywords: []string{"churn"},
					Prompt:   "my custom churn prompt",
					IsBase:   true,
					Enabled:  true,
					Priority: 1,
				},
			},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		AnalysisAreas: map[string]string{"churn": "default churn prompt"},
	}
	dpAreas := []domainpack.AnalysisArea{
		{ID: "churn", Name: "Churn Risks", IsBase: true, Priority: 1},
	}

	prompts, areas := o.resolvePrompts(dpPrompts, dpAreas)

	if prompts.AnalysisAreas["churn"] != "my custom churn prompt" {
		t.Errorf("churn prompt = %q, want custom", prompts.AnalysisAreas["churn"])
	}
	if len(areas) != 1 || areas[0].Name != "Custom Churn" {
		t.Error("should use project area name")
	}
}

func TestResolvePrompts_DisabledArea(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			AnalysisAreas: map[string]models.AnalysisAreaConfig{
				"churn": {Name: "Churn", Enabled: true, Prompt: "x", Priority: 1},
				"engagement": {Name: "Engagement", Enabled: false, Prompt: "x", Priority: 2},
			},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		AnalysisAreas: map[string]string{"churn": "c", "engagement": "e"},
	}
	dpAreas := []domainpack.AnalysisArea{
		{ID: "churn"}, {ID: "engagement"},
	}

	prompts, areas := o.resolvePrompts(dpPrompts, dpAreas)

	if _, ok := prompts.AnalysisAreas["engagement"]; ok {
		t.Error("disabled area should be removed from prompts")
	}
	if len(areas) != 1 || areas[0].ID != "churn" {
		t.Errorf("areas = %v, want only churn", areas)
	}
}

func TestResolvePrompts_CustomArea(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			AnalysisAreas: map[string]models.AnalysisAreaConfig{
				"churn": {Name: "Churn", Enabled: true, Prompt: "c", IsBase: true, Priority: 1},
				"whales": {Name: "Whale Analysis", Enabled: true, Prompt: "find whales",
					IsCustom: true, Keywords: []string{"whale", "spend"}, Priority: 10},
			},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		AnalysisAreas: map[string]string{"churn": "c"},
	}
	dpAreas := []domainpack.AnalysisArea{
		{ID: "churn", IsBase: true},
	}

	prompts, areas := o.resolvePrompts(dpPrompts, dpAreas)

	if prompts.AnalysisAreas["whales"] != "find whales" {
		t.Error("custom area prompt should be added")
	}
	found := false
	for _, a := range areas {
		if a.ID == "whales" {
			found = true
		}
	}
	if !found {
		t.Error("custom area should appear in areas list")
	}
}

func TestResolvePrompts_DomainPackAreasPreservedIfNotInProject(t *testing.T) {
	o := &Orchestrator{
		projectPrompts: &models.ProjectPrompts{
			AnalysisAreas: map[string]models.AnalysisAreaConfig{
				"churn": {Name: "Churn", Enabled: true, Prompt: "c", Priority: 1},
				// engagement is NOT in project — should come from domain pack
			},
		},
	}

	dpPrompts := domainpack.PromptTemplates{
		AnalysisAreas: map[string]string{"churn": "c", "engagement": "e"},
	}
	dpAreas := []domainpack.AnalysisArea{
		{ID: "churn", IsBase: true},
		{ID: "engagement", IsBase: true},
	}

	_, areas := o.resolvePrompts(dpPrompts, dpAreas)

	ids := map[string]bool{}
	for _, a := range areas {
		ids[a.ID] = true
	}
	if !ids["engagement"] {
		t.Error("domain pack areas not in project should be preserved")
	}
}

// --- DiscoveryResult run type ---

func TestDiscoveryResult_RunType(t *testing.T) {
	result := models.DiscoveryResult{
		RunType:        "partial",
		AreasRequested: []string{"churn", "levels"},
	}

	if result.RunType != "partial" {
		t.Error("RunType should be partial")
	}
	if len(result.AreasRequested) != 2 {
		t.Error("AreasRequested should have 2 areas")
	}
}

func TestDiscoveryResult_FullRun(t *testing.T) {
	result := models.DiscoveryResult{RunType: "full"}

	if result.RunType != "full" {
		t.Error("RunType should be full")
	}
	if result.AreasRequested != nil {
		t.Error("AreasRequested should be nil for full run")
	}
}
