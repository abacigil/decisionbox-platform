package handler

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/decisionbox-io/decisionbox/services/api/internal/database"
)

// DiscoveriesHandler handles discovery result endpoints.
type DiscoveriesHandler struct {
	repo        *database.DiscoveryRepository
	projectRepo *database.ProjectRepository
	runRepo     *database.RunRepository
}

func NewDiscoveriesHandler(repo *database.DiscoveryRepository, projectRepo *database.ProjectRepository, runRepo *database.RunRepository) *DiscoveriesHandler {
	return &DiscoveriesHandler{repo: repo, projectRepo: projectRepo, runRepo: runRepo}
}

// List returns discovery results for a project.
// GET /api/v1/projects/{id}/discoveries
func (h *DiscoveriesHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	p, err := h.projectRepo.GetByID(r.Context(), projectID)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	results, err := h.repo.List(r.Context(), projectID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list discoveries: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, results)
}

// GetLatest returns the most recent discovery for a project.
// GET /api/v1/projects/{id}/discoveries/latest
func (h *DiscoveriesHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	result, err := h.repo.GetLatest(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get discovery: "+err.Error())
		return
	}
	if result == nil {
		writeError(w, http.StatusNotFound, "no discoveries found")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetByDate returns a discovery for a specific date.
// GET /api/v1/projects/{id}/discoveries/{date}
func (h *DiscoveriesHandler) GetByDate(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	dateStr := r.PathValue("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	result, err := h.repo.GetByDate(r.Context(), projectID, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get discovery: "+err.Error())
		return
	}
	if result == nil {
		writeError(w, http.StatusNotFound, "no discovery found for date "+dateStr)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// TriggerDiscovery triggers a discovery run for a project.
// POST /api/v1/projects/{id}/discover
func (h *DiscoveriesHandler) TriggerDiscovery(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	p, err := h.projectRepo.GetByID(r.Context(), projectID)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	// Check if there's already a running discovery
	running, _ := h.runRepo.GetRunningByProject(r.Context(), projectID)
	if running != nil {
		writeJSON(w, http.StatusConflict, map[string]string{
			"status":  "already_running",
			"run_id":  running.ID,
			"message": "A discovery is already running for this project",
		})
		return
	}

	// Create a run record
	runID, err := h.runRepo.Create(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create run: "+err.Error())
		return
	}

	// Spawn the agent as a background subprocess
	cmd := exec.Command("decisionbox-agent",
		"--project-id", projectID,
		"--run-id", runID,
	)

	if err := cmd.Start(); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to start agent: %s", err.Error()))
		return
	}

	// Don't wait — let it run in the background
	go cmd.Wait()

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "started",
		"run_id": runID,
		"message": "Discovery agent started",
	})
}

// GetStatus returns the live discovery status for a project.
// GET /api/v1/projects/{id}/status
func (h *DiscoveriesHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	p, err := h.projectRepo.GetByID(r.Context(), projectID)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	// Get the latest run (for live status)
	latestRun, _ := h.runRepo.GetLatestByProject(r.Context(), projectID)

	status := map[string]interface{}{
		"project_id": projectID,
	}

	if latestRun != nil {
		status["run"] = latestRun
	}

	// Also include latest completed discovery stats
	latest, _ := h.repo.GetLatest(r.Context(), projectID)
	if latest != nil {
		status["last_discovery"] = map[string]interface{}{
			"date":            latest.DiscoveryDate,
			"insights_count":  len(latest.Insights),
			"total_steps":     latest.TotalSteps,
		}
	}

	writeJSON(w, http.StatusOK, status)
}

// GetRun returns a specific discovery run by ID.
// GET /api/v1/runs/{runId}
func (h *DiscoveriesHandler) GetRun(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("runId")

	run, err := h.runRepo.GetByID(r.Context(), runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get run: "+err.Error())
		return
	}
	if run == nil {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}

	writeJSON(w, http.StatusOK, run)
}
