## Vacation Rental Marketplace Context

This is a **vacation rental marketplace** connecting property owners (hosts) with travelers (guests). Key aspects to explore:

- **Two-sided marketplace**: Both host and guest behavior matter. A great guest experience requires responsive hosts with quality listings. Analyze both sides of the marketplace.
- **Property diversity**: Listings range from apartments and houses to villas, cabins, and unique stays. Segment analysis by property type — what works for a city apartment is different from a beachfront villa.
- **Seasonal demand**: Travel is highly seasonal. Look for peak/off-peak patterns, holiday surges, and shoulder season opportunities. Compare destination seasonality.
- **Booking flow**: Guests typically search → view listings → compare → book. Some platforms use instant booking, others use request-to-book. Analyze the full funnel by device and channel.
- **Pricing dynamics**: Nightly rates vary by property type, location, season, and demand. Look for pricing patterns, competitive gaps, and revenue optimization opportunities.
- **Host ecosystem**: Some hosts manage one property, others manage dozens. Professional hosts vs individual hosts may have very different performance profiles.
- **Trust signals**: Reviews, verification badges, superhost status, and response rates all affect booking conversion. Analyze how these trust signals impact guest behavior.
- **Cancellation patterns**: Both guest and host cancellations affect the platform. Identify cancellation timing, reasons, and their impact on revenue and guest satisfaction.
- **Support operations**: Customer support tickets reveal friction points. Analyze ticket volume by category, resolution time, and correlation with review scores.
- **Destination patterns**: Some destinations are drive-to weekend getaways, others are fly-to vacation spots. Booking lead times, stay durations, and price sensitivity vary significantly.

### Vacation Rental Example Queries

**Property Type Performance**:
```sql
SELECT pr.property_type,
       COUNT(DISTINCT pr.property_id) AS total_listings,
       COUNT(DISTINCT b.booking_id) AS total_bookings,
       AVG(b.total_price) AS avg_booking_value,
       AVG(r.rating) AS avg_review_score
FROM `{{DATASET}}.properties` pr
LEFT JOIN `{{DATASET}}.bookings` b ON pr.property_id = b.property_id
LEFT JOIN `{{DATASET}}.reviews` r ON b.booking_id = r.booking_id
{{FILTER}}
GROUP BY pr.property_type
ORDER BY total_bookings DESC
```

**Host Response Analysis**:
```sql
SELECT
  CASE
    WHEN response_time_hours <= 1 THEN 'within_1h'
    WHEN response_time_hours <= 4 THEN '1-4h'
    WHEN response_time_hours <= 24 THEN '4-24h'
    ELSE 'over_24h'
  END AS response_bucket,
  COUNT(DISTINCT host_id) AS host_count,
  AVG(acceptance_rate) AS avg_acceptance_rate,
  AVG(avg_review_score) AS avg_review
FROM `{{DATASET}}.hosts`
{{FILTER}}
GROUP BY response_bucket
ORDER BY host_count DESC
```

**Destination Seasonality**:
```sql
SELECT d.destination_name,
       DATE_TRUNC(b.check_in_date, MONTH) AS month,
       COUNT(DISTINCT b.booking_id) AS bookings,
       AVG(b.total_price) AS avg_price
FROM `{{DATASET}}.bookings` b
JOIN `{{DATASET}}.properties` pr ON b.property_id = pr.property_id
JOIN `{{DATASET}}.destinations` d ON pr.destination_id = d.destination_id
{{FILTER}}
GROUP BY d.destination_name, month
ORDER BY d.destination_name, month
```

**Cancellation Timing**:
```sql
SELECT
  CASE
    WHEN DATEDIFF(check_in_date, cancellation_date) > 30 THEN '30+ days before'
    WHEN DATEDIFF(check_in_date, cancellation_date) > 7 THEN '7-30 days before'
    WHEN DATEDIFF(check_in_date, cancellation_date) > 1 THEN '1-7 days before'
    ELSE 'same day or after'
  END AS cancellation_window,
  COUNT(*) AS cancellations,
  SUM(total_price) AS revenue_lost
FROM `{{DATASET}}.bookings`
WHERE status = 'cancelled'
{{FILTER}}
GROUP BY cancellation_window
ORDER BY cancellations DESC
```

**Guest Browsing to Booking Path**:
```sql
SELECT
  user_id,
  COUNT(CASE WHEN event_type = 'search' THEN 1 END) AS searches,
  COUNT(CASE WHEN event_type = 'view' THEN 1 END) AS views,
  COUNT(CASE WHEN event_type = 'book' THEN 1 END) AS bookings,
  COUNT(DISTINCT property_id) AS properties_viewed,
  MIN(event_timestamp) AS first_action,
  MAX(event_timestamp) AS last_action
FROM `{{DATASET}}.clickstream`
{{FILTER}}
  AND event_date >= DATE_SUB(CURRENT_DATE(), INTERVAL 30 DAY)
GROUP BY user_id
HAVING bookings > 0
LIMIT 100
```
