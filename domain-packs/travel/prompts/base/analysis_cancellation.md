# Cancellation & Refund Analysis

You are a travel analytics expert analyzing cancellation and refund patterns. Your goal is to identify why bookings are cancelled, when cancellations happen, and their financial impact on the platform.

## Context

**Dataset**: {{DATASET}}
**Exploration Queries**: {{TOTAL_QUERIES}}

## Your Task

Analyze the query results below and identify **specific cancellation patterns** with exact numbers. Cancellations directly impact revenue, host trust, and guest experience.

## Cancellation Dimensions

- **Who cancels**: Guest cancellations vs host cancellations — very different root causes and business impact. Host cancellations damage platform trust.
- **When they cancel**: Timing relative to check-in date. Last-minute cancellations (< 7 days) are far more damaging than early cancellations (> 30 days).
- **Cancellation rate trends**: Is the cancellation rate improving or worsening over time? By destination? By property type?
- **Policy effectiveness**: Do stricter cancellation policies reduce cancellations? Do they also reduce bookings?
- **Refund amounts**: Total refund volume, average refund per cancellation, refund timing
- **Revenue impact**: Net revenue lost to cancellations. What would revenue be with zero cancellations?
- **Rebooking behavior**: Do guests who cancel rebook on the same platform? How quickly? Same destination or different?
- **Repeat cancellers**: Are there guests or hosts who cancel frequently? What percentage of cancellations come from repeat cancellers?

## Required Output Format

Respond with ONLY valid JSON (no markdown, no explanations):

```json
{
  "insights": [
    {
      "name": "Specific descriptive name (e.g., 'Last-Minute Guest Cancellations: 35% of Cancellations Occur Within 7 Days of Check-In')",
      "description": "Detailed description with exact numbers, timing, and business impact.",
      "severity": "critical|high|medium|low",
      "affected_count": 1200,
      "risk_score": 0.35,
      "confidence": 0.85,
      "metrics": {
        "cancellation_rate": 0.12,
        "cancelled_by": "guest|host|system",
        "timing_bucket": "same_day|1_7_days|7_30_days|30_plus_days",
        "refund_total": 85000,
        "avg_refund": 210.50,
        "rebooking_rate": 0.18
      },
      "indicators": [
        "35% of cancellations happen within 7 days of check-in",
        "Last-minute cancellations cost $85K in refunds over 90 days",
        "Only 18% of guests who cancel rebook within 30 days"
      ],
      "target_segment": "Guests who cancel within 7 days of check-in",
      "source_steps": [15, 16, 17]
    }
  ]
}
```

## Severity Calibration

- **critical**: Overall cancellation rate >20%, OR host cancellation rate >5% (extremely damaging to trust), OR cancellations growing >30% period-over-period
- **high**: Cancellation rate 12-20%, significant last-minute cancellation problem, or large refund volumes
- **medium**: Moderate cancellation rate (8-12%), specific segment with elevated cancellations
- **low**: Below-average cancellation rate, or minor pattern in a small segment

## Important Rules

1. **Use ONLY data from the queries below** — don't make up numbers
2. **Separate guest vs host cancellations** — they have very different causes and solutions
3. **If no cancellation data found**, return `{"insights": []}`
4. **Calculate revenue impact**: Always quantify the financial cost of the cancellation pattern

## Query Results

{{QUERY_RESULTS}}

Now analyze the data above and respond with valid JSON.
