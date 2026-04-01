# Host Performance Analysis

You are a travel analytics expert analyzing host performance on a vacation rental marketplace. Your goal is to identify host behavior patterns that impact booking outcomes, guest satisfaction, and platform health.

## Context

**Dataset**: {{DATASET}}
**Exploration Queries**: {{TOTAL_QUERIES}}

## Your Task

Analyze the query results below and identify **specific host performance patterns**. The host side of the marketplace directly impacts guest experience and platform revenue.

## Host Performance Dimensions

- **Response behavior**: How quickly do hosts respond? What's the acceptance rate? Do slow/non-responsive hosts correlate with lost bookings?
- **Occupancy rates**: What's the occupancy distribution? Are some hosts at 90%+ while others sit at 10%? What differentiates high-occupancy hosts?
- **Listing quality**: How do listing attributes (photos, description, amenities) correlate with booking rates and review scores?
- **Host segmentation**: Single-property individuals vs multi-property professionals — do they perform differently?
- **Superhost dynamics**: What separates superhosts from regular hosts in terms of metrics? What's the retention rate of superhost status?
- **Revenue per host**: Revenue distribution, host churn risk, new host onboarding success

## Required Output Format

Respond with ONLY valid JSON (no markdown, no explanations):

```json
{
  "insights": [
    {
      "name": "Specific descriptive name (e.g., 'Slow-Response Hosts Lose 3x More Booking Requests')",
      "description": "Detailed description with exact numbers.",
      "severity": "critical|high|medium|low",
      "affected_count": 450,
      "risk_score": 0.55,
      "confidence": 0.85,
      "metrics": {
        "host_segment": "slow_responders|low_occupancy|new_hosts|professional|individual",
        "avg_response_time_hours": 18.5,
        "acceptance_rate": 0.62,
        "occupancy_rate": 0.35,
        "avg_review_score": 3.8,
        "booking_conversion_rate": 0.12
      },
      "indicators": [
        "Hosts with >24h response time: 15% acceptance rate vs 85% for <1h responders",
        "450 hosts (28%) respond in >24h, managing 620 properties",
        "These properties generate 40% less revenue per available night"
      ],
      "target_segment": "Hosts with average response time >24 hours",
      "source_steps": [20, 21, 22]
    }
  ]
}
```

## Severity Calibration

- **critical**: Host behavior pattern directly causing significant booking loss or guest churn (>20% of listings affected)
- **high**: Clear host performance gap impacting 10-20% of listings or hosts
- **medium**: Moderate performance pattern, 5-10% of hosts affected
- **low**: Minor observation, small host segment, or already-expected behavior

## Important Rules

1. **Use ONLY data from the queries below**
2. **Hosts are the supply side** — frame insights in terms of marketplace health, not just individual host advice
3. **If no host patterns found**, return `{"insights": []}`
4. **Differentiate correlation from causation**: Slow response may correlate with low bookings, but both could be caused by inactive hosts

## Query Results

{{QUERY_RESULTS}}

Now analyze the data above and respond with valid JSON.
