package models

import "time"

// ProjectContext represents accumulated knowledge about a project.
// Enables continuous learning across discovery runs.
type ProjectContext struct {
	ProjectID string    `bson:"project_id" json:"project_id"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// Schema knowledge
	KnownSchemas map[string]SchemaKnowledge `bson:"known_schemas" json:"known_schemas"`

	// Query history
	SuccessfulQueries []QueryHistory `bson:"successful_queries" json:"successful_queries"`
	FailedQueries     []QueryHistory `bson:"failed_queries" json:"failed_queries"`

	// Pattern history (domain-agnostic)
	HistoricalPatterns []HistoricalPattern `bson:"historical_patterns" json:"historical_patterns"`

	// Discovery metadata
	TotalDiscoveries    int       `bson:"total_discoveries" json:"total_discoveries"`
	LastDiscoveryDate   time.Time `bson:"last_discovery_date" json:"last_discovery_date"`
	FirstDiscoveryDate  time.Time `bson:"first_discovery_date" json:"first_discovery_date"`
	ConsecutiveFailures int       `bson:"consecutive_failures" json:"consecutive_failures"`

	// Learning notes
	Notes []ContextNote `bson:"notes" json:"notes"`
}

// SchemaKnowledge tracks what we know about a warehouse table.
type SchemaKnowledge struct {
	TableName     string      `bson:"table_name" json:"table_name"`
	FirstSeen     time.Time   `bson:"first_seen" json:"first_seen"`
	LastSeen      time.Time   `bson:"last_seen" json:"last_seen"`
	SchemaVersion int         `bson:"schema_version" json:"schema_version"`
	CurrentSchema TableSchema `bson:"current_schema" json:"current_schema"`

	UsefulColumns []string `bson:"useful_columns" json:"useful_columns"`
	CommonFilters []string `bson:"common_filters" json:"common_filters"`

	EstimatedRowCount int64 `bson:"estimated_row_count" json:"estimated_row_count"`
}

// HistoricalPattern tracks a discovered pattern over time.
type HistoricalPattern struct {
	PatternID    string    `bson:"pattern_id" json:"pattern_id"`
	AnalysisArea string    `bson:"analysis_area" json:"analysis_area"`
	Name         string    `bson:"name" json:"name"`
	Description  string    `bson:"description" json:"description"`
	FirstSeen    time.Time `bson:"first_seen" json:"first_seen"`
	LastSeen     time.Time `bson:"last_seen" json:"last_seen"`
	SeenCount    int       `bson:"seen_count" json:"seen_count"`
	Status       string    `bson:"status" json:"status"` // active, resolved, worsening, improving
}

// ContextNote represents a learning note.
type ContextNote struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Category  string    `bson:"category" json:"category"`
	Note      string    `bson:"note" json:"note"`
	Relevance float64   `bson:"relevance" json:"relevance"`
}

// NewProjectContext creates a new context for a project.
func NewProjectContext(projectID string) *ProjectContext {
	now := time.Now()
	return &ProjectContext{
		ProjectID:          projectID,
		CreatedAt:          now,
		UpdatedAt:          now,
		KnownSchemas:       make(map[string]SchemaKnowledge),
		SuccessfulQueries:  make([]QueryHistory, 0),
		FailedQueries:      make([]QueryHistory, 0),
		HistoricalPatterns: make([]HistoricalPattern, 0),
		Notes:              make([]ContextNote, 0),
		FirstDiscoveryDate: now,
	}
}

// AddSuccessfulQuery records a successful query (keeps last 100).
func (ctx *ProjectContext) AddSuccessfulQuery(query QueryHistory) {
	query.Success = true
	ctx.SuccessfulQueries = append(ctx.SuccessfulQueries, query)
	if len(ctx.SuccessfulQueries) > 100 {
		ctx.SuccessfulQueries = ctx.SuccessfulQueries[len(ctx.SuccessfulQueries)-100:]
	}
	ctx.UpdatedAt = time.Now()
}

// AddFailedQuery records a failed query (keeps last 50).
func (ctx *ProjectContext) AddFailedQuery(query QueryHistory) {
	query.Success = false
	ctx.FailedQueries = append(ctx.FailedQueries, query)
	if len(ctx.FailedQueries) > 50 {
		ctx.FailedQueries = ctx.FailedQueries[len(ctx.FailedQueries)-50:]
	}
	ctx.UpdatedAt = time.Now()
}

// AddNote adds a learning note (keeps last 200).
func (ctx *ProjectContext) AddNote(category, note string, relevance float64) {
	ctx.Notes = append(ctx.Notes, ContextNote{
		Timestamp: time.Now(),
		Category:  category,
		Note:      note,
		Relevance: relevance,
	})
	if len(ctx.Notes) > 200 {
		ctx.Notes = ctx.Notes[len(ctx.Notes)-200:]
	}
	ctx.UpdatedAt = time.Now()
}

// RecordDiscovery updates context after a discovery run.
func (ctx *ProjectContext) RecordDiscovery(success bool) {
	ctx.TotalDiscoveries++
	ctx.LastDiscoveryDate = time.Now()
	ctx.UpdatedAt = time.Now()
	if success {
		ctx.ConsecutiveFailures = 0
	} else {
		ctx.ConsecutiveFailures++
	}
}

// UpdatePatterns merges new insights into historical patterns.
// Tracks seen count and updates status.
func (ctx *ProjectContext) UpdatePatterns(insights []Insight) {
	existing := make(map[string]*HistoricalPattern)
	for i := range ctx.HistoricalPatterns {
		existing[ctx.HistoricalPatterns[i].PatternID] = &ctx.HistoricalPatterns[i]
	}

	now := time.Now()
	for _, ins := range insights {
		// Use analysis_area:name as pattern ID
		pid := ins.AnalysisArea + ":" + ins.Name
		if p, ok := existing[pid]; ok {
			p.SeenCount++
			p.LastSeen = now
			p.Status = "recurring"
		} else {
			ctx.HistoricalPatterns = append(ctx.HistoricalPatterns, HistoricalPattern{
				PatternID:    pid,
				AnalysisArea: ins.AnalysisArea,
				Name:         ins.Name,
				Description:  ins.Description,
				FirstSeen:    now,
				LastSeen:     now,
				SeenCount:    1,
				Status:       "active",
			})
		}
	}

	// Keep last 200 patterns
	if len(ctx.HistoricalPatterns) > 200 {
		ctx.HistoricalPatterns = ctx.HistoricalPatterns[len(ctx.HistoricalPatterns)-200:]
	}
	ctx.UpdatedAt = now
}

// InsightSummary is a compact representation of a previous insight for LLM context.
type InsightSummary struct {
	Name          string `json:"name"`
	AnalysisArea  string `json:"analysis_area"`
	Severity      string `json:"severity"`
	AffectedCount int    `json:"affected_count"`
	Date          string `json:"date"`
}

// FeedbackSummary is a compact representation of user feedback for LLM context.
type FeedbackSummary struct {
	InsightName string `json:"insight_name"`
	Rating      string `json:"rating"`
	Comment     string `json:"comment,omitempty"`
}

// RecommendationSummary is a compact representation of a previous recommendation.
type RecommendationSummary struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Priority int    `json:"priority"`
}
