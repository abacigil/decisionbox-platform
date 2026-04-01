# Generate Actionable Recommendations

You are a travel and hospitality analytics expert creating **specific, actionable recommendations** based on discovered patterns. Every recommendation must be concrete enough that a platform operator could implement it immediately.

## Context

**Discovery Date**: {{DISCOVERY_DATE}}
**Insights Found**: {{INSIGHTS_SUMMARY}}

## Your Task

Generate **specific, actionable recommendations** that can be immediately implemented. Each recommendation must include:

1. **Clear action** — What exactly to do, with specific parameters
2. **Target segment** — Who to target, with exact criteria that can be used for segmentation
3. **Expected impact** — Quantified expected results based on the data
4. **Implementation steps** — Concrete steps to implement this recommendation

## Output Format

Respond with ONLY valid JSON:

```json
{
  "recommendations": [
    {
      "title": "Action — Context (e.g., 'Send Rebooking Incentive to Guests 60 Days After Checkout')",
      "description": "Detailed explanation with numbers. What is the problem? How big is the impact? Why does this recommendation address it? What evidence supports it?",
      "category": "booking_conversion|revenue|guest_retention|host_performance|guest_satisfaction|growth",
      "priority": 1,
      "effort": "quick_win|moderate|major_initiative",
      "target_segment": "Exact segment definition with measurable criteria (e.g., 'Guests who completed one stay, left a 4+ star review, and have not booked again within 60 days')",
      "segment_size": 1234,
      "expected_impact": {
        "metric": "conversion_rate|revenue|retention_rate|avg_booking_value|review_score|occupancy_rate",
        "estimated_improvement": "15-20%",
        "reasoning": "Why we expect this improvement, with supporting data points from the analysis"
      },
      "actions": [
        "Specific implementation step 1 with parameters",
        "Specific implementation step 2 with parameters",
        "Specific implementation step 3 with parameters"
      ],
      "success_metrics": [
        "Track repeat booking rate for targeted segment (target: increase from 18% to 25%)",
        "Monitor revenue per guest for this cohort (target: improve by 15%)"
      ],
      "related_insight_ids": ["insight-id-1", "insight-id-2"],
      "confidence": 0.85
    }
  ]
}
```

**IMPORTANT:** Each recommendation MUST include `related_insight_ids` — an array of insight `id` values from the input data that this recommendation addresses. Copy the exact `id` values from the insights provided below.

## Requirements

### DO create recommendations that are:
- **Specific**: Exact numbers, time windows, thresholds, segments, property types
- **Actionable**: A product or ops team knows exactly what to implement after reading this
- **Measurable**: Clear success metrics with baseline and target values
- **Data-backed**: Every recommendation is grounded in the insights, not generic advice

### DON'T create recommendations that are:
- Generic ("improve the booking experience", "optimize pricing")
- Vague ("send targeted emails", "A/B test the checkout flow")
- Missing numbers or specifics
- Not supported by the discovered insights
- Duplicating another recommendation with different wording

### Effort Scale:
- **quick_win**: Can be implemented in hours to a day. Email campaign, pricing adjustment, feature flag toggle.
- **moderate**: Requires development work, typically 1-2 weeks. New feature, UI change, algorithm update.
- **major_initiative**: Significant engineering effort, typically weeks to months. New system, marketplace redesign, new product offering.

### Priority Scale (P1 = highest):
- **1 (Critical)**: Large guest/host base affected AND high revenue/retention impact. Do this first.
- **2 (High)**: Significant impact, implement soon. Strong evidence and clear path to implementation.
- **3 (Medium)**: Moderate impact. Worth doing but not urgent.
- **4 (Low)**: Nice to have. Small improvement or low-confidence opportunity.
- **5 (Optional)**: Minor improvement. Consider if resources are available.

## Recommendation Quality Checklist

Before including a recommendation, verify:
- Does it reference specific insights from the data? (related_insight_ids)
- Is the target segment precisely defined with measurable criteria?
- Can an ops/product team implement this without asking clarifying questions?
- Are the success_metrics specific enough to measure impact?
- Is the expected_impact realistic based on the data?

## Discovered Insights

{{INSIGHTS_DATA}}

---

Generate 3-8 specific, actionable recommendations. Prioritize by impact and urgency. Focus on recommendations where the data clearly supports the expected outcome.
