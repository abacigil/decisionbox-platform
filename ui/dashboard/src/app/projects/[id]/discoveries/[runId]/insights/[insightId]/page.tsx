'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import {
  Badge, Button, Card, Code, Group, Loader, Stack, Table, Text, Title,
} from '@mantine/core';
import {
  IconAlertTriangle, IconArrowLeft, IconCheck, IconX,
} from '@tabler/icons-react';
import Link from 'next/link';
import Shell from '@/components/layout/AppShell';
import { api, Insight } from '@/lib/api';

const severityColor: Record<string, string> = {
  critical: 'red', high: 'orange', medium: 'yellow', low: 'gray',
};

export default function InsightDetailPage() {
  const { id, runId, insightId } = useParams<{ id: string; runId: string; insightId: string }>();
  const [insight, setInsight] = useState<Insight | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getDiscoveryById(runId)
      .then((discovery) => {
        const insights = discovery?.insights || [];
        // Find by ID or index
        const found = insights.find((i) => i.id === insightId) || insights[parseInt(insightId)] || null;
        setInsight(found);
      })
      .catch(() => null)
      .finally(() => setLoading(false));
  }, [runId, insightId]);

  if (loading) return <Shell><Loader /></Shell>;
  if (!insight) return <Shell><Text>Insight not found</Text></Shell>;

  return (
    <Shell>
      <Stack gap="lg" maw={800}>
        {/* Back */}
        <Button variant="subtle" component={Link}
          href={`/projects/${id}/discoveries/${runId}`}
          leftSection={<IconArrowLeft size={16} />} size="sm" w="fit-content">
          Back to Discovery
        </Button>

        {/* Header */}
        <div>
          <Group gap="sm" mb={4}>
            <IconAlertTriangle size={20}
              color={`var(--mantine-color-${severityColor[insight.severity] || 'gray'}-6)`} />
            <Title order={2}>{insight.name}</Title>
          </Group>
          <Group gap="xs">
            <Badge color={severityColor[insight.severity] || 'gray'} variant="light">
              {insight.severity}
            </Badge>
            <Badge variant="outline">{insight.analysis_area}</Badge>
            {insight.affected_count > 0 && (
              <Badge variant="outline">{insight.affected_count.toLocaleString()} affected</Badge>
            )}
          </Group>
        </div>

        {/* Description */}
        <Card withBorder p="lg">
          <Text size="sm">{insight.description}</Text>
        </Card>

        {/* Indicators */}
        {insight.indicators && insight.indicators.length > 0 && (
          <Card withBorder p="lg">
            <Title order={4} mb="sm">Key Indicators</Title>
            <Stack gap={6}>
              {insight.indicators.map((ind, i) => (
                <Group key={i} gap="xs">
                  <Text size="xs" c="dimmed">-</Text>
                  <Text size="sm">{ind}</Text>
                </Group>
              ))}
            </Stack>
          </Card>
        )}

        {/* Metrics */}
        {insight.metrics && Object.keys(insight.metrics).length > 0 && (
          <Card withBorder p="lg">
            <Title order={4} mb="sm">Metrics</Title>
            <Table>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Metric</Table.Th>
                  <Table.Th>Value</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {Object.entries(insight.metrics).map(([key, value]) => (
                  <Table.Tr key={key}>
                    <Table.Td><Text size="sm">{key.replace(/_/g, ' ')}</Text></Table.Td>
                    <Table.Td><Text size="sm" fw={600}>{String(value)}</Text></Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          </Card>
        )}

        {/* Scores */}
        <Card withBorder p="lg">
          <Title order={4} mb="sm">Assessment</Title>
          <Group gap="xl">
            <div>
              <Text size="xs" c="dimmed">Risk Score</Text>
              <Text size="lg" fw={700} c={insight.risk_score > 0.7 ? 'red' : insight.risk_score > 0.4 ? 'orange' : 'green'}>
                {(insight.risk_score * 100).toFixed(0)}%
              </Text>
            </div>
            <div>
              <Text size="xs" c="dimmed">Confidence</Text>
              <Text size="lg" fw={700}>{(insight.confidence * 100).toFixed(0)}%</Text>
            </div>
            {insight.target_segment && (
              <div>
                <Text size="xs" c="dimmed">Target Segment</Text>
                <Text size="sm">{insight.target_segment}</Text>
              </div>
            )}
          </Group>
        </Card>

        {/* Validation */}
        {insight.validation && (
          <Card withBorder p="lg">
            <Group mb="sm">
              <Title order={4}>Validation</Title>
              <Badge
                color={insight.validation.status === 'confirmed' ? 'green' :
                       insight.validation.status === 'adjusted' ? 'yellow' :
                       insight.validation.status === 'rejected' ? 'red' : 'gray'}
                leftSection={insight.validation.status === 'confirmed' ? <IconCheck size={12} /> : <IconX size={12} />}>
                {insight.validation.status}
              </Badge>
            </Group>

            {(insight.validation.original_count || insight.validation.verified_count) && (
              <Group gap="xl" mb="sm">
                {insight.validation.original_count != null && (
                  <div>
                    <Text size="xs" c="dimmed">Claimed Count</Text>
                    <Text size="sm" fw={600}>{insight.validation.original_count.toLocaleString()}</Text>
                  </div>
                )}
                {insight.validation.verified_count != null && (
                  <div>
                    <Text size="xs" c="dimmed">Verified Count</Text>
                    <Text size="sm" fw={600}>{insight.validation.verified_count.toLocaleString()}</Text>
                  </div>
                )}
              </Group>
            )}

            {insight.validation.reasoning && (
              <Text size="xs" c="dimmed">{insight.validation.reasoning}</Text>
            )}
          </Card>
        )}

        {/* Discovered At */}
        {insight.discovered_at && (
          <Text size="xs" c="dimmed">
            Discovered: {new Date(insight.discovered_at).toLocaleString()}
          </Text>
        )}
      </Stack>
    </Shell>
  );
}
