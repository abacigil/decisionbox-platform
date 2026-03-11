'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import {
  Badge, Button, Card, Chip, Grid, Group, Loader, Stack, Text, Title,
} from '@mantine/core';
import {
  IconAlertTriangle, IconArrowLeft, IconBulb, IconTrendingUp,
} from '@tabler/icons-react';
import Link from 'next/link';
import Shell from '@/components/layout/AppShell';
import { api, DiscoveryResult, Insight, Recommendation } from '@/lib/api';

const severityColor: Record<string, string> = {
  critical: 'red', high: 'orange', medium: 'yellow', low: 'gray',
};
const severityOrder: Record<string, number> = {
  critical: 0, high: 1, medium: 2, low: 3,
};

export default function DiscoveryDetailPage() {
  const { id, runId } = useParams<{ id: string; runId: string }>();
  const [discovery, setDiscovery] = useState<DiscoveryResult | null>(null);
  const [loading, setLoading] = useState(true);

  // Filters
  const [areaFilter, setAreaFilter] = useState<string[]>([]);
  const [severityFilter, setSeverityFilter] = useState<string[]>([]);

  useEffect(() => {
    api.getDiscoveryById(runId)
      .then(setDiscovery)
      .catch(() => null)
      .finally(() => setLoading(false));
  }, [runId]);

  if (loading) return <Shell><Loader /></Shell>;
  if (!discovery) return <Shell><Text>Discovery not found</Text></Shell>;

  // Get unique areas and severities
  const allAreas = [...new Set((discovery.insights || []).map((i) => i.analysis_area))];
  const allSeverities = [...new Set((discovery.insights || []).map((i) => i.severity))];

  // Apply filters
  let filtered = discovery.insights || [];
  if (areaFilter.length > 0) {
    filtered = filtered.filter((i) => areaFilter.includes(i.analysis_area));
  }
  if (severityFilter.length > 0) {
    filtered = filtered.filter((i) => severityFilter.includes(i.severity));
  }

  // Sort by severity then risk score
  filtered.sort((a, b) => {
    const sevDiff = (severityOrder[a.severity] || 9) - (severityOrder[b.severity] || 9);
    if (sevDiff !== 0) return sevDiff;
    return b.risk_score - a.risk_score;
  });

  const durationSec = discovery.duration ? Math.round(discovery.duration / 1000000000) : 0;

  return (
    <Shell>
      <Stack gap="lg">
        {/* Header */}
        <Group>
          <Button variant="subtle" component={Link} href={`/projects/${id}`}
            leftSection={<IconArrowLeft size={16} />} size="sm">Back</Button>
        </Group>

        <Group justify="space-between">
          <div>
            <Title order={2}>
              {new Date(discovery.discovery_date).toLocaleDateString('en-US', {
                month: 'long', day: 'numeric', year: 'numeric',
              })}
            </Title>
            <Group gap="xs" mt={4}>
              <Badge variant="light" color={discovery.run_type === 'partial' ? 'violet' : 'blue'}>
                {discovery.run_type || 'full'}
              </Badge>
              {discovery.areas_requested && discovery.areas_requested.length > 0 && (
                <Text size="sm" c="dimmed">{discovery.areas_requested.join(', ')}</Text>
              )}
              <Text size="sm" c="dimmed">{discovery.total_steps} steps</Text>
              {durationSec > 0 && <Text size="sm" c="dimmed">{durationSec}s</Text>}
            </Group>
          </div>
        </Group>

        {/* KPI Row */}
        <Grid>
          <Grid.Col span={{ base: 6, md: 3 }}>
            <Card withBorder p="md" ta="center">
              <Text size="xl" fw={700} c="blue">{discovery.summary?.total_insights || 0}</Text>
              <Text size="sm" c="dimmed">Insights</Text>
            </Card>
          </Grid.Col>
          <Grid.Col span={{ base: 6, md: 3 }}>
            <Card withBorder p="md" ta="center">
              <Text size="xl" fw={700} c="red">
                {(discovery.insights || []).filter((i) => i.severity === 'critical').length}
              </Text>
              <Text size="sm" c="dimmed">Critical</Text>
            </Card>
          </Grid.Col>
          <Grid.Col span={{ base: 6, md: 3 }}>
            <Card withBorder p="md" ta="center">
              <Text size="xl" fw={700} c="violet">{discovery.summary?.total_recommendations || 0}</Text>
              <Text size="sm" c="dimmed">Recommendations</Text>
            </Card>
          </Grid.Col>
          <Grid.Col span={{ base: 6, md: 3 }}>
            <Card withBorder p="md" ta="center">
              <Text size="xl" fw={700} c="green">{discovery.summary?.queries_executed || 0}</Text>
              <Text size="sm" c="dimmed">Queries</Text>
            </Card>
          </Grid.Col>
        </Grid>

        {/* Filters */}
        {(discovery.insights || []).length > 0 && (
          <Card withBorder p="md">
            <Group gap="lg">
              <div>
                <Text size="xs" fw={600} mb={4}>Area</Text>
                <Chip.Group multiple value={areaFilter} onChange={setAreaFilter}>
                  <Group gap={4}>
                    {allAreas.map((area) => (
                      <Chip key={area} value={area} size="xs" variant="outline">
                        {area}
                      </Chip>
                    ))}
                  </Group>
                </Chip.Group>
              </div>
              <div>
                <Text size="xs" fw={600} mb={4}>Severity</Text>
                <Chip.Group multiple value={severityFilter} onChange={setSeverityFilter}>
                  <Group gap={4}>
                    {allSeverities.sort((a, b) => (severityOrder[a] || 9) - (severityOrder[b] || 9)).map((sev) => (
                      <Chip key={sev} value={sev} size="xs" variant="outline"
                        color={severityColor[sev] || 'gray'}>
                        {sev}
                      </Chip>
                    ))}
                  </Group>
                </Chip.Group>
              </div>
              {(areaFilter.length > 0 || severityFilter.length > 0) && (
                <Button variant="subtle" size="xs" onClick={() => { setAreaFilter([]); setSeverityFilter([]); }}>
                  Clear filters
                </Button>
              )}
            </Group>
          </Card>
        )}

        {/* Insights Report */}
        {filtered.length > 0 && (
          <>
            <Group justify="space-between">
              <Title order={3}>Insights ({filtered.length})</Title>
            </Group>
            <Stack gap="sm">
              {filtered.map((insight, idx) => (
                <Card key={idx} withBorder p="md" radius="md" component={Link}
                  href={`/projects/${id}/discoveries/${runId}/insights/${insight.id || idx}`}
                  style={{ textDecoration: 'none', cursor: 'pointer',
                    borderLeft: `3px solid var(--mantine-color-${severityColor[insight.severity] || 'gray'}-6)` }}>
                  <Group justify="space-between" mb={4}>
                    <Group gap="xs">
                      <IconAlertTriangle size={14}
                        color={`var(--mantine-color-${severityColor[insight.severity] || 'gray'}-6)`} />
                      <Text size="sm" fw={600}>{insight.name}</Text>
                    </Group>
                    <Group gap="xs">
                      <Badge size="xs" color={severityColor[insight.severity] || 'gray'} variant="light">
                        {insight.severity}
                      </Badge>
                      <Badge size="xs" variant="outline">{insight.analysis_area}</Badge>
                      {insight.affected_count > 0 && (
                        <Text size="xs" c="dimmed">{insight.affected_count.toLocaleString()} affected</Text>
                      )}
                    </Group>
                  </Group>
                  <Text size="xs" c="dimmed" lineClamp={2}>{insight.description}</Text>
                </Card>
              ))}
            </Stack>
          </>
        )}

        {(discovery.insights || []).length === 0 && (
          <Card withBorder p="xl" ta="center">
            <Text c="dimmed">No insights found in this discovery run.</Text>
          </Card>
        )}

        {/* Recommendations */}
        {(discovery.recommendations || []).length > 0 && (
          <>
            <Title order={3}>
              <IconBulb size={20} style={{ verticalAlign: 'middle', marginRight: 8 }} />
              Recommendations ({discovery.recommendations.length})
            </Title>
            <Stack gap="sm">
              {discovery.recommendations
                .sort((a, b) => b.priority - a.priority)
                .map((rec, idx) => (
                  <RecommendationCard key={idx} rec={rec} />
                ))}
            </Stack>
          </>
        )}
      </Stack>
    </Shell>
  );
}

function RecommendationCard({ rec }: { rec: Recommendation }) {
  const priorityColor = rec.priority >= 5 ? 'red' : rec.priority >= 4 ? 'orange' : 'blue';
  return (
    <Card withBorder p="md" radius="md"
      style={{ borderLeft: `3px solid var(--mantine-color-${priorityColor}-6)` }}>
      <Group justify="space-between" mb={4}>
        <Text size="sm" fw={600}>{rec.title}</Text>
        <Badge color={priorityColor} variant="light" size="xs">P{rec.priority}</Badge>
      </Group>
      <Text size="xs" c="dimmed" mb="sm">{rec.description}</Text>
      {rec.expected_impact && (
        <Group gap="xs" mb="xs">
          <IconTrendingUp size={12} />
          <Text size="xs" c="dimmed">
            {rec.expected_impact.metric}: {rec.expected_impact.estimated_improvement}
          </Text>
        </Group>
      )}
      {rec.actions && rec.actions.length > 0 && (
        <Stack gap={2}>
          {rec.actions.slice(0, 3).map((action, i) => (
            <Text key={i} size="xs" c="dimmed">- {action}</Text>
          ))}
        </Stack>
      )}
    </Card>
  );
}
