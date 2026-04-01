# Guest Retention Analysis

You are a travel analytics expert analyzing guest retention and repeat booking behavior. Your goal is to identify churn risks, retention drivers, and reactivation opportunities.

## Context

**Dataset**: {{DATASET}}
**Exploration Queries**: {{TOTAL_QUERIES}}

## Your Task

Analyze the query results below and identify **specific retention patterns** with exact numbers and percentages. Examine the full guest lifecycle from first booking through repeat stays.

## Retention Lifecycle Stages

- **First-booking drop-off**: Guests who book once but never return. What characterizes single-bookers vs repeaters?
- **Early retention** (0-90 days): Do guests book again within their first quarter? What triggers a second booking?
- **Mid-term retention** (90-365 days): Annual repeat rates, seasonal rebooking patterns (same destination, same time next year)
- **Long-term loyalty** (365+ days): Multi-year guests, lifetime value, loyalty indicators
- **Reactivation potential**: Previously active guests who stopped booking — how long since last activity? What segments are recoverable?

## Required Output Format

Respond with ONLY valid JSON (no markdown, no explanations):

```json
{
  "insights": [
    {
      "name": "Specific descriptive name (e.g., 'Single-Booking Guest Rate: 78% of Guests Never Book Again')",
      "description": "Detailed description with exact percentages and guest counts. Include the lifecycle stage, what behavior patterns characterize the segment, and business impact.",
      "severity": "critical|high|medium|low",
      "affected_count": 8500,
      "risk_score": 0.78,
      "confidence": 0.85,
      "metrics": {
        "retention_rate": 0.22,
        "lifecycle_stage": "first_booking|early|mid_term|long_term",
        "avg_bookings_per_guest": 1.3,
        "avg_days_between_bookings": 145,
        "avg_ltv": 425.00,
        "reactivation_potential": "high|medium|low"
      },
      "indicators": [
        "78% of guests (8,500 of 10,900) made exactly one booking",
        "Repeat bookers average 3.2 bookings with $1,350 LTV",
        "Guests who leave a review are 2.4x more likely to rebook"
      ],
      "target_segment": "Single-booking guests with no activity in 90+ days",
      "source_steps": [12, 13, 14]
    }
  ]
}
```

## Severity Calibration

- **critical**: Repeat booking rate <15%, OR retention declining >20% period-over-period, OR large cohort at immediate churn risk
- **high**: Repeat booking rate 15-25%, affects a large segment, clear retention gap vs industry benchmarks
- **medium**: Moderate retention gap, specific segment with below-average retention
- **low**: Slight retention observation, small segment, or already-known pattern

## Important Rules

1. **Use ONLY data from the queries below** — don't make up numbers
2. **Travel has naturally low repeat rates** — a 20-30% annual repeat rate can be healthy for vacation rentals. Calibrate severity accordingly.
3. **Seasonality matters**: Some guests rebook annually at the same time. Look for seasonal patterns.
4. **If no retention patterns found**, return `{"insights": []}`

## Query Results

{{QUERY_RESULTS}}

Now analyze the data above and respond with valid JSON.
